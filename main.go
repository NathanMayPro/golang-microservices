package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// We gonna create a function that present the api
func home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"output":  nil,
		"message": "Welcome to the home page",
		"status":  200,
	})
}

// create the function that handle the request GET /ping
func pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// creation of function tha handle arguments in the url
func giveNumber(c *gin.Context) {
	id := c.Param("id")
	response := ""
	// check if id is a number
	if id, err := strconv.Atoi(id); err == nil {
		response = fmt.Sprintf("The number is %d", id)
	} else {
		response = "The id is not a number"
	}
	c.JSON(http.StatusOK, gin.H{
		"output": response,
	})
}

func main() {
	r := gin.Default()

	// create a new handler function that will respond to the request
	r.GET("/", home)

	// create a new handler function that will respond to the request
	r.GET("/ping", pong)

	// function that give a number given in the url
	r.GET("/print/:id", giveNumber)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

	// To start the server, run the following command:
	// go run main.go
}
