package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getContainer(username string, h *Handler, c *gin.Context) {
	var container Container
	if result := h.db.Where("username = ?", username).First(&container); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "container not found",
		})
		return
	}
	c.IndentedJSON(http.StatusOK, container)
}

func postContainer(username string, h *Handler, c *gin.Context) {
	var request ContainerRequest

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	if result := h.db.Where("image_name = ?", request.ImageName).First(&Image{}); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "image not found",
		})
		return
	}

	newContainer, err := createResources(username, request.ImageName, request.Password, h)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if result := h.db.Create(&newContainer); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error,
		})
		return
	}

	c.IndentedJSON(http.StatusOK, newContainer)
}

func deleteContainer(username string, h *Handler, c *gin.Context) {
	err := deleteResources(username, h)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Remove resource from database
	var container Container
	if result := h.db.Where("username = ?", username).First(&container); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "container not found",
		})
		return
	}
	if result := h.db.Delete(&container); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error,
		})
		return
	}

	c.Status(http.StatusOK)
}
