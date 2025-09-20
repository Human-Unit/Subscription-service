package handlers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"

    "service/internal/database"
    "service/internal/models"
)

 

func CreateSubscription(c *gin.Context) {
	db := database.GetDB()
    var sub models.Subscription
    if err := c.ShouldBindJSON(&sub); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

 
    if sub.ID == uuid.Nil {
        sub.ID = uuid.New()
    }
    if err := db.Create(&sub).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, sub)
}

 
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

 
func UpdateSubscription(c *gin.Context) {
	db := database.GetDB()
    id := c.Param("id")
    var sub models.Subscription

    if err := db.First(&sub, "id = ?", id).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
        return
    }

    var input models.Subscription
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

 
    sub.ServiceName = input.ServiceName
    sub.Price = input.Price
    sub.UserID = input.UserID
    sub.StartDate = input.StartDate
    sub.EndDate = input.EndDate

    if err := db.Save(&sub).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, sub)
}

 
func DeleteSubscription(c *gin.Context) {
	db := database.GetDB()
    id := c.Param("id")
    if err := db.Delete(&models.Subscription{}, "id = ?", id).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

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
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, subs)
}

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
        t, _ := time.Parse("2006-01", start) // YYYY-MM
        query = query.Where("start_date >= ?", t)
    }
    if end := c.Query("end_date"); end != "" {
        t, _ := time.Parse("2006-01", end)
        query = query.Where("end_date <= ?", t)
    }

    if err := query.Select("SUM(price)").Scan(&total).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"total": total})
}
