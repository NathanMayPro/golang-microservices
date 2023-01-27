package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"golang-microservices/customlog" // import loggers package
)

type Document struct {
	Filename  string `json:"filename"`
	Content   string `json:"content"`
	Extension string `json:"extension"`
	Uid       string `json:"uid"`
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	//Data    interface{} `json:"data"`
}

// decoding b64 encoded document
func decodeb64(Content string) []byte {
	// decode base64 string
	decoded, err := base64.StdEncoding.DecodeString(Content)
	must(err, "decode base64 string")

	// check if decoded is empty
	if decoded == nil || len(decoded) == 0 {
		logAndKill(err, "decode base64 string")
	}
	return decoded
}

func homepage(w http.ResponseWriter, r *http.Request) {
	log.Printf("endpoint:homepage hit %v", r.RemoteAddr)
}

func save_file(decoded []byte, filename string, extension string, path string) error {
	// create file
	file, err := os.Create(path + filename + "." + extension)
	must(err, "create file")

	// write to file
	_, err = file.Write(decoded)
	must(err, "write to file")

	// close file
	err = file.Close()
	must(err, "close file")

	return err
}

func save(w http.ResponseWriter, r *http.Request) (Response, error) {
	// log activity
	log.Printf("save endpoint hit %v", r.RemoteAddr)

	// check if there is a path as query parameter
	path := r.URL.Query().Get("path") // http://localhost:8080/save?path=/home/username

	// if path is not empty
	if path == "" {
		log.Printf("endpoint:homepage hit %v", r.RemoteAddr)
		return Response{
			Status:  400,
			Message: "You must provide a path.",
		}, nil // to stop the execution of the function
	}

	// Declare a new Document struct.
	var d Document

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		return Response{
			Status:  400,
			Message: "Invalid request payload",
		}, nil
	}
	// decode base64 string
	decoded, err := base64.StdEncoding.DecodeString(d.Content)
	if err != nil {
		return Response{
			Status:  400,
			Message: "Invalid request payload",
		}, nil // to stop the execution of the function
	}

	// save file
	err = save_file(decoded, d.Filename, d.Extension, path)

	if err != nil {
		return Response{
			Status:  400,
			Message: "Internal file save error",
		}, nil // to stop the execution of the function
	}

	log.Printf("save endpoint done %v", r.RemoteAddr)
	return Response{
		Status:  200,
		Message: "success",
	}, nil

}

func retrieve(w http.ResponseWriter, r *http.Request) (Response, error) {
	// log activity
	log.Printf("retrieve endpoint hit %v", r.RemoteAddr)

	// check if there is a path as query parameter
	path := r.URL.Query().Get("path") // http://localhost:8080/retrieve?path=/home/username
	// if path is empty
	if path == "" || len(path) == 0 {
		return Response{
			Status:  400,
			Message: "You must provide a path.",
		}, nil // to stop the execution of the function
	}

	// check if file exists
	fileInfo, error := os.Stat(path)
	if error != nil {
		return Response{
			Status:  400,
			Message: "File does not exist.",
		}, nil // to stop the execution of the function
	}

	// check if file is a directory
	if fileInfo.IsDir() {
		return Response{
			Status:  400,
			Message: "Filepath is a directory.",
		}, nil // to stop the execution of the function

	}

	// open file using READ & WRITE permission
	http.ServeFile(w, r, path)

	log.Printf("retrieve endpoint done %v", r.RemoteAddr)
	return Response{
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

type customHandler func(w http.ResponseWriter, r *http.Request) (Response, error)

func responseMiddleware(next customHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Call the custom handler
		response, err := next(w, r)
		// check the response
		if err != nil {
			json.NewEncoder(w).Encode(Response{
				Status:  500,
				Message: "internal server error",
			})
			customlog.NewLog("Internal server error.", r)
		}
		//in case of success return the respons of the handler
		json.NewEncoder(w).Encode(Response{
			Status:  response.Status,
			Message: response.Message,
		})
		customlog.NewLog(response.Message, r)

	})
}

func main() {
	// // set up the log file
	// file, err := os.OpenFile("./logging/api.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)

	// must(err, "open log file")

	// set the output of the logs to be the file
	// log.SetOutput(file)

	// defer file.Close()

	mux := http.NewServeMux()
	// mux.Handle("/", responseMiddleware(homepage))
	// mux.Handle("/save", responseMiddleware(http.HandlerFunc(save)))
	mux.Handle("/retrieve", responseMiddleware(retrieve))

	errServer := http.ListenAndServe(":8080", mux)
	must(errServer, "start server")
}
