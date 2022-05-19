package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getImages(c *gin.Context) {
	var images []Image
	h.db.Find(&images)
	c.IndentedJSON(http.StatusOK, images)
}

func (h *Handler) getImage(c *gin.Context) {
	imageName := c.Param("image_name")
	var image Image
	if result := h.db.Where("image_name = ?", imageName).First(&image); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "image not found",
		})
		return
	}
	c.IndentedJSON(http.StatusOK, image)
}

func (h *Handler) postImage(c *gin.Context) {
	username := getUsername(c)
	if !isAdmin(username, h) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "not allowed",
		})
		return
	}
	var image Image
	if err := c.BindJSON(&image); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if result := h.db.Create(&image); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error,
		})
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) putImage(c *gin.Context) {
	username := getUsername(c)
	if !isAdmin(username, h) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "not allowed",
		})
		return
	}
	imageName := c.Param("image_name")
	var image Image
	if result := h.db.Where("image_name = ?", imageName).First(&image); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "image not found",
		})
		return
	}
	if err := c.BindJSON(&image); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if result := h.db.Save(&image); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error,
		})
		return
	}
	c.Status(http.StatusOK)
}

func (h *Handler) deleteImage(c *gin.Context) {
	username := getUsername(c)
	if !isAdmin(username, h) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "not allowed",
		})
		return
	}
	imageName := c.Param("image_name")
	var image Image
	if result := h.db.Where("image_name = ?", imageName).First(&image); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "image not found",
		})
		return
	}
	if result := h.db.Delete(&image); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error,
		})
		return
	}
	c.Status(http.StatusOK)
}
