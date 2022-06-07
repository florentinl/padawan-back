package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	handler := newHandler()

	router := gin.Default()
	router.Use(CORSMiddleware())
	router.GET("/api/containers", handler.getAllContainers)

	router.GET("/api/login", handler.login)

	containerRoute := router.Group("/api/container")
	containerRoute.GET("/me", handler.getMyContainer)
	containerRoute.POST("/me", handler.postMyContainer)
	containerRoute.DELETE("/me", handler.deleteMyContainer)
	containerRoute.GET("/:username", handler.getUserContainer)
	containerRoute.POST("/:username", handler.postUserContainer)
	containerRoute.DELETE("/:username", handler.deleteUserContainer)

	imageRoute := router.Group("/api/images")
	imageRoute.GET(":image_name", handler.getImage)
	imageRoute.GET("", handler.getImages)
	imageRoute.POST("", handler.postImage)
	imageRoute.PUT(":image_name", handler.putImage)
	imageRoute.DELETE(":image_name", handler.deleteImage)

	router.Run("0.0.0.0:3000")
}
