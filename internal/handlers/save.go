package handlers


import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"golang-microservices/internal/datamodel"
)

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

func HandlerSave(w http.ResponseWriter, r *http.Request) (datamodel.Response, error) {
	// log activity
	log.Printf("save endpoint hit %v", r.RemoteAddr)

	// check if there is a path as query parameter
	path := r.URL.Query().Get("path") // http://localhost:8080/save?path=/home/username

	// if path is not empty
	if path == "" {
		log.Printf("endpoint:homepage hit %v", r.RemoteAddr)
		return datamodel.Response{
			Status:  400,
			Message: "You must provide a path.",
		}, nil // to stop the execution of the function
	}

	// Declare a new Document struct.
	var d datamodel.Document

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		return datamodel.Response{
			Status:  400,
			Message: "Invalid request payload",
		}, nil
	}
	// decode base64 string
	decoded, err := base64.StdEncoding.DecodeString(d.Content)
	if err != nil {
		return datamodel.Response{
			Status:  400,
			Message: "Invalid request payload",
		}, nil // to stop the execution of the function
	}

	// save file
	err = save_file(decoded, d.Filename, d.Extension, path)

	if err != nil {
		return datamodel.Response{
			Status:  400,
			Message: "Internal file save error",
		}, nil // to stop the execution of the function
	}

	log.Printf("save endpoint done %v", r.RemoteAddr)
	return datamodel.Response{
		Status:  200,
		Message: "success",
	}, nil

}
