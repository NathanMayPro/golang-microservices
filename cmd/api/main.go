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


//to import an interface from another package
// import (
// 	"golang-microservices/internal/handlers"
// )
//
// type Document struct {
// 	Filename  string `json:"filename"`
// 	Content   string `json:"content"`
// 	Extension string `json:"extension"`
// 	Uid       string `json:"uid"`




func homepage(w http.ResponseWriter, r *http.Request) {
	log.Printf("endpoint:homepage hit %v", r.RemoteAddr)
}




func retrieve(w http.ResponseWriter, r *http.Request) (datamodel.Response, error) {
	// log activity
	log.Printf("retrieve endpoint hit %v", r.RemoteAddr)

	// check if there is a path as query parameter
	path := r.URL.Query().Get("path") // http://localhost:8080/retrieve?path=/home/username
	// if path is empty
	if path == "" || len(path) == 0 {
		return datamodel.Response{
			Status:  400,
			Message: "You must provide a path.",
		}, nil // to stop the execution of the function
	}

	// check if file exists
	fileInfo, error := os.Stat(path)
	if error != nil {
		return datamodel.Response{
			Status:  400,
			Message: "File does not exist.",
		}, nil // to stop the execution of the function
	}

	// check if file is a directory
	if fileInfo.IsDir() {
		return datamodel.Response{
			Status:  400,
			Message: "Filepath is a directory.",
		}, nil // to stop the execution of the function

	}

	// open file using READ & WRITE permission
	http.ServeFile(w, r, path)

	log.Printf("retrieve endpoint done %v", r.RemoteAddr)
	return datamodel.Response{
		Status:  200,
		Message: "success",
	}, nil

}

func logAndKill(err error, msg string) {
	log.Printf("failed to %s: %v", msg, err)
}

func must(err error, msg string) {
	if err != nil {
		logAndKill(err, msg)
	}
}

type customHandler func(w http.ResponseWriter, r *http.Request) (datamodel.Response, error)

func ResponseMiddleware(next customHandler) http.Handler {
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
	// set up the log file
	file, err := os.OpenFile("./logging/api.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)

	must(err, "open log file")

	//set the output of the logs to be the file
	log.SetOutput(file)

	defer file.Close()

	mux := http.NewServeMux()
	// mux.Handle("/", datamodel.ResponseMiddleware(homepage))
	//mux.Handle("/save", datamodel.ResponseMiddleware(handlers.HandlerSave))
	mux.Handle("/retrieve", ResponseMiddleware(handlers.HandlerRetrieve))

	errServer := http.ListenAndServe(":8080", mux)
	must(errServer, "start server")
}
