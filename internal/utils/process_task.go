package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"github.com/rajiknows/webhook-mock/internal/models"
	"github.com/rajiknows/webhook-mock/internal/services"
	"github.com/rajiknows/webhook-mock/pkg/db"
	"github.com/rajiknows/webhook-mock/pkg/queue"
)

func ProcessDeliveryTask(ctx context.Context, task *asynq.Task) error {
	var payload queue.DeliveryPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	// Retrieve IncomingWebhook
	var incomingWebhook models.IncomingWebhook
	if err := db.DB.First(&incomingWebhook, "id = ?", payload.IncomingWebhookID).Error; err != nil {
		return fmt.Errorf("failed to find incoming webhook: %w", err)
	}

	// Retrieve Subscription
	var subscription models.Subscription
	if err := db.DB.First(&subscription, "id = ?", incomingWebhook.SubscriptionID).Error; err != nil {
		return fmt.Errorf("failed to find subscription: %w", err)
	}

	// Attempt delivery
	_, httpStatusCode, err := services.DeliverWebhook(subscription.TargetURL, incomingWebhook.Payload)
	if err != nil {
		// Log failure
		log.Printf("Delivery attempt %d failed for webhook %s: %v", payload.AttemptNumber, incomingWebhook.ID, err)
		// Create delivery log
		deliveryLog := models.DeliveryLog{
			IncomingWebhookID: incomingWebhook.ID,
			AttemptNumber:     payload.AttemptNumber,
			Status:            "failed",
			HTTPStatusCode:    &httpStatusCode,
			ErrorDetails:      err.Error(),
		}
		if err := db.DB.Create(&deliveryLog).Error; err != nil {
			log.Printf("Failed to create delivery log: %v", err)
		}
		// Retry if attempts left
		// if payload.AttemptNumber < MaxRetries {
		// 	return
		// }
		return nil
	}

	// Log success
	log.Printf("Delivery attempt %d succeeded for webhook %s", payload.AttemptNumber, incomingWebhook.ID)
	// Create delivery log
	deliveryLog := models.DeliveryLog{
		IncomingWebhookID: incomingWebhook.ID,
		AttemptNumber:     payload.AttemptNumber,
		Status:            "success",
		HTTPStatusCode:    &httpStatusCode,
	}
	if err := db.DB.Create(&deliveryLog).Error; err != nil {
		log.Printf("Failed to create delivery log: %v", err)
	}
	return nil
}
