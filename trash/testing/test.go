package main

// This file will be used to made a request to the api and get the response

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

// function to test if string is a valid string
func isNotEmpty(x interface{}) bool {
	if x == "" {
		return false
	}
	return true
}

// function to test if contains spaces
func containsSpaces(str string) bool {
	for _, char := range str {
		if char == ' ' {
			return true
		}
	}
	return false
}

func checkPort(port string) (bool, string) {
	port_cast, err := strconv.Atoi(port)

	if err != nil {
		return false, "Port argument is not a number"
	}

	if port_cast < 0 || port_cast > 65535 {
		return false, "Port arguments is not a valid port"
	}
	return true, ""
}

func checkEndpoint(endpoint string) (bool, string) {
	if !isNotEmpty(endpoint) {
		return false, "Endpoint argument is empty"
	}
	if containsSpaces(endpoint) {
		return false, "Endpoint argument contains spaces"
	}
	return true, ""
}

// function to parse the arguments
func parseArgs(args []string) (int, string) {
	// check if len of args is 0
	if len(args) == 0 {
		fmt.Println("Please provide a command")
		return 0, ""
	}
	port := args[0]
	if _, error_message := checkPort(port); error_message != "" {
		fmt.Println(error_message)
		return 0, ""
	}
	endpoint := args[1]
	if _, error_message := checkEndpoint(endpoint); error_message != "" {
		fmt.Println(error_message)
		return 0, ""
	}

	// cast port to int
	port_cast, _ := strconv.Atoi(port) // error is already checked

	// return port, endpoint
	return port_cast, endpoint
}

func main() {
	// parse arguments
	args := os.Args[1:]
	// check if len of args is 0
	if len(args) == 0 {
		fmt.Println("Please provide a command")
		return
	}

	port := args[0]
	endpoint := args[1]

	// create a new request to the api
	url := fmt.Sprintf("http://localhost:%v%v", port, endpoint)
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s	", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(data))
	}
}
