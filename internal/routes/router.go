package routes

import (
	"github.com/gin-gonic/gin"
	"service/internal/handlers"
)

func SetupRoutes(r *gin.Engine) {
    r.POST("/subscriptions", handlers.CreateSubscription)
    r.GET("/subscriptions/:id", handlers.GetSubscription)
    r.PUT("/subscriptions/:id", handlers.UpdateSubscription)
    r.DELETE("/subscriptions/:id", handlers.DeleteSubscription)
    r.GET("/subscriptions", handlers.ListSubscriptions)
    r.GET("/subscriptions/summary", handlers.GetSummary)
}