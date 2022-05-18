package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	handler := newHandler()
	handler.db.AutoMigrate(&Image{})

	router := gin.Default()
	router.POST("/container", handler.postContainer)
	router.GET("/container", handler.existsContainer)
	router.DELETE("/container", handler.deleteContainer)
	router.Run("0.0.0.0:3000")
}
