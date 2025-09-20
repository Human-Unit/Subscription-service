package main

import (
    "fmt"
    "log"

    "service/internal/config"
    "service/internal/routes"
    "service/internal/database"

    "github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	db := database.InitDB(cfg.DBDsn)
	defer database.CloseDB()
	_ = db 


	router := gin.Default()
	router.Use(gin.Recovery())
	routes.SetupRoutes(router) 

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("🚀 Starting server on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}


