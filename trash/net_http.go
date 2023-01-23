/*
This file will be used to create a simple app that will respond pong to ping

	and will also respond with the current time when asked for it.
*/

package main

import (
	"encoding/json"
	"net/http"
)

type album struct {
	Title  string
	Artist string
	Price  float64
}

var albums = []album{
	{Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
} // this is an example of a slice of data for album struct

// create a basic functions that will present the api
func home(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Welcome to the home page")
}

// create the function that handle the request GET /albums
func getAlbums(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(albums)
}

// first let's declare main function
func main() {
	// create a new serve with net/http package
	http.HandleFunc("/", home)
	// create a new handler function that will respond to the request
	http.HandleFunc("/albums", getAlbums)

	// start the server
	http.ListenAndServe(":8080", nil)
}
