package handlers

import (
	"fmt"
	"log"
	"net/http"
	"service/internal/database"
	"service/internal/models"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"fmt"
)

func CreateEntry(c *gin.Context){
	var entry models.Subscription
	if err := c.ShouldBindJSON(&entry); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error":"Invalid input format",
			"details": err.Error(),
		})
		log.Printf("Accured error %s", err.Error())
	}
	db := database.InitDB()
    if db == nil {
        log.Println("Database connection is nil")
        status.Error(codes.Internal, "database connection error")
    }
	result := db.Create(&entry)
    if result.Error != nil {
        log.Printf("Failed to create entry: %v", result.Error)
        status.Error(codes.Internal, "failed to create entry")
    }
}

func GetEntryById(c *gin.Context) {
	id := c.Query("id")
	var entry models.Subscription
	if id == "" {
    	status.Error(codes.InvalidArgument, "entry ID is required")
    }
	db := database.InitDB()
    if db == nil {
        log.Println("Database connection is nil")
        status.Error(codes.Internal, "database connection error")
    }
	result:= db.Where("id = ?", id).First(&entry)
	if result ==nil{
		fmt.Print("lox")
	}
	c.JSON(200, entry)
}