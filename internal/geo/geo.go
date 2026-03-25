package geo

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"time"
)

const maxResponseSize = 1 * 1024 * 1024 // 1MB

const (
	earthRadiusKm = 6371.0
	issAltitudeKm = 408.0
	ipGeoURL      = "http://ip-api.com/json/"
)

type Location struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	City      string  `json:"city"`
	Country   string  `json:"country"`
}

type ipGeoResponse struct {
	Status  string  `json:"status"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	City    string  `json:"city"`
	Country string  `json:"country"`
}

func GetMyLocation() (*Location, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(ipGeoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get location: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geo API returned status %d", resp.StatusCode)
	}

	var data ipGeoResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxResponseSize)).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode location: %w", err)
	}

	if data.Status != "success" {
		return nil, fmt.Errorf("geo API returned non-success status")
	}

	return &Location{
		Latitude:  data.Lat,
		Longitude: data.Lon,
		City:      data.City,
		Country:   data.Country,
	}, nil
}

func DistanceToISS(userLat, userLon, issLat, issLon float64) float64 {
	groundDist := haversine(userLat, userLon, issLat, issLon)
	return math.Sqrt(groundDist*groundDist + issAltitudeKm*issAltitudeKm)
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := toRad(lat2 - lat1)
	dLon := toRad(lon2 - lon1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(toRad(lat1))*math.Cos(toRad(lat2))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}

func toRad(deg float64) float64 {
	return deg * math.Pi / 180.0
}

// EstimateNextPass provides a rough estimate of when the ISS will next be
// near the user's location. The ISS orbits at ~7.66 km/s with a ~90 min period.
// This is a simplified calculation - not astronomically precise.
func EstimateNextPass(userLat, userLon, issLat, issLon float64) time.Time {
	groundDist := haversine(userLat, userLon, issLat, issLon)
	issSpeedKmPerSec := 7.66
	secondsAway := groundDist / issSpeedKmPerSec

	// ISS orbital period is ~92 minutes
	orbitalPeriod := 92.0 * 60.0

	// If ISS is very far, estimate based on orbital mechanics
	if secondsAway > orbitalPeriod {
		// Normalize to within a few orbits
		orbits := math.Ceil(secondsAway / orbitalPeriod)
		secondsAway = orbits * orbitalPeriod * 0.3 // rough correction factor
	}

	// The ISS visibility window is within ~10 degrees of latitude
	latDiff := math.Abs(userLat - issLat)
	if latDiff < 10 {
		// ISS is relatively close in latitude, could pass soon
		lonDiff := math.Abs(userLon - issLon)
		if lonDiff > 180 {
			lonDiff = 360 - lonDiff
		}
		// Earth rotates ~15 degrees/hour, ISS moves ~360 degrees/92 min
		secondsAway = lonDiff / (360.0 / orbitalPeriod)
	}

	if secondsAway < 60 {
		secondsAway = orbitalPeriod // at least one orbit away
	}

	return time.Now().Add(time.Duration(secondsAway) * time.Second)
}
