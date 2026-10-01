package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
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

	if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusSeeOther || resp.StatusCode == http.StatusTemporaryRedirect {
		loc := resp.Header.Get("Location")
		if strings.Contains(strings.ToLower(loc), "login") {
			return &APIError{
				StatusCode: http.StatusUnauthorized,
				Code:       "AUTH_ERROR",
				Message:    "authentication required (server redirected to login page; please supply a valid API key or credentials)",
			}
		}
	}

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
		var opResp OperationResponse
		if err := json.Unmarshal(bodyBytes, &opResp); err == nil && opResp.Message != "" {
			return &APIError{
				StatusCode: resp.StatusCode,
				Code:       "API_ERROR",
				Message:    opResp.Message,
			}
		}
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

// GetServiceRecords calls /api/vehicle/servicerecords?vehicleId={id}.
func (c *Client) GetServiceRecords(ctx context.Context, vehicleID int) ([]ServiceRecord, error) {
	path := fmt.Sprintf("/api/vehicle/servicerecords?vehicleId=%d", vehicleID)
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var records []ServiceRecord
	if err := c.doRequest(req, &records); err != nil {
		return nil, err
	}
	return records, nil
}

// GetNotes calls /api/vehicle/notes?vehicleId={id}.
func (c *Client) GetNotes(ctx context.Context, vehicleID int) ([]NoteRecord, error) {
	path := fmt.Sprintf("/api/vehicle/notes?vehicleId=%d", vehicleID)
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var notes []NoteRecord
	if err := c.doRequest(req, &notes); err != nil {
		return nil, err
	}
	return notes, nil
}

// GetGasRecords calls /api/vehicle/gasrecords?vehicleId={id}.
func (c *Client) GetGasRecords(ctx context.Context, vehicleID int) ([]GasRecord, error) {
	path := fmt.Sprintf("/api/vehicle/gasrecords?vehicleId=%d", vehicleID)
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var records []GasRecord
	if err := c.doRequest(req, &records); err != nil {
		return nil, err
	}
	return records, nil
}

// DownloadFile fetches a file attachment from the LubeLogger server.
// Returns the reader for file bytes, Content-Type, Content-Length, and error.
// The caller is responsible for closing the returned ReadCloser.
func (c *Client) DownloadFile(ctx context.Context, location string) (io.ReadCloser, string, int64, error) {
	cleanPath := location
	if !strings.HasPrefix(cleanPath, "/") && !strings.HasPrefix(cleanPath, "http://") && !strings.HasPrefix(cleanPath, "https://") {
		cleanPath = "/" + cleanPath
	}

	var req *http.Request
	var err error
	if strings.HasPrefix(cleanPath, "http://") || strings.HasPrefix(cleanPath, "https://") {
		req, err = http.NewRequestWithContext(ctx, http.MethodGet, cleanPath, nil)
	} else {
		req, err = http.NewRequestWithContext(ctx, http.MethodGet, c.BaseURL+cleanPath, nil)
	}
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to create download request: %w", err)
	}

	// Apply authentication
	if c.APIKey != "" {
		req.Header.Set("x-api-key", c.APIKey)
	} else if c.Username != "" && c.Password != "" {
		creds := base64.StdEncoding.EncodeToString([]byte(c.Username + ":" + c.Password))
		req.Header.Set("Authorization", "Basic "+creds)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, "", 0, &APIError{
			StatusCode: 0,
			Code:       "NETWORK_ERROR",
			Message:    err.Error(),
		}
	}

	if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusSeeOther || resp.StatusCode == http.StatusTemporaryRedirect {
		loc := resp.Header.Get("Location")
		if strings.Contains(strings.ToLower(loc), "login") {
			resp.Body.Close()
			return nil, "", 0, &APIError{
				StatusCode: http.StatusUnauthorized,
				Code:       "AUTH_ERROR",
				Message:    "authentication required to download file (redirected to login)",
			}
		}
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		resp.Body.Close()
		return nil, "", 0, &APIError{
			StatusCode: resp.StatusCode,
			Code:       "AUTH_ERROR",
			Message:    "authentication failed or access denied for file download",
		}
	}

	if resp.StatusCode == http.StatusNotFound {
		resp.Body.Close()
		return nil, "", 0, &APIError{
			StatusCode: http.StatusNotFound,
			Code:       "NOT_FOUND",
			Message:    fmt.Sprintf("file not found at location '%s'", location),
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		resp.Body.Close()
		return nil, "", 0, &APIError{
			StatusCode: resp.StatusCode,
			Code:       "HTTP_ERROR",
			Message:    fmt.Sprintf("unexpected HTTP status %d downloading file", resp.StatusCode),
		}
	}

	contentType := resp.Header.Get("Content-Type")
	contentLength := resp.ContentLength
	return resp.Body, contentType, contentLength, nil
}

