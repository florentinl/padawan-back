package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	handler := newHandler()

	router := gin.Default()
	router.Use(CORSMiddleware())
	router.GET("/containers", handler.getAllContainers)

	router.LoadHTMLFiles("login.html")
	router.GET("/login", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	containerRoute := router.Group("/container")
	containerRoute.GET("/me", handler.getMyContainer)
	containerRoute.POST("/me", handler.postMyContainer)
	containerRoute.DELETE("/me", handler.deleteMyContainer)
	containerRoute.GET("/:username", handler.getUserContainer)
	containerRoute.POST("/:username", handler.postUserContainer)
	containerRoute.DELETE("/:username", handler.deleteUserContainer)

	imageRoute := router.Group("/images")
	imageRoute.GET(":image_name", handler.getImage)
	imageRoute.GET("", handler.getImages)
	imageRoute.POST("", handler.postImage)
	imageRoute.PUT(":image_name", handler.putImage)
	imageRoute.DELETE(":image_name", handler.deleteImage)

	router.Run("0.0.0.0:3000")
}
