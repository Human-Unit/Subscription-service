package main

import (
	"fmt"
	"log"
	"net/http"

	"service/internal/config"
	"service/internal/routes"

	"github.com/gin-gonic/gin"
	"service/internal/database"
)

func main() {
	cfg := config.LoadConfig()

	if err := database.InitDB(); err != nil {
        log.Fatalf("Database connection failed: %v", err)
    }
	router := gin.Default()
	router.Use(gin.Recovery())
	routes.SetupRouter()
	
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Starting server on %s", addr)

	http.ListenAndServe(addr, nil)
}
