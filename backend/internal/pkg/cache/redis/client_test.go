package redis

import (
	"context"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User model for testing
type User struct {
	ID   string `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}
	err = db.AutoMigrate(&User{})
	if err != nil {
		t.Fatalf("failed to migrate database: %v", err)
	}
	return db
}

func TestClient(t *testing.T) {
	cfg := Config{
		Addr:         "192.168.6.130:6379",
		Password:     "123456",
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 2,
	}

	Init(cfg)
	ctx := context.Background()
	RDB.Set(ctx, "test_key", "test_value", 0)
	val, err := RDB.Get(ctx, "test_key").Result()
	if err != nil {
		t.Fatalf("failed to get value: %v", err)
	}
	if val != "test_value" {
		t.Fatalf("expected 'test_value', got '%s'", val)
	}
	t.Logf("value: %s", val)
}

func TestGetCache(t *testing.T) {
	cfg := Config{
		Addr:         "192.168.6.130:6379",
		Password:     "123456",
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 2,
	}
	Init(cfg)
	db := setupTestDB(t)

	// Insert test data into DB
	user := User{ID: "1", Name: "Alice", Age: 25}
	db.Create(&user)

	key := "user:1"
	// Ensure key doesn't exist in Redis
	RDB.Del(context.Background(), key)

	// 1. Test Cache Miss (Fetch from DB)
	called := false
	val, err := GetCache[User](key, func() (User, error) {
		called = true
		var u User
		if err := db.First(&u, "id = ?", "1").Error; err != nil {
			return User{}, err
		}
		return u, nil
	}, time.Minute)

	if err != nil {
		t.Fatalf("GetCache failed: %v", err)
	}
	if val.Name != "Alice" {
		t.Errorf("expected Alice, got %s", val.Name)
	}
	if !called {
		t.Error("expected fn to be called on cache miss")
	}

	// 2. Test Cache Hit (Should NOT call DB)
	called = false
	val, err = GetCache[User](key, func() (User, error) {
		called = true
		return User{}, nil
	}, time.Minute)

	if err != nil {
		t.Fatalf("GetCache failed: %v", err)
	}
	if val.Name != "Alice" {
		t.Errorf("expected Alice, got %s", val.Name)
	}
	if called {
		t.Error("expected fn NOT to be called on cache hit")
	}

	// 3. Test String Type
	strKey := "str_key"
	RDB.Del(context.Background(), strKey)
	strVal, err := GetCache[string](strKey, func() (string, error) {
		return "hello", nil
	}, time.Minute)
	if err != nil {
		t.Fatalf("GetCache string failed: %v", err)
	}
	if strVal != "hello" {
		t.Errorf("expected hello, got %s", strVal)
	}
}

func TestBatchGetCache(t *testing.T) {
	cfg := Config{
		Addr:         "192.168.6.130:6379",
		Password:     "123456",
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 2,
	}
	Init(cfg)
	db := setupTestDB(t)

	// Insert test data
	users := []User{
		{ID: "101", Name: "User101", Age: 20},
		{ID: "102", Name: "User102", Age: 21},
		{ID: "103", Name: "User103", Age: 22},
	}
	db.Create(&users)

	keys := []string{"user:101", "user:102", "user:103"}
	// Ensure keys don't exist
	RDB.Del(context.Background(), keys...)

	// 1. Test Batch Cache Miss
	called := false
	res, err := BatchGetCache[User](keys, func(missingKeys []string) (map[string]User, error) {
		called = true
		if len(missingKeys) != 3 {
			t.Errorf("expected 3 missing keys, got %d", len(missingKeys))
		}

		// Extract IDs from keys (e.g., "user:101" -> "101")
		var ids []string
		for _, k := range missingKeys {
			ids = append(ids, k[5:])
		}

		var dbUsers []User
		if err := db.Find(&dbUsers, "id IN ?", ids).Error; err != nil {
			return nil, err
		}

		ret := make(map[string]User)
		for _, u := range dbUsers {
			ret["user:"+u.ID] = u
		}
		return ret, nil
	}, time.Minute)

	if err != nil {
		t.Fatalf("BatchGetCache failed: %v", err)
	}
	if len(res) != 3 {
		t.Errorf("expected 3 results, got %d", len(res))
	}
	if res["user:101"].Name != "User101" {
		t.Errorf("expected User101, got %s", res["user:101"].Name)
	}
	if !called {
		t.Error("expected fn to be called")
	}

	// 2. Test Batch Cache Hit
	called = false
	res, err = BatchGetCache[User](keys, func(missingKeys []string) (map[string]User, error) {
		called = true
		return nil, nil
	}, time.Minute)

	if err != nil {
		t.Fatalf("BatchGetCache failed: %v", err)
	}
	if len(res) != 3 {
		t.Errorf("expected 3 results, got %d", len(res))
	}
	if called {
		t.Error("expected fn NOT to be called on full cache hit")
	}

	// 3. Test Partial Hit
	RDB.Del(context.Background(), "user:102") // Delete one key
	called = false
	res, err = BatchGetCache[User](keys, func(missingKeys []string) (map[string]User, error) {
		called = true
		if len(missingKeys) != 1 || missingKeys[0] != "user:102" {
			t.Errorf("expected missing key user:102, got %v", missingKeys)
		}
		return map[string]User{
			"user:102": {ID: "102", Name: "User102", Age: 21},
		}, nil
	}, time.Minute)

	if err != nil {
		t.Fatalf("BatchGetCache partial failed: %v", err)
	}
	if len(res) != 3 {
		t.Errorf("expected 3 results, got %d", len(res))
	}
	if !called {
		t.Error("expected fn to be called for partial miss")
	}
}
