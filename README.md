# weather-api 

This is a simple Weather API application built in Go using the `go-chi/chi` router. It fetches current weather data by interacting with external public APIs.

## Overview
Initially the project had a basic setup to fetch the weather using latitude and longitude coordinates via the Open-Meteo API. I expanded on this  by adding a Geocoding client and multiple new endpoints that allow users to search for weather by city name, country names, and get recommendations.

## Features I Added

1. **Geocoding Integration**:
   - Created a new `GeocodingClient` to interact with `geocoding-api.open-meteo.com` and `countriesnow.space`. 
   - This allows the application to resolve a city name into geographical coordinates (latitude and longitude) and to list all cities inside a specific country.

2. **Search Weather by City** (`GET /weather/{city}`):
   - Users can now get the weather by simply typing the city name instead of coordinates.
   - Using the Geocoding API, the app first finds the coordinates of the city and then fetches the current weather.
   - **Clothing Recommendations**: I implemented a feature that analyzes the fetched temperature and provides simple clothing suggestions (e.g., "Тёплая одежда" for cold days, "Лёгкая одежда" for warm days).

3. **Search Weather by Country** (`GET /weather/country/{country}`):
   - Retrieves all cities for a given country and fetches the weather for the first few cities.

4. **Top Warmest Cities** (`GET /weather/country/{country}/top`):
   - Fetches weather for a country's cities and sorts the results to return the top 3 warmest locations right now.

5. **Project Restructuring**:
   - I moved `main.go` into `cmd/app/main.go` for better project organization, following Go standard directory layout.

## Running the Application

To start the server, run:
```bash
go run cmd/app/main.go
```
The server will start on port `8080`.

## Endpoints

- `GET /health` - Check if API is running
- `GET /api/weather?lat={lat}&lon={lon}` - Get weather by coordinates
- `GET /weather/{city}` - Get weather by city name
- `GET /weather/country/{country}` - Get weather for several cities in a country
- `GET /weather/country/{country}/top` - Get top 3 warmest cities in a country
