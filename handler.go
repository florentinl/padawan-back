package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Handler struct {
	db        *gorm.DB
	client    *kubernetes.Clientset
	namespace string
	admins    []string
}

type Container struct {
	Username  string `json:"username" gorm:"primary_key"`
	ImageName string `json:"image_name"`
	Port      int    `json:"port"`
}

type ContainerRequest struct {
	ImageName string `json:"image_name"`
	Password  string `json:"password"`
}

type Image struct {
	ImageName  string `json:"image_name" gorm:"primary_key"`
	Repository string `json:"repository"`
	Shell      string `json:"shell"`
}

func getUsername(c *gin.Context) string {
	username := c.Request.Header.Get("X-Forwarded-User")
	return strings.ReplaceAll(strings.Split(username, "@")[0], ".", "-")
}

func newHandler() *Handler {
	dbuser := strings.Trim(os.Getenv("DB_USER"), "\n")
	dbpass := strings.Trim(os.Getenv("DB_PASSWORD"), "\n")
	dbname := strings.Trim(os.Getenv("DB_NAME"), "\n")
	dbhost := strings.Trim(os.Getenv("DB_HOST"), "\n")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s", dbuser, dbpass, dbhost, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Image{})
	db.AutoMigrate(&Container{})

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	b, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err != nil {
		panic(err.Error())
	}

	namespace := string(b)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	admins := strings.Split(os.Getenv("ADMINS"), ",")

	return &Handler{
		db:        db,
		client:    clientset,
		namespace: namespace,
		admins:    admins,
	}
}

func isAdmin(username string, h *Handler) bool {
	for _, admin := range h.admins {
		if username == admin {
			return true
		}
	}
	return false
}
