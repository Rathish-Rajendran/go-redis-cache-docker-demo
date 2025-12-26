package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedisClient() {
	addr := os.Getenv("REDIS_ADDR")
	if addr != ""{
		addr = "redis:6379"
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password
		DB:       0,  // use default DB
		Protocol: 2,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("❌ Failed to connect to Redis: %v", err)
	}

	fmt.Println("✅ Connected to Redis at", addr)
}