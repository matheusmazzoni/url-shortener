package api

import (
	"net/http"
	"time" // Import needed for duration

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// Middleware is a function type that takes an http.Handler and returns another one.
type Middleware func(http.Handler) http.Handler

// Chain applies a list of middlewares to an http.Handler.
// The middlewares are applied in the order they are provided in the call.
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	// To apply in the correct order (the first in the list is the outermost),
	// we iterate over the slice backwards.
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// --- ResponseWriter Interceptor to capture status code ---

// responseWriterInterceptor is a wrapper for http.ResponseWriter that allows
// capturing the HTTP status code written to the response.
type responseWriterInterceptor struct {
	http.ResponseWriter
	statusCode int
}

// newResponseWriterInterceptor creates a new interceptor.
func newResponseWriterInterceptor(w http.ResponseWriter) *responseWriterInterceptor {
	// Default to 200 OK if WriteHeader is not called.
	return &responseWriterInterceptor{w, http.StatusOK}
}

// WriteHeader captures the status code and calls the original WriteHeader.
func (rwi *responseWriterInterceptor) WriteHeader(code int) {
	rwi.statusCode = code
	rwi.ResponseWriter.WriteHeader(code)
}

// --- Middlewares ---

// ContextualLogMiddleware injects a request-scoped logger into the context and
// logs a final summary event after the request has been handled.
func ContextualLogMiddleware(baseLogger *zerolog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestID := uuid.New().String()

			// Create a child logger with the request_id field.
			// This logger will be passed down in the context for handlers to use.
			ctxLogger := baseLogger.With().Str("request_id", requestID).Logger()
			ctx := ctxLogger.WithContext(r.Context())

			// Create the interceptor to capture the status code.
			rwi := newResponseWriterInterceptor(w)

			// Call the next handler in the chain with the new context and the interceptor.
			next.ServeHTTP(rwi, r.WithContext(ctx))

			// After the handler has finished, log the final summary event.
			ctxLogger.Info().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status_code", rwi.statusCode). // Log the captured status code.
				Dur("duration_ms", time.Since(start)).
				Msg("Request completed")
		})
	}
}
