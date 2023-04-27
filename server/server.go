package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/martian/log"
	"net/http"
)

type Car struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Color string `json:"color"`
}

var Cars = Car{
	ID: "1", Title: "BMW", Color: "Black",
}

func main() {
	router := gin.Default()
	router.GET("/cars", getCars)

	err := router.Run("localhost:8080")
	if err != nil {
		log.Infof("Failed to start you service")
	}
}

func getCars(c *gin.Context) {

	c.IndentedJSON(http.StatusOK, Cars)
}
