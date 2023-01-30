package customlog

import (
	"log"
	"net/http"
)

func NewLog(logging_message string, r *http.Request) {
	log.Printf(logging_message + " " + r.RemoteAddr)
}