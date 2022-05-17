package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func newHandler(db *gorm.DB) *Handler {
	return &Handler{
		db: db,
	}
}

var db *gorm.DB

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	handler := newHandler(db)
	handler.db.AutoMigrate(&Container{})

	router := gin.Default()
	router.POST("/container", handler.postContainer)
	router.GET("/container", handler.getContainer)
	router.DELETE("/container", handler.deleteContainer)
	router.Run("localhost:8080")
}

type Container struct {
	Username string `json:"username" gorm:"primary_key"`
	Image    string `json:"image"`
}

type ContainerRequest struct {
	Image    string `json:"image"`
	Password string `json:"password"`
}

func (h *Handler) getContainer(c *gin.Context) {
	username := c.GetHeader("X-Forwarded-User")
	var container Container

	if result := h.db.Where("username = ?", username).First(&container); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, container)
}

func (h *Handler) postContainer(c *gin.Context) {
	username := c.GetHeader("X-Forwarded-User")
	var request ContainerRequest

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	container := Container{
		Username: username,
		Image:    request.Image,
	}

	if result := h.db.Create(&container); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, container)
}

func (h *Handler) deleteContainer(c *gin.Context) {
	username := c.GetHeader("X-Forwarded-User")
	fmt.Println(username)

	if result := h.db.Where("username = ?", username).Delete(&Container{}); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	c.Status(http.StatusNoContent)
}
