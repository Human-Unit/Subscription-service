package handlers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"

    "service/internal/database"
    "service/internal/models"
    "service/internal/logger"
    "go.uber.org/zap"
)

 
// @Summary Create a subscription
// @Description Create a new subscription for a user
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body models.Subscription true "Subscription info"
// @Success 201 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func CreateSubscription(c *gin.Context) {
    db := database.GetDB()
    var sub models.Subscription
    if err := c.ShouldBindJSON(&sub); err != nil {
        logger.Log.Error("Failed to bind subscription JSON", zap.Error(err))
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if sub.ID == uuid.Nil {
        sub.ID = uuid.New()
    }

    if err := db.Create(&sub).Error; err != nil {
        logger.Log.Error("Failed to create subscription", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    logger.Log.Info("Subscription created",
        zap.String("id", sub.ID.String()),
        zap.String("service", sub.ServiceName),
        zap.String("user_id", sub.UserID.String()),
        zap.Float64("price", float64(sub.Price)),
    )
    c.JSON(http.StatusCreated, sub)
}

 // @Summary Get a subscription
// @Description Get subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} models.Subscription
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [get]
func GetSubscription(c *gin.Context) {
	db := database.GetDB()
    id := c.Param("id")
    var sub models.Subscription

    if err := db.First(&sub, "id = ?", id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
        return
    }
    c.JSON(http.StatusOK, sub)
}

// @Summary Update a subscription
// @Description Update subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param subscription body models.Subscription true "Subscription info"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [put] 
func UpdateSubscription(c *gin.Context) {
    db := database.GetDB()
    id := c.Param("id")
    var sub models.Subscription

    if err := db.First(&sub, "id = ?", id).Error; err != nil {
        logger.Log.Warn("Subscription not found for update", zap.String("id", id))
        c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
        return
    }

    var input models.Subscription
    if err := c.ShouldBindJSON(&input); err != nil {
        logger.Log.Error("Failed to bind subscription JSON for update", zap.Error(err))
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    sub.ServiceName = input.ServiceName
    sub.Price = input.Price
    sub.UserID = input.UserID
    sub.StartDate = input.StartDate
    sub.EndDate = input.EndDate

    if err := db.Save(&sub).Error; err != nil {
        logger.Log.Error("Failed to update subscription", zap.Error(err), zap.String("id", id))
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    logger.Log.Info("Subscription updated",
        zap.String("id", id),
        zap.String("service", sub.ServiceName),
        zap.String("user_id", sub.UserID.String()),
    )
    c.JSON(http.StatusOK, sub)
}

// @Summary Delete a subscription
// @Description Delete subscription by ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [delete] 
func DeleteSubscription(c *gin.Context) {
    db := database.GetDB()
    id := c.Param("id")

    if err := db.Delete(&models.Subscription{}, "id = ?", id).Error; err != nil {
        logger.Log.Error("Failed to delete subscription", zap.Error(err), zap.String("id", id))
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    logger.Log.Info("Subscription deleted", zap.String("id", id))
    c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// @Summary List subscriptions
// @Description Get list of subscriptions with optional filters
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "User ID"
// @Param service_name query string false "Service name"
// @Success 200 {array} models.Subscription
// @Failure 500 {object} map[string]string
// @Router /subscriptions [get]
func ListSubscriptions(c *gin.Context) {
    db := database.GetDB()
    var subs []models.Subscription
    query := db

    if userID := c.Query("user_id"); userID != "" {
        query = query.Where("user_id = ?", userID)
    }
    if service := c.Query("service_name"); service != "" {
        query = query.Where("service_name ILIKE ?", "%"+service+"%")
    }

    if err := query.Find(&subs).Error; err != nil {
        logger.Log.Error("Failed to list subscriptions", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    logger.Log.Info("Subscriptions listed",
        zap.Int("count", len(subs)),
        zap.String("user_id", c.Query("user_id")),
        zap.String("service_name", c.Query("service_name")),
    )
    c.JSON(http.StatusOK, subs)
}

// @Summary Get summary
// @Description Get total price summary of subscriptions
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "User ID"
// @Param service_name query string false "Service name"
// @Param start_date query string false "Start date YYYY-MM"
// @Param end_date query string false "End date YYYY-MM"
// @Success 200 {object} map[string]int64
// @Failure 500 {object} map[string]string
// @Router /subscriptions/summary [get]
func GetSummary(c *gin.Context) {
    db := database.GetDB()
    var total int64
    query := db.Model(&models.Subscription{})

    if userID := c.Query("user_id"); userID != "" {
        query = query.Where("user_id = ?", userID)
    }
    if service := c.Query("service_name"); service != "" {
        query = query.Where("service_name ILIKE ?", "%"+service+"%")
    }
    if start := c.Query("start_date"); start != "" {
        t, _ := time.Parse("2006-01", start)
        query = query.Where("start_date >= ?", t)
    }
    if end := c.Query("end_date"); end != "" {
        t, _ := time.Parse("2006-01", end)
        query = query.Where("end_date <= ?", t)
    }

    if err := query.Select("SUM(price)").Scan(&total).Error; err != nil {
        logger.Log.Error("Failed to calculate summary", zap.Error(err))
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    logger.Log.Info("Summary calculated",
        zap.Int64("total", total),
        zap.String("user_id", c.Query("user_id")),
        zap.String("service_name", c.Query("service_name")),
        zap.String("start_date", c.Query("start_date")),
        zap.String("end_date", c.Query("end_date")),
    )

    c.JSON(http.StatusOK, gin.H{"total": total})
}
