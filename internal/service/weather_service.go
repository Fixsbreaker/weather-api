package service

import (
	"context"
	"fmt"
	"sort"
)

// WeatherProvider is the interface for fetching weather by coordinates
type WeatherProvider interface {
	GetCurrentWeather(ctx context.Context, lat, lon float64) (*ProviderWeatherResponse, error)
}

// GeocodingProvider is the interface for geocoding operations
type GeocodingProvider interface {
	GetCityCoordinates(ctx context.Context, city string) (*GeoLocation, error)
	GetCitiesByCountry(ctx context.Context, country string) ([]string, error)
}

// ProviderWeatherResponse is the raw weather data returned by the external API
type ProviderWeatherResponse struct {
	Temperature float64
	WindSpeed   float64
	WeatherCode int
	Time        string
}

// GeoLocation holds coordinates and metadata for a city
type GeoLocation struct {
	Name      string
	Latitude  float64
	Longitude float64
	Country   string
}

// WeatherResult is the response for coordinate-based weather lookup
type WeatherResult struct {
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Temperature float64 `json:"temperature"`
	WindSpeed   float64 `json:"wind_speed"`
	WeatherCode int     `json:"weather_code"`
	Time        string  `json:"time"`
	Description string  `json:"description"`
}

// CityWeatherResult is the response for city-based weather lookup
type CityWeatherResult struct {
	City        string  `json:"city"`
	Country     string  `json:"country"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Temperature float64 `json:"temperature"`
	WindSpeed   float64 `json:"wind_speed"`
	WeatherCode int     `json:"weather_code"`
	Time        string  `json:"time"`
	Description string  `json:"description"`
	Clothing    string  `json:"clothing_recommendation"`
}

// WeatherService handles all weather business logic
type WeatherService struct {
	provider  WeatherProvider
	geocoding GeocodingProvider
}

// NewWeatherService creates a new WeatherService
func NewWeatherService(provider WeatherProvider, geocoding GeocodingProvider) *WeatherService {
	return &WeatherService{
		provider:  provider,
		geocoding: geocoding,
	}
}

// GetWeather returns weather data for a given latitude and longitude
func (s *WeatherService) GetWeather(ctx context.Context, lat, lon float64) (*WeatherResult, error) {
	resp, err := s.provider.GetCurrentWeather(ctx, lat, lon)
	if err != nil {
		return nil, fmt.Errorf("get weather from provider: %w", err)
	}

	return &WeatherResult{
		Latitude:    lat,
		Longitude:   lon,
		Temperature: resp.Temperature,
		WindSpeed:   resp.WindSpeed,
		WeatherCode: resp.WeatherCode,
		Time:        resp.Time,
		Description: mapWeatherCode(resp.WeatherCode),
	}, nil
}

// GetWeatherByCity returns weather for a city by its name.
func (s *WeatherService) GetWeatherByCity(ctx context.Context, city string) (*CityWeatherResult, error) {
	loc, err := s.geocoding.GetCityCoordinates(ctx, city)
	if err != nil {
		return nil, fmt.Errorf("get city coordinates: %w", err)
	}

	resp, err := s.provider.GetCurrentWeather(ctx, loc.Latitude, loc.Longitude)
	if err != nil {
		return nil, fmt.Errorf("get weather: %w", err)
	}

	return &CityWeatherResult{
		City:        loc.Name,
		Country:     loc.Country,
		Latitude:    loc.Latitude,
		Longitude:   loc.Longitude,
		Temperature: resp.Temperature,
		WindSpeed:   resp.WindSpeed,
		WeatherCode: resp.WeatherCode,
		Time:        resp.Time,
		Description: mapWeatherCode(resp.WeatherCode),
		Clothing:    clothingRecommendation(resp.Temperature),
	}, nil
}

// GetWeatherByCountry returns weather for cities in the given country
// it limits results for the first 7 cities to avoid excessive API calls
func (s *WeatherService) GetWeatherByCountry(ctx context.Context, country string) ([]CityWeatherResult, error) {
	cities, err := s.geocoding.GetCitiesByCountry(ctx, country)
	if err != nil {
		return nil, fmt.Errorf("get cities by country: %w", err)
	}

	const maxCities = 7
	if len(cities) > maxCities {
		cities = cities[:maxCities]
	}

	var results []CityWeatherResult
	for _, city := range cities {
		w, err := s.GetWeatherByCity(ctx, city)
		if err != nil {
			// skip cities that could not be resolved and continue with the rest
			continue
		}
		results = append(results, *w)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("failed to fetch weather for any city in the country")
	}

	return results, nil
}

// GetTopCitiesByCountry returns the top 3 warmest cities in a country
func (s *WeatherService) GetTopCitiesByCountry(ctx context.Context, country string) ([]CityWeatherResult, error) {
	results, err := s.GetWeatherByCountry(ctx, country)
	if err != nil {
		return nil, err
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Temperature > results[j].Temperature
	})

	if len(results) > 3 {
		results = results[:3]
	}

	return results, nil
}

// clothingRecommendation returns a clothing suggestion based on temperature
func clothingRecommendation(temp float64) string {
	switch {
	case temp < 10:
		return "Тёплая одежда"
	case temp < 20:
		return "Куртка"
	default:
		return "Лёгкая одежда"
	}
}

// mapWeatherCode converts an Open-Meteo weather code to a human-readable description
func mapWeatherCode(code int) string {
	switch code {
	case 0:
		return "Ясно"
	case 1, 2, 3:
		return "Переменная облачность"
	case 45, 48:
		return "Туман"
	case 51, 53, 55:
		return "Морось"
	case 61, 63, 65:
		return "Дождь"
	case 71, 73, 75:
		return "Снег"
	case 95:
		return "Гроза"
	default:
		return "Неизвестно"
	}
}
