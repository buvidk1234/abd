package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
)

// This script clears all test data from Kafka, Redis, and MySQL.
// Run it before starting a fresh stress test.

const (
	mysqlDSN      = "root:123@tcp(localhost:3306)/abd"
	redisAddr     = "localhost:6379"
	redisPassword = "123456"
)

func main() {
	fmt.Println("=== Starting Environment Cleanup ===")

	// 1. Clear Redis
	clearRedis()

	// 2. Clear MySQL
	clearMySQL()

	// 3. Clear Kafka (via existing script or direct command)
	clearKafka()

	fmt.Println("\n=== Cleanup Complete. Environment is Ready for Testing. ===")
}

func clearRedis() {
	fmt.Print("-> Clearing Redis... ")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0,
	})
	err := rdb.FlushAll(context.Background()).Err()
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		fmt.Println("SUCCESS")
	}
}

func clearMySQL() {
	fmt.Print("-> Clearing MySQL Tables... ")
	db, err := sql.Open("mysql", mysqlDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Disable foreign key checks for truncation
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	defer db.Exec("SET FOREIGN_KEY_CHECKS = 1")

	tables := []string{
		"messages",
		"conversations",
		"seq_conversations",
		"seq_users",
		"user_timelines",
		"users",
		"friends",
		"friend_applies",
		"groups",
		"group_members",
	}

	for _, table := range tables {
		// Use backticks for reserved words like 'groups'
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE `%s`", table))
		if err != nil {
			// Some tables might not exist yet if it's the first run
			if !strings.Contains(err.Error(), "doesn't exist") {
				fmt.Printf("\nWarning: Failed to truncate %s: %v", table, err)
			}
		}
	}
	fmt.Println("SUCCESS")
}

func clearKafka() {
	fmt.Print("-> Clearing Kafka Topics... ")
	// Using the existing clean_kafka.go script logic
	cmd := exec.Command("go", "run", "scripts/clean_kafka.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Printf("FAILED: %v\n", err)
	} else {
		fmt.Println("SUCCESS")
	}
}
