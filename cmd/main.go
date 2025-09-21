// @title Subscription Service API
// @version 1.0
// @description This is the API documentation for the subscription service.
// @host localhost:8080
// @BasePath /
package main

import (
    "fmt"
    "log"

    "service/internal/config"
    "service/internal/routes"
    "service/internal/database"

	"github.com/gin-gonic/gin"
    ginSwagger "github.com/swaggo/gin-swagger"
    swaggerFiles "github.com/swaggo/files"
     _ "service/docs"
)

func main() {
	cfg := config.LoadConfig()

	db := database.InitDB(cfg.DBDsn)
	defer database.CloseDB()
	_ = db 

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Use(gin.Recovery())
	routes.SetupRoutes(router) 

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("🚀 Starting server on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}


