package handlers

import (
	"log"
	"net/http"
	"os"
	"golang-microservices/internal/datamodel"
)



func HandlerRetrieve(w http.ResponseWriter, r *http.Request) (datamodel.Response, error) {
	// log activity
	log.Printf("retrieve endpoint hit %v", r.RemoteAddr)

	//Response datamodel.Response

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

