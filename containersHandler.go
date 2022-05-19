package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getAllContainers(c *gin.Context) {
	username := getUsername(c)
	if !isAdmin(username, h) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "not allowed",
		})
		return
	}
	var containers []Container
	if result := h.db.Find(&containers); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error,
		})
		return
	}
	c.IndentedJSON(http.StatusOK, containers)
}

func (h *Handler) getMyContainer(c *gin.Context) {
	username := getUsername(c)
	getContainer(username, h, c)
}

func (h *Handler) postMyContainer(c *gin.Context) {
	username := getUsername(c)
	postContainer(username, h, c)
}

func (h *Handler) deleteMyContainer(c *gin.Context) {
	username := getUsername(c)
	deleteContainer(username, h, c)
}

func (h *Handler) getUserContainer(c *gin.Context) {
	username := getUsername(c)
	if !isAdmin(username, h) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "not allowed",
		})
		return
	}
	user := c.Param("username")
	getContainer(user, h, c)

}

func (h *Handler) postUserContainer(c *gin.Context) {
	username := getUsername(c)
	if !isAdmin(username, h) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "not allowed",
		})
		return
	}
	user := c.Param("username")
	postContainer(user, h, c)
}

func (h *Handler) deleteUserContainer(c *gin.Context) {
	username := getUsername(c)
	if !isAdmin(username, h) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "not allowed",
		})
		return
	}
	user := c.Param("username")
	deleteContainer(user, h, c)
}
