package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/williabk198/timeclock/internal/services/admin/endpoints"
	"github.com/williabk198/timeclock/internal/services/admin/transport"
)

func main() {
	// Get the database connection info from the environment, and attempt to open a connect to the db using that info.
	dbUri := os.Getenv("DB_URI")
	dbSession, err := sql.Open("postgres", dbUri)
	if err != nil {
		log.Fatalf("failed to parse database connection data: %v", err)
	}

	// Initialize the route handler(s) and the server that will handle incoming requests.
	adminEndpointHandlers := endpoints.NewAdminEndpointHandlers(dbSession)
	server := http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler: transport.NewHttpHandler(adminEndpointHandlers),
	}

	// Start the server on another thread, and send any errors that may occur to a channel so that it can be handled later, if need be.
	log.Println("Starting server at", server.Addr)
	errChan := make(chan error)
	go func() {
		errChan <- server.ListenAndServe()
	}()

	// Create and register a channel that listens for a terminate or interrupt signal from the operating system
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	// Blocking select statement. This waits for a message to come in on either errChan, or sigChan.
	select {
	case err := <-errChan:
		// If we get something on the errChan, then an error occurred while starting the server.
		// So, just log it and return immediately since the server never started.
		log.Println(err)
		return
	case sig := <-sigChan:
		// If we get a terminate or interrupt signal from the OS, log it and the attempt to gracefully shutdown the server
		log.Printf("Received %q signal. Shutting down", sig)
	}

	// Create a timeout context for server.Shutdown. Otherwise, it will wait indefinitely for connections to close
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Try to allow the server to wrap up any operations on active connections before shutting down.
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Println(err)
	}
}
