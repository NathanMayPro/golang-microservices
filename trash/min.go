package main

import (
	"golang-microservices/customlog" // import loggers package

	"net/http"
)


func main() {
	// log activity
	customlog.NewLog("endpoint:homepage hit", &http.Request{
		RemoteAddr: "",
	})
}