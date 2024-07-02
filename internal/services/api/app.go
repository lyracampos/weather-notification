package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"weather-notification/configs"
	"weather-notification/internal/domain/usecases"
	"weather-notification/internal/gateways/broker/rabbitmq"
	"weather-notification/internal/gateways/database/postgres"
	api "weather-notification/internal/gateways/http"
	"weather-notification/internal/services/api/handlers"
	"weather-notification/internal/services/websocket"

	"github.com/go-openapi/runtime/middleware"
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

	// initialize client dependences
	databaseClient, err := postgres.NewClient(sugar, config)
	if err != nil {
		sugar.Fatalf("failed to initialize postgres client: %w", err)
	}
	defer databaseClient.Close()
	userDatabase := postgres.NewUserDatabase(databaseClient)

	brokerClient, err := rabbitmq.NewClient(sugar, config)
	if err != nil {
		sugar.Fatalf("failed to initialize rabbitmq client: %w", err)
	}
	defer brokerClient.Close()
	publishBroker := rabbitmq.NewPublisher(brokerClient)

	weatherAPI := api.NewWeatherAPI(sugar, config)

	registerUsecase := usecases.NewRegisterUseCase(userDatabase, weatherAPI)
	unsubscribeUseCase := usecases.NewUnsubscribeUseUseCase(userDatabase)
	queueNotificationsUseCase := usecases.NewEnqueueNotificationsUseCase(sugar, publishBroker)

	router := mux.NewRouter()

	healthHandler := handlers.NewHealthHandler(sugar)
	usersHandler := handlers.NewUserHandler(sugar, registerUsecase, unsubscribeUseCase)
	notificationsHandler := handlers.NewNotificationHandler(sugar, queueNotificationsUseCase)

	websocketHandler := websocket.NewWebSocketHandler(sugar)

	router.HandleFunc("/health", healthHandler.Health).Methods(http.MethodGet)

	router.HandleFunc("/users", usersHandler.Register).Methods(http.MethodPost)
	router.HandleFunc("/users/{email}/unsubscribe", usersHandler.Unsubscribe).Methods(http.MethodPut)
	router.HandleFunc("/notifications", notificationsHandler.Notify).Methods(http.MethodPost)

	router.HandleFunc("/ws/connect", websocketHandler.Connect)
	router.HandleFunc("/ws/clients", websocketHandler.Clients).Methods(http.MethodGet)
	router.HandleFunc("/ws/notify", websocketHandler.NotifyUser).Methods(http.MethodPost)

	router.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
	opts := middleware.SwaggerUIOpts{SpecURL: "swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	router.Handle("/docs", sh)

	http.Handle("/", router)

	address := fmt.Sprintf("%s:%d", config.API.Host, config.API.Port)
	server := &http.Server{
		Addr:         address,
		Handler:      router,
		WriteTimeout: time.Second * time.Duration(config.API.WriteTimeout),
		ReadTimeout:  time.Second * time.Duration(config.API.ReadTimeout),
		IdleTimeout:  time.Second * time.Duration(config.API.IdleTimeout),
	}

	go func() {
		log.Printf("running API HTTP server at: %s", address)

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
