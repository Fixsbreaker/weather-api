package service

import "testing"

func TestClothingRecommendation(t *testing.T) {
	tests := []struct {
		temp     float64
		expected string
	}{
		{-10, "Тёплая одежда"},
		{0, "Тёплая одежда"},
		{9.9, "Тёплая одежда"},
		{10, "Куртка"},
		{15, "Куртка"},
		{19.9, "Куртка"},
		{20, "Лёгкая одежда"},
		{35, "Лёгкая одежда"},
	}

	for _, tt := range tests {
		got := clothingRecommendation(tt.temp)
		if got != tt.expected {
			t.Errorf("clothingRecommendation(%.1f) = %q, want %q", tt.temp, got, tt.expected)
		}
	}
}

func TestMapWeatherCode(t *testing.T) {
	tests := []struct {
		code     int
		expected string
	}{
		{0, "Ясно"},
		{1, "Переменная облачность"},
		{2, "Переменная облачность"},
		{3, "Переменная облачность"},
		{45, "Туман"},
		{48, "Туман"},
		{51, "Морось"},
		{53, "Морось"},
		{55, "Морось"},
		{61, "Дождь"},
		{63, "Дождь"},
		{65, "Дождь"},
		{71, "Снег"},
		{73, "Снег"},
		{75, "Снег"},
		{95, "Гроза"},
		{999, "Неизвестно"},
	}

	for _, tt := range tests {
		got := mapWeatherCode(tt.code)
		if got != tt.expected {
			t.Errorf("mapWeatherCode(%d) = %q, want %q", tt.code, got, tt.expected)
		}
	}
}
