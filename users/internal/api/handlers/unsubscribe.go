package handlers

import (
	"log"
	"net/http"
)

func Unsubscribe(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("unsubscribe handler"))
	if err != nil {
		log.Printf("Health - write failed: %v", err)
	}
}
