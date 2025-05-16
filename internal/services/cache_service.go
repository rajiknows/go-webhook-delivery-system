package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/rajiknows/webhook-mock/internal/models"
	"log"
	"os"
	"time"
)

var ctx = context.Background()
var rdb *redis.Client

func InitRedis() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatal("Failed to parse Redis URL:", err)
	}
	rdb = redis.NewClient(opt)
}

func CacheSubscription(sub models.Subscription) {
	key := fmt.Sprintf("subscription:%s", sub.ID)
	data, err := json.Marshal(sub)
	if err != nil {
		log.Println("Failed to marshal subscription:", err)
		return
	}
	if err := rdb.Set(ctx, key, data, 5*time.Minute).Err(); err != nil {
		log.Println("Failed to cache subscription:", err)
	}
}

func GetSubscriptionFromCache(id uuid.UUID) (*models.Subscription, error) {
	key := fmt.Sprintf("subscription:%s", id)
	data, err := rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var sub models.Subscription
	if err := json.Unmarshal(data, &sub); err != nil {
		return nil, err
	}
	return &sub, nil
}

func DeleteSubscriptionCache(id uuid.UUID) {
	key := fmt.Sprintf("subscription:%s", id)
	if err := rdb.Del(ctx, key).Err(); err != nil {
		log.Println("Failed to delete cache:", err)
	}
}
