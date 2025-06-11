package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/matheusmazzoni/url-shortener/internal/api"
	"github.com/matheusmazzoni/url-shortener/internal/config"
	"github.com/matheusmazzoni/url-shortener/internal/storage"
	"github.com/rs/zerolog"
)

func main() {

	// TODO: Create a logging.go file for more advanced logger configurations (e.g., dev vs. prod).
	// For now, we use a simple logger that writes to standard output.
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	// Load application configuration from environment variables.
	cfg, err := config.New()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize the database connection.
	db, err := storage.NewSQLiteStore(cfg.DBPath)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to initialize database")
	}

	// Set up the router and middleware chain.
	router := api.NewRouter(cfg, db, &logger)
	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: *router,
	}

	// Start the server in a new goroutine so it doesn't block.
	go func() {
		logger.Info().Str("address", server.Addr).Msg("Server started")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Failed to start the server")
		}
	}()

	// Block until a shutdown signal is received.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	sig := <-stop
	logger.Info().Str("signal", sig.String()).Msg("Signal received. Shutting down server gracefully.")

	// Create a context with a timeout to allow ongoing requests to finish.
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelShutdown()

	logger.Info().Msg("Stopping server")

	// Attempt to gracefully shut down the server.
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatal().Err(err).Msg("Error during server shutdown")
	}

	// Now that the server is shut down, it's safe to close the database connection.
	logger.Info().Msg("Closing database connection...")
	if err := db.Close(); err != nil {
		logger.Error().Err(err).Msg("Error closing the database")
	}

	logger.Info().Msg("Server shut down successfully.")
}
