package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) login(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "http://padawan.kube.test.viarezo.fr/")
}
