package handlers

import (
	"log"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("register handler"))
	if err != nil {
		log.Printf("Health - write failed: %v", err)
	}
}
