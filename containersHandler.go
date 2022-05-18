package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getMyContainer(c *gin.Context) {
	username := c.GetHeader("X-Forwarded-User")
	getContainer(username, h, c)
}

func (h *Handler) postMyContainer(c *gin.Context) {
	username := c.GetHeader("X-Forwarded-User")
	postContainer(username, h, c)
}

func (h *Handler) deleteMyContainer(c *gin.Context) {
	username := c.GetHeader("X-Forwarded-User")
	deleteContainer(username, h, c)
}

func (h *Handler) getUserContainer(c *gin.Context) {
	username := c.GetHeader("X-Forwarded-User")
	if !isAdmin(username, h) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "not allowed",
		})
		return
	}
	user := c.Param("user")
	getContainer(user, h, c)

}

func (h *Handler) postUserContainer(c *gin.Context) {
	username := c.GetHeader("X-Forwarded-User")
	if !isAdmin(username, h) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "not allowed",
		})
		return
	}
	user := c.Param("user")
	postContainer(user, h, c)
}

func (h *Handler) deleteUserContainer(c *gin.Context) {
	username := c.GetHeader("X-Forwarded-User")
	if !isAdmin(username, h) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "not allowed",
		})
		return
	}
	user := c.Param("user")
	deleteContainer(user, h, c)
}
