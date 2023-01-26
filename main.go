package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Document struct {
	Filename  string `json:"filename"`
	Content   string `json:"content"`
	Extension string `json:"extension"`
	Uid       string `json:"uid"`
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

func save(w http.ResponseWriter, r *http.Request) {
	// log activity
	log.Printf("save endpoint hit %v", r.RemoteAddr)

	// check if there is a path as query parameter
	path := r.URL.Query().Get("path") // http://localhost:8080/save?path=/home/username

	// if path is not empty
	if path == "" {
		log.Fatal("path is empty.")
	}

	// Declare a new Document struct.
	var d Document

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// decode base64 string
	decoded, err := base64.StdEncoding.DecodeString(d.Content)
	if err != nil {
		log.Fatal(err)
	}

	// save file
	err = save_file(decoded, d.Filename, d.Extension, path)

	must(err, "save file")

	log.Printf("save endpoint done %v", r.RemoteAddr)
}

func retrieve(w http.ResponseWriter, r *http.Request) {
	// log activity
	log.Printf("retrieve endpoint hit %v", r.RemoteAddr)

	// check if there is a path as query parameter
	path := r.URL.Query().Get("path") // http://localhost:8080/retrieve?path=/home/username

	// if path is not empty
	if path == "" {
		logAndKill(nil, "path is empty.")
	}

	// check if file exists
	fileInfo, error := os.Stat(path)
	must(error, "check if file exists")

	// check if file is a directory
	if fileInfo.IsDir() {
		logAndKill(nil, "file is a directory")
	}

	// open file using READ & WRITE permission
	http.ServeFile(w, r, path)

	log.Printf("retrieve endpoint done %v", r.RemoteAddr)
}

func logAndKill(err error, msg string) {
	log.Printf("failed to %s: %v", msg, err)
	panic(err)

}

func must(err error, msg string) {
	if err != nil {
		logAndKill(err, msg)
	}
}

func main() {
	// set up the log file
	file, err := os.OpenFile("./logging/api.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)

	must(err, "open log file")

	// set the output of the logs to be the file
	log.SetOutput(file)

	defer file.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("/", homepage)
	mux.HandleFunc("/retrieve", retrieve)
	mux.HandleFunc("/save", save)

	// log activity
	log.Printf("Server started on port 8080\n")

	errServer := http.ListenAndServe(":8080", mux)
	must(errServer, "start server")
}
