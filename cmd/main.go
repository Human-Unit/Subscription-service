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
	"service/internal/logger"

	"github.com/gin-gonic/gin"
    ginSwagger "github.com/swaggo/gin-swagger"
    swaggerFiles "github.com/swaggo/files"
     _ "service/docs"
	"go.uber.org/zap" 			
)

func main() {
	cfg := config.LoadConfig()

	logger.Init()

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

func LoggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        path := c.Request.URL.Path
        method := c.Request.Method
        logger.Log.Info("Incoming request",
            zap.String("method", method),
            zap.String("path", path),
        )

        c.Next()

        status := c.Writer.Status()
        logger.Log.Info("Request completed",
            zap.String("method", method),
            zap.String("path", path),
            zap.Int("status", status),
        )
    }
}


