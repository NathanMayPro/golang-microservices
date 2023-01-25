package main

import (
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	return decoded
}

func read_csv_from_bytes(decoded []byte) ([][]string, error) {
	// read csv from b64
	csvReader := csv.NewReader(strings.NewReader(string(decoded)))

	// read all the records
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return records, err
}

// curl -X GET "http://localhost:8080/person/create" -d '{"name":"John", "age":30}'
func converter(w http.ResponseWriter, r *http.Request) {
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

	// read csv from b64
	records, err := read_csv_from_bytes(decoded)

	if err != nil {
		log.Fatal(err)
	}

	// Do something with the Person struct...
	w.Write([]byte(fmt.Sprintf(records[0][0])))
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage Endpoint Hit")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homepage)
	mux.HandleFunc("/converter", converter)

	//mux.HandleFunc("/person/list", personList)

	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
