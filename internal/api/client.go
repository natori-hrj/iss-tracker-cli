package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

const maxResponseSize = 1 * 1024 * 1024 // 1MB

const (
	issNowURL  = "http://api.open-notify.org/iss-now.json"
	astrosURL  = "http://api.open-notify.org/astros.json"
)

type ISSPosition struct {
	Latitude  float64
	Longitude float64
	Timestamp time.Time
}

type Astronaut struct {
	Name  string `json:"name"`
	Craft string `json:"craft"`
}

type AstrosResponse struct {
	Message string      `json:"message"`
	Number  int         `json:"number"`
	People  []Astronaut `json:"people"`
}

type issNowResponse struct {
	Message    string `json:"message"`
	Timestamp  int64  `json:"timestamp"`
	ISSPosition struct {
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
	} `json:"iss_position"`
}

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) GetISSPosition() (*ISSPosition, error) {
	resp, err := c.httpClient.Get(issNowURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch ISS position: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ISS position API returned status %d", resp.StatusCode)
	}

	var data issNowResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxResponseSize)).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode ISS position: %w", err)
	}

	if data.Message != "success" {
		return nil, fmt.Errorf("API returned non-success message: %s", data.Message)
	}

	lat, err := strconv.ParseFloat(data.ISSPosition.Latitude, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse latitude: %w", err)
	}

	lon, err := strconv.ParseFloat(data.ISSPosition.Longitude, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to parse longitude: %w", err)
	}

	return &ISSPosition{
		Latitude:  lat,
		Longitude: lon,
		Timestamp: time.Unix(data.Timestamp, 0),
	}, nil
}

func (c *Client) GetAstronauts() (*AstrosResponse, error) {
	resp, err := c.httpClient.Get(astrosURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch astronauts: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("astronauts API returned status %d", resp.StatusCode)
	}

	var data AstrosResponse
	if err := json.NewDecoder(io.LimitReader(resp.Body, maxResponseSize)).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode astronauts: %w", err)
	}

	if data.Message != "success" {
		return nil, fmt.Errorf("API returned non-success message: %s", data.Message)
	}

	return &data, nil
}
