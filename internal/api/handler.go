package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/matheusmazzoni/url-shortener/internal/config"
	"github.com/matheusmazzoni/url-shortener/internal/shortener"
	"github.com/matheusmazzoni/url-shortener/internal/storage"
	"github.com/rs/zerolog"
)

const maxCollisionRetries = 10

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

// Handler groups the API handlers' dependencies, such as database access and configuration.
type Handler struct {
	Store  storage.Storage
	Config *config.Config
}

func (h *Handler) ShortenURLHandler(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.Ctx(r.Context())

	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error().Err(err).Msg("Failed to decode request body")
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		logger.Warn().Msg("Received shorten request with empty URL")
		http.Error(w, `{"error": "URL cannot be empty"}`, http.StatusBadRequest)
		return
	}

	// Check if the long URL has already been shortened.
	existingKey, err := h.Store.GetShortKey(req.URL)
	if err == nil {
		response := ShortenResponse{
			ShortURL: h.Config.AppBaseURL + "/" + existingKey,
		}
		logger.Info().Str("original_url", req.URL).Str("short_key", existingKey).Msg("URL already shortened. Returning existing key.")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // 200 OK because we are returning existing data, not creating.
		json.NewEncoder(w).Encode(response)
		return
	}
	// Only proceed if the error is specifically 'sql.ErrNoRows'. Any other error is a server problem.
	if err != sql.ErrNoRows {
		logger.Error().Err(err).Msg("Database error when checking for existing URL")
		http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
		return
	}

	// --- Collision handling loop ---
	var newShortKey string
	var success bool

	for i := 0; i < maxCollisionRetries; i++ {
		candidateKey := shortener.GenerateShortKey()

		exists, err := h.Store.KeyExists(candidateKey)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to check if key exists in database")
			http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
			return
		}

		if !exists {
			err = h.Store.SaveURL(candidateKey, req.URL)
			if err != nil {
				// This could happen in a rare race condition. Logging it is important for monitoring.
				logger.Error().Err(err).Str("key", candidateKey).Msg("Failed to save URL, possibly a race condition")
				http.Error(w, `{"error": "Internal server error"}`, http.StatusInternalServerError)
				return
			}

			newShortKey = candidateKey
			success = true
			logger.Info().Int("attempts", i+1).Str("key", newShortKey).Msg("Unique key generated and saved successfully")
			break // Successfully saved, exit the loop.
		}
		logger.Warn().Str("key", candidateKey).Int("attempt", i+1).Msg("Key collision detected. Generating new key.")
	}

	if !success {
		logger.Error().Int("retries", maxCollisionRetries).Msg("Failed to generate unique key after max retries")
		http.Error(w, `{"error": "Service unavailable, please try again later"}`, http.StatusServiceUnavailable)
		return
	}

	response := ShortenResponse{
		ShortURL: h.Config.AppBaseURL + "/" + newShortKey,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created because a new resource was created.
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	logger := zerolog.Ctx(r.Context())

	shortKey := r.PathValue("shortKey")

	originalURL, err := h.Store.GetURL(shortKey)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Warn().Str("short_key", shortKey).Msg("Short key not found in database")
			// A simple text response is fine for a 404 on a redirect attempt.
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		logger.Error().Err(err).Str("short_key", shortKey).Msg("Error fetching URL from database")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
}
