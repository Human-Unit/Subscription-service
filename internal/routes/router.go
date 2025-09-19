package routes

import (
	"github.com/gin-gonic/gin"
	"service/internal/handlers"
)

func SetupRouter() *gin.Engine {
	router:= gin.Default()
	router.POST("/Create", handlers.CreateEntry )
	router.GET("/GetEntry/:id",)
	router.PUT("/UpdateEntry/:id",)
	router.DELETE("/DeleteEntry/:id",)
	router.GET("/ListEntries",)

	router.GET("/summarise")
	return router
}