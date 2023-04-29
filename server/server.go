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

var Cars = []Car{
	{ID: "1", Title: "BMW", Color: "Black"},
	{ID: "2", Title: "Tesla", Color: "Red"},
}

func main() {
	router := gin.Default()
	router.GET("/cars/:id", getCarByID)
	router.POST("/cars", createCar)

	err := router.Run("localhost:8080")
	if err != nil {
		log.Infof("Failed to start you service")
	}
}

func createCar(c *gin.Context) {
	var newCar Car
	err := c.BindJSON(&newCar)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest,
			"Provided data is invalid")
		return
	}

	Cars = append(Cars, newCar)
	c.IndentedJSON(http.StatusCreated, newCar)
}

func getCarByID(c *gin.Context) {
	for _, car := range Cars {
		if car.ID == c.Param("id") {
			c.IndentedJSON(http.StatusOK, car)
			return
		}
	}

	// Return 404 Status Code and error message if no car was found.
	c.IndentedJSON(http.StatusNotFound,
		gin.H{"message": "Requested car is not found"})
}
