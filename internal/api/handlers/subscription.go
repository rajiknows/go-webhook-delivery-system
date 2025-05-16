package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rajiknows/webhook-mock/internal/models"
	"github.com/rajiknows/webhook-mock/internal/services"
	"github.com/rajiknows/webhook-mock/pkg/db"
)

func CreateSubscription(c *gin.Context) {
	var sub models.Subscription
	if err := c.ShouldBindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.DB.Create(&sub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
		return
	}
	services.CacheSubscription(sub)
	c.JSON(http.StatusCreated, sub)
}

func GetSubscription(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}
	var sub models.Subscription
	if err := db.DB.First(&sub, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}
	c.JSON(http.StatusOK, sub)
}

func UpdateSubscription(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}
	var sub models.Subscription
	if err := db.DB.First(&sub, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}
	if err := c.ShouldBindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.DB.Save(&sub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription"})
		return
	}
	services.DeleteSubscriptionCache(id)
	c.JSON(http.StatusOK, sub)
}

func DeleteSubscription(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}
	if err := db.DB.Delete(&models.Subscription{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete subscription"})
		return
	}
	services.DeleteSubscriptionCache(id)
	c.Status(http.StatusNoContent)
}

func ListSubscriptionDeliveries(c *gin.Context) {
	subscriptionID, err := uuid.Parse(c.Param("subscription_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription ID"})
		return
	}

	var deliveries []models.DeliveryLog
	if err := db.DB.
		Where("subscription_id = ? AND status = ?", subscriptionID, "completed").
		Find(&deliveries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch deliveries"})
		return
	}

	c.JSON(http.StatusOK, deliveries)
}
