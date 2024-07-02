package main

import (
	"errors"
	"flag"
	"fmt"
	"weather-notification/configs"
	"weather-notification/internal/services/api"
	"weather-notification/internal/services/worker"
)

const (
	defaultConfigFilePath = "../configs/config.yaml"
	apiEntrypoint         = "api"
	workerEntrypoint      = "worker"
	websocketWorkerType   = "websocket"
	emailWorkerType       = "email"
	smsWorkerType         = "sms"
	pushWorkerType        = "push"
)

var errInvalidAppEntrypoint = errors.New("invalid entrypoint, must be one of [api, worker]")

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var configFilePath string
	var appEntrypoint string
	var workerType string

	flag.StringVar(&configFilePath, "c", defaultConfigFilePath, "File path with app configs file.")
	flag.StringVar(&appEntrypoint, "e", apiEntrypoint, "Entrypoint to define which application will be started. [api, worker]")
	flag.StringVar(&workerType, "t", websocketWorkerType, "Type to define which worker will be started. [websocket, email, sms, push]")
	flag.Parse()

	config, err := configs.NewConfig(configFilePath)
	if err != nil {
		return fmt.Errorf("failed to build config: %w", err)
	}

	switch appEntrypoint {
	case apiEntrypoint:
		api.Run(config)
	case workerEntrypoint:
		worker.Run(config, workerType)
	default:
		return errInvalidAppEntrypoint
	}

	return nil
}