// UploadDocuments streams local files to POST /api/documents/upload using multipart/form-data.
// Returns the uploaded file descriptors returned by the server.
func (c *Client) UploadDocuments(ctx context.Context, filePaths []string) ([]UploadedFile, error) {
	if len(filePaths) == 0 {
		return nil, nil
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for _, fp := range filePaths {
		file, err := os.Open(fp)
		if err != nil {
			return nil, fmt.Errorf("failed to open file '%s': %w", fp, err)
		}

		stat, err := file.Stat()
		if err != nil {
			file.Close()
			return nil, fmt.Errorf("failed to stat file '%s': %w", fp, err)
		}
		if stat.IsDir() {
			file.Close()
			return nil, fmt.Errorf("path '%s' is a directory, not a file", fp)
		}

		part, err := writer.CreateFormFile("documents", filepath.Base(fp))
		if err != nil {
			file.Close()
			return nil, fmt.Errorf("failed to create form file for '%s': %w", fp, err)
		}
		if _, err := io.Copy(part, file); err != nil {
			file.Close()
			return nil, fmt.Errorf("failed to copy file content for '%s': %w", fp, err)
		}
		file.Close()
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/api/documents/upload", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var uploaded []UploadedFile
	if err := c.doRequest(req, &uploaded); err != nil {
		return nil, err
	}

	return uploaded, nil
}

// AddServiceRecord calls POST /api/vehicle/servicerecords/add?vehicleId={id}.
func (c *Client) AddServiceRecord(ctx context.Context, vehicleID int, input ServiceRecordInput) (*MutationResult, error) {
	if vehicleID <= 0 {
		return nil, fmt.Errorf("invalid vehicle ID: %d", vehicleID)
	}

	payload, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal service record payload: %w", err)
	}

	path := fmt.Sprintf("/api/vehicle/servicerecords/add?vehicleId=%d", vehicleID)
	req, err := c.newRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var resp OperationResponse
	if err := c.doRequest(req, &resp); err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("server rejected service record: %s", resp.Message)
	}

	return &MutationResult{
		Success:     true,
		RecordType:  "service",
		RecordID:    parseRecordID(resp.AdditionalData),
		VehicleID:   vehicleID,
		Date:        input.Date,
		Odometer:    input.Odometer,
		Cost:        input.Cost,
		Description: input.Description,
		Files:       input.Files,
	}, nil
}

// AddGasRecord calls POST /api/vehicle/gasrecords/add?vehicleId={id}.
func (c *Client) AddGasRecord(ctx context.Context, vehicleID int, input GasRecordInput) (*MutationResult, error) {
	if vehicleID <= 0 {
		return nil, fmt.Errorf("invalid vehicle ID: %d", vehicleID)
	}

	payload, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal gas record payload: %w", err)
	}

	path := fmt.Sprintf("/api/vehicle/gasrecords/add?vehicleId=%d", vehicleID)
	req, err := c.newRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var resp OperationResponse
	if err := c.doRequest(req, &resp); err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("server rejected gas record: %s", resp.Message)
	}

	return &MutationResult{
		Success:      true,
		RecordType:   "fuel",
		RecordID:     parseRecordID(resp.AdditionalData),
		VehicleID:    vehicleID,
		Date:         input.Date,
		Odometer:     input.Odometer,
		Cost:         input.Cost,
		FuelConsumed: input.FuelConsumed,
		Files:        input.Files,
	}, nil
}

// AddNoteRecord calls POST /api/vehicle/notes/add?vehicleId={id}.
func (c *Client) AddNoteRecord(ctx context.Context, vehicleID int, input NoteRecordInput) (*MutationResult, error) {
	if vehicleID <= 0 {
		return nil, fmt.Errorf("invalid vehicle ID: %d", vehicleID)
	}

	payload, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal note record payload: %w", err)
	}

	path := fmt.Sprintf("/api/vehicle/notes/add?vehicleId=%d", vehicleID)
	req, err := c.newRequest(ctx, http.MethodPost, path, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var resp OperationResponse
	if err := c.doRequest(req, &resp); err != nil {
		return nil, err
	}

	if !resp.Success {
		return nil, fmt.Errorf("server rejected note record: %s", resp.Message)
	}

	return &MutationResult{
		Success:     true,
		RecordType:  "note",
		RecordID:    parseRecordID(resp.AdditionalData),
		VehicleID:   vehicleID,
		Description: input.Description,
		Files:       input.Files,
	}, nil
}

func parseRecordID(data map[string]interface{}) interface{} {
	if data == nil {
		return nil
	}
	if val, ok := data["recordId"]; ok {
		if f, ok := val.(float64); ok {
			return int64(f)
		}
		return val
	}
	return nil
}

