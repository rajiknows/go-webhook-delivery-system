package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rajiknows/webhook-mock/internal/models"
	// "github.com/rajiknows/webhook-mock/internal/services"
	"github.com/rajiknows/webhook-mock/pkg/db"
	"github.com/rajiknows/webhook-mock/pkg/queue"
)

func IngestWebhook(c *gin.Context) {
	// Ingest an incoming webhook payload, acknowledge immediately, and queue for asynchronous delivery.

	// 1. Parse subscription ID from URL parameter
	subscriptionID, err := uuid.Parse(c.Param("subscription_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	// 2. Retrieve subscription from database
	var subscription models.Subscription
	if err := db.DB.First(&subscription, "id = ?", subscriptionID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	// 3. Check if subscription is active
	if !subscription.Active {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found or inactive"})
		return
	}

	// 4. Read raw payload from request body
	rawPayload, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read payload"})
		return
	}

	// 5. Parse JSON payload
	var payload models.IncomingWebhook
	if err := json.Unmarshal(rawPayload, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	// 6. Event Type Filtering (Bonus)
	eventType := c.GetHeader("X-Event-Type")
	if subscription.EventTypeFilter != "" && eventType != subscription.EventTypeFilter {
		// Log that this event type was filtered out (for now, print to console)
		println("Event type '" + eventType + "' filtered out for subscription " + subscriptionID.String())
		// Return 202 Accepted, as the ingestion was processed but not queued
		c.JSON(http.StatusAccepted, gin.H{"message": "Webhook received, but filtered out by subscription settings"})
		return
	}

	// 7. Create IncomingWebhook in database
	incomingWebhook := models.IncomingWebhook{
		SubscriptionID: subscriptionID,
		Payload:        string(rawPayload),
		EventType:      eventType,
	}

	if err := db.DB.Create(&incomingWebhook).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create incoming webhook"})
		return
	}

	// 8. Enqueue delivery task
	if err = queue.EnqueueDeliveryTask(incomingWebhook.ID, 5); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enqueue delivery task"})
		return
	}

	// 9. Return 202 Accepted with delivery task ID
	c.JSON(http.StatusAccepted, gin.H{
		"message":          "Webhook received and queued for delivery",
		"delivery_task_id": incomingWebhook.ID,
	})
}
