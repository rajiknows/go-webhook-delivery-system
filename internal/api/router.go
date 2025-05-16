package api

import (
	"github.com/gin-gonic/gin"
	"github.com/rajiknows/webhook-mock/internal/api/handlers"
	// "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/subscriptions", handlers.CreateSubscription)
	r.GET("/subscriptions/:id", handlers.GetSubscription)
	r.PUT("/subscriptions/:id", handlers.UpdateSubscription)
	r.DELETE("/subscriptions/:id", handlers.DeleteSubscription)
	r.POST("/ingest/:subscription_id", handlers.IngestWebhook)
	r.GET("/status/:incoming_webhook_id", handlers.GetWebhookStatus)
	r.GET("/subscriptions/:subscription_id/deliveries", handlers.ListSubscriptionDeliveries)
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
