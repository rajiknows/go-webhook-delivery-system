package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rajiknows/webhook-mock/internal/models"
	"github.com/rajiknows/webhook-mock/pkg/db"
)

func GetWebhookStatus(c *gin.Context) {
	webhookID, err := uuid.Parse(c.Param("incoming_webhook_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "please provide a valid incoming webhook id"})
		return
	}

	var webhook models.IncomingWebhook
	if err := db.DB.First(&webhook, "id = ?", webhookID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "webhook not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
