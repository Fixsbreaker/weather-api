package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"weather-api/internal/service"
)

// GeocodingClient fetches city coordinates and country city lists
type GeocodingClient struct {
	httpClient    *http.Client
	geocodingURL  string
	countryCities string
}

// NewGeocodingClient creates a new client
func NewGeocodingClient(httpClient *http.Client) *GeocodingClient {
	return &GeocodingClient{
		httpClient:    httpClient,
		geocodingURL:  "https://geocoding-api.open-meteo.com/v1/search",
		countryCities: "https://countriesnow.space/api/v0.1/countries/cities",
	}
}

type geocodingResponse struct {
	Results []struct {
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Country   string  `json:"country"`
	} `json:"results"`
}

// GetCityCoordinates returns coordinates and metadata for a city name
func (c *GeocodingClient) GetCityCoordinates(ctx context.Context, city string) (*service.GeoLocation, error) {
	u, err := url.Parse(c.geocodingURL)
	if err != nil {
		return nil, fmt.Errorf("parse geocoding url: %w", err)
	}

	q := u.Query()
	q.Set("name", city)
	q.Set("count", "1")
	q.Set("language", "en")
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call geocoding api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geocoding api returned status: %d", resp.StatusCode)
	}

	var result geocodingResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode geocoding response: %w", err)
	}

	if len(result.Results) == 0 {
		return nil, fmt.Errorf("city not found: %s", city)
	}

	r := result.Results[0]
	return &service.GeoLocation{
		Name:      r.Name,
		Latitude:  r.Latitude,
		Longitude: r.Longitude,
		Country:   r.Country,
	}, nil
}

type countryCitiesResponse struct {
	Error bool     `json:"error"`
	Data  []string `json:"data"`
}

// GetCitiesByCountry returns a list of city names for the given country
func (c *GeocodingClient) GetCitiesByCountry(ctx context.Context, country string) ([]string, error) {
	body, err := json.Marshal(map[string]string{"country": country})
	if err != nil {
		return nil, fmt.Errorf("marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.countryCities, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call cities api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cities api returned status: %d", resp.StatusCode)
	}

	var result countryCitiesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode cities response: %w", err)
	}

	if result.Error {
		return nil, fmt.Errorf("country not found: %s", country)
	}

	return result.Data, nil
}
