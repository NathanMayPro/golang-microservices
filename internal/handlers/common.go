package handlers

import (
	"log"
)


func logAndKill(err error, msg string) {
	log.Printf("failed to %s: %v", msg, err)
}

func must(err error, msg string) {
	if err != nil {
		logAndKill(err, msg)
	}
}