package api

import (
	"net/http"

	"github.com/matheusmazzoni/url-shortener/internal/config"
	"github.com/matheusmazzoni/url-shortener/internal/storage"
	"github.com/rs/zerolog"
)

// NewRouter sets up the application's routing and middleware chain.
// It initializes the API handlers and applies all necessary middlewares
// in the correct order, returning a single, final http.Handler.
func NewRouter(cfg *config.Config, store storage.Storage, logger *zerolog.Logger) *http.Handler {
	// Initialize the API handler struct with its dependencies.
	handler := &Handler{
		Store:  store,
		Config: cfg,
	}

	// Create a new ServeMux to register the routes.
	mux := http.NewServeMux()

	// Register application routes.
	mux.HandleFunc("POST /shorten", handler.ShortenURLHandler)
	mux.HandleFunc("GET /{shortKey}", handler.RedirectHandler)

	// Create the middleware chain. The order is important: outer -> inner.
	// 1. ContextualLogMiddleware adds request-scoped logging.
	handlerChain := Chain(mux,
		ContextualLogMiddleware(logger),
		// Other middlewares would be added here...
	)

	return &handlerChain
}
