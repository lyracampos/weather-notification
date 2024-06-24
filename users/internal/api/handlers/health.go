package handlers

import (
	"log"
	"net/http"
)

func Health(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("user API is running"))
	if err != nil {
		log.Printf("Health - write failed: %v", err)
	}
}
