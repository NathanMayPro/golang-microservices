package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"golang-microservices/internal/customlog" // import loggers package
	"golang-microservices/internal/handlers"
	"golang-microservices/internal/datamodel"
)

type customHandler func(w http.ResponseWriter, r *http.Request) (datamodel.Response, error)

func responseMiddleware(next customHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Call the custom handler
		response, err := next(w, r)
		// check the datamodel.Response
		if err != nil {
			json.NewEncoder(w).Encode(datamodel.Response{
				Status:  500,
				Message: "internal server error",
			})
			customlog.NewLog("Internal server error.", r)
		}
		//in case of success return the respons of the handler
		json.NewEncoder(w).Encode(datamodel.Response{
			Status:  response.Status,
			Message: response.Message,
		})
		customlog.NewLog(response.Message, r)

	})
}

func main() {
	/* SET UP THE LOG FILE */
	// set up the log file
	file, err := os.OpenFile("./logging/api.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	//set the output of the logs to be the file
	log.SetOutput(file)
	// close the file when the program ends
	defer file.Close()


	/* START THE SERVER */
	// set up the server
	mux := http.NewServeMux()
	
	mux.Handle("/", responseMiddleware(handlers.HandlerHomepage))
	mux.Handle("/save", responseMiddleware(handlers.HandlerSave))
	mux.Handle("/retrieve", responseMiddleware(handlers.HandlerRetrieve))

	errServer := http.ListenAndServe(":8080", mux)
	must(errServer, "start server")
}
