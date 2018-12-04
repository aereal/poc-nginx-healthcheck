package web

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	HostPort        string
	ShutdownTimeout time.Duration
}

func Run(config *Config, handler http.Handler) error {
	server := &http.Server{
		Addr:    config.HostPort,
		Handler: handler,
	}
	go graceful(server, config.ShutdownTimeout)

	log.Printf("starting server (config:%#v) ...", config)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func graceful(server *http.Server, timeout time.Duration) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	sig := <-sigChan
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	log.Printf("shutting down server (%v) ...", sig)
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("failed to shutdown: %v", err)
	}
}

func RespondJSON(w http.ResponseWriter, body interface{}) error {
	w.Header().Set("content-type", "application/json")
	return json.NewEncoder(w).Encode(body)
}

type errorRes struct {
	Error string
}

func RespondErrorJSON(w http.ResponseWriter, statusCode int, err error) error {
	w.WriteHeader(statusCode)
	return RespondJSON(w, &errorRes{Error: err.Error()})
}
