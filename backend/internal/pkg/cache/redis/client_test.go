package redis

import (
	"context"
	"testing"
)

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

type AppConfig struct {
	Redis Config `yaml:"redis"`
}
