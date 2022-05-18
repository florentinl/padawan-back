package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) existsContainer(c *gin.Context) {
	username := c.GetHeader("X-Forwarded-User")
	port, err := getResources(username, h)
	c.IndentedJSON(http.StatusOK, gin.H{"exists": err == nil, "port": port})
}

type ContainerRequest struct {
	ImageName string `json:"image_name"`
	Password  string `json:"password"`
}

func (h *Handler) postContainer(c *gin.Context) {
	username := c.GetHeader("X-Forwarded-User")
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

	_, err := createResources(username, request.ImageName, request.Password, h)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	h.existsContainer(c)
}

func (h *Handler) deleteContainer(c *gin.Context) {
	username := c.GetHeader("X-Forwarded-User")
	fmt.Println(username)
	err := deleteResources(username, h)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}
