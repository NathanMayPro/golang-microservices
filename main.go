package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
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
	if err != nil {
		log.Fatal(err)
	}
	if decoded == nil || len(decoded) == 0 {
		log.Fatal("Something went wrong while decoding base64 string.")
	}
	return decoded
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Endpoint Hit")
}

func save_file(decoded []byte, filename string, extension string, path string) error {
	// create file
	file, err := os.Create(path + filename + "." + extension)
	if err != nil {
		log.Fatal(err)
	}

	// write to file
	_, err = file.Write(decoded)
	if err != nil {
		log.Fatal(err)
	}

	// close file
	err = file.Close()
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func save(w http.ResponseWriter, r *http.Request) {

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
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Log: save endpoint hit %v", path)

}

func retrieve(w http.ResponseWriter, r *http.Request) {

	fmt.Println("retrieve endpoint hit")
	// check if there is a path as query parameter
	path := r.URL.Query().Get("path") // http://localhost:8080/retrieve?path=/home/username

	// if path is not empty
	if path == "" {
		log.Fatal("path is empty.")
	}

	// open file using READ & WRITE permission
	http.ServeFile(w, r, path)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homepage)
	mux.HandleFunc("/retrieve", retrieve)
	mux.HandleFunc("/save", save)

	// log activity
	fmt.Printf("Server started on port 8080\n")

	//mux.HandleFunc("/person/list", personList)

	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
