package websocket

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"weather-notification/configs"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func Run(config *configs.Config) {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Fatalf("failed to defer logger sync: %v", err)
		}
	}()

	sugar := logger.Sugar()

	router := mux.NewRouter()

	websocketHandler := NewWebSocketHandler(sugar)
	router.HandleFunc("/ws/connect", websocketHandler.Connect)
	router.HandleFunc("/ws/clients", websocketHandler.Clients).Methods(http.MethodGet)
	router.HandleFunc("/ws/notify", websocketHandler.NotifyUser).Methods(http.MethodPost)

	http.Handle("/", router)

	address := fmt.Sprintf("%s:%d", config.WebSocketClient.Host, config.WebSocketClient.Port)
	server := &http.Server{
		Addr:         address,
		Handler:      router,
		WriteTimeout: time.Second * time.Duration(config.WebSocketClient.WriteTimeout),
		ReadTimeout:  time.Second * time.Duration(config.WebSocketClient.ReadTimeout),
		IdleTimeout:  time.Second * time.Duration(config.WebSocketClient.IdleTimeout),
	}

	go func() {
		log.Printf("running Websocket client on: %s", address)

		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	if err := server.Shutdown(ctx); err != nil {
		log.Println("Error shutting down server: %w", err)
	}

	log.Println("shutting down")
}
