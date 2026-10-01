package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// APIError represents an error returned by LubeLogger API.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error [%s] (status %d): %s", e.Code, e.StatusCode, e.Message)
}

// Client is the HTTP client for LubeLogger REST API.
type Client struct {
	BaseURL    string
	APIKey     string
	Username   string
	Password   string
	HTTPClient *http.Client
}

// NewClient creates a new LubeLogger API client.
func NewClient(baseURL string, apiKey string, username string, password string, timeout time.Duration) *Client {
	cleanURL := strings.TrimRight(baseURL, "/")
	if timeout <= 0 {
		timeout = 10 * time.Second
	}
	return &Client{
		BaseURL:  cleanURL,
		APIKey:   apiKey,
		Username: username,
		Password: password,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	fullURL := c.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("culture-invariant", "true")

	// Apply authentication
	if c.APIKey != "" {
		req.Header.Set("x-api-key", c.APIKey)
	} else if c.Username != "" && c.Password != "" {
		creds := base64.StdEncoding.EncodeToString([]byte(c.Username + ":" + c.Password))
		req.Header.Set("Authorization", "Basic "+creds)
	}

	return req, nil
}

func (c *Client) doRequest(req *http.Request, target interface{}) error {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return &APIError{
			StatusCode: 0,
			Code:       "NETWORK_ERROR",
			Message:    err.Error(),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return &APIError{
			StatusCode: resp.StatusCode,
			Code:       "AUTH_ERROR",
			Message:    "authentication failed or access denied (verify URL and API key)",
		}
	}

	if resp.StatusCode == http.StatusNotFound {
		return &APIError{
			StatusCode: resp.StatusCode,
			Code:       "NOT_FOUND",
			Message:    "resource not found",
		}
	}

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return &APIError{
			StatusCode: resp.StatusCode,
			Code:       "SERVER_ERROR",
			Message:    string(bodyBytes),
		}
	}

	if target != nil {
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			return fmt.Errorf("failed to decode JSON response: %w", err)
		}
	}

	return nil
}

// GetServerHealth calls /health.
func (c *Client) GetServerHealth(ctx context.Context) (*ServerHealth, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/health", nil)
	if err != nil {
		return nil, err
	}
	var health ServerHealth
	if err := c.doRequest(req, &health); err != nil {
		return nil, err
	}
	return &health, nil
}

// GetWhoAmI calls /api/whoami.
func (c *Client) GetWhoAmI(ctx context.Context) (*UserExportModel, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/whoami", nil)
	if err != nil {
		return nil, err
	}
	var user UserExportModel
	if err := c.doRequest(req, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// GetVehicles calls /api/vehicles.
func (c *Client) GetVehicles(ctx context.Context) ([]Vehicle, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/vehicles", nil)
	if err != nil {
		return nil, err
	}
	var vehicles []Vehicle
	if err := c.doRequest(req, &vehicles); err != nil {
		return nil, err
	}
	return vehicles, nil
}

// GetVehicleInfo calls /api/vehicle/info?vehicleId={id}.
func (c *Client) GetVehicleInfo(ctx context.Context, vehicleID int) (*VehicleInfo, error) {
	path := fmt.Sprintf("/api/vehicle/info?vehicleId=%d", vehicleID)
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var infos []VehicleInfo
	if err := c.doRequest(req, &infos); err != nil {
		return nil, err
	}
	if len(infos) == 0 {
		return nil, &APIError{
			StatusCode: http.StatusNotFound,
			Code:       "VEHICLE_NOT_FOUND",
			Message:    fmt.Sprintf("no vehicle found with ID %d", vehicleID),
		}
	}
	return &infos[0], nil
}
