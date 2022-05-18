package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	handler := newHandler()
	handler.db.AutoMigrate(&Image{})

	router := gin.Default()
	router.POST("/container", handler.postContainer)
	router.GET("/container", handler.getContainer)
	router.DELETE("/container", handler.deleteContainer)
	router.GET("/images", handler.getImages)
	router.GET("/images/:image_name", handler.getImage)
	router.POST("/images", handler.postImage)
	router.PUT("/images/:image_name", handler.putImage)
	router.DELETE("/images/:image_name", handler.deleteImage)
	router.Run("0.0.0.0:3000")
}
