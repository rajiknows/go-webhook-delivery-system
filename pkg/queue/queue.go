package queue

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

const (
	DeliveryTaskType = "webhook:delivery"
	MaxRetries       = 5
)

var Client *asynq.Client
var Server *asynq.Server

type DeliveryPayload struct {
	IncomingWebhookID uuid.UUID `json:"incoming_webhook_id"`
	AttemptNumber     int       `json:"attempt_number"`
}

func InitQueue() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "redis://localhost:6379"
	}
	Client = asynq.NewClient(asynq.RedisClientOpt{Addr: redisURL})
	Server = asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisURL},
		asynq.Config{
			Concurrency: 10,
			RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
				return time.Duration(10*n) * time.Second // Exponential backoff: 10s, 20s, 30s, etc.
			},
		},
	)
}

func EnqueueDeliveryTask(incomingWebhookID uuid.UUID, attemptNumber int) error {
	payload, err := json.Marshal(DeliveryPayload{
		IncomingWebhookID: incomingWebhookID,
		AttemptNumber:     attemptNumber,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	task := asynq.NewTask(DeliveryTaskType, payload)
	_, err = Client.Enqueue(task, asynq.MaxRetry(MaxRetries))
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}
	return nil
}
