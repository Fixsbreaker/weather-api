package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"weather-api/internal/service"
)

// Service defines all operations the handler depends on
type Service interface {
	GetWeather(ctx context.Context, lat, lon float64) (*service.WeatherResult, error)
	GetWeatherByCity(ctx context.Context, city string) (*service.CityWeatherResult, error)
	GetWeatherByCountry(ctx context.Context, country string) ([]service.CityWeatherResult, error)
	GetTopCitiesByCountry(ctx context.Context, country string) ([]service.CityWeatherResult, error)
}

// WeatherHandler holds the service and handles HTTP requests
type WeatherHandler struct {
	service Service
}

// NewWeatherHandler creates a new WeatherHandler
func NewWeatherHandler(service Service) *WeatherHandler {
	return &WeatherHandler{service: service}
}

// ErrorResponse is returned when an error occurs
type ErrorResponse struct {
	Error string `json:"error"`
}

// GetWeather handles GET
func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")

	if latStr == "" || lonStr == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "query params lat and lon are required"})
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid lat"})
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid lon"})
		return
	}

	result, err := h.service.GetWeather(r.Context(), lat, lon)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// GetWeatherByCity handles GET /weather/{city}
func (h *WeatherHandler) GetWeatherByCity(w http.ResponseWriter, r *http.Request) {
	city := chi.URLParam(r, "city")
	if city == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "city is required"})
		return
	}

	result, err := h.service.GetWeatherByCity(r.Context(), city)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// GetWeatherByCountry handles GET /weather/country/{country}
func (h *WeatherHandler) GetWeatherByCountry(w http.ResponseWriter, r *http.Request) {
	country := chi.URLParam(r, "country")
	if country == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "country is required"})
		return
	}

	results, err := h.service.GetWeatherByCountry(r.Context(), country)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, results)
}

// GetTopCitiesByCountry handles GET /weather/country/{country}/top
func (h *WeatherHandler) GetTopCitiesByCountry(w http.ResponseWriter, r *http.Request) {
	country := chi.URLParam(r, "country")
	if country == "" {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "country is required"})
		return
	}

	results, err := h.service.GetTopCitiesByCountry(r.Context(), country)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, results)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, `{"error":"failed to encode json"}`, http.StatusInternalServerError)
	}
}
