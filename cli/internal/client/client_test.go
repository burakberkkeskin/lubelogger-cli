package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestClient_GetServerHealth(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" {
			t.Errorf("expected path /health, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"version": "1.4.2", "status": "pass", "totalDuration": "00:00:00.005", "checks": [{"name": "Postgres", "status": "pass"}]}`))
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "test-api-key", "", "", 2*time.Second)
	health, err := c.GetServerHealth(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if health.Version != "1.4.2" {
		t.Errorf("expected version 1.4.2, got %s", health.Version)
	}
	if health.Status != "pass" {
		t.Errorf("expected status pass, got %s", health.Status)
	}
}

func TestClient_AuthHeaderInjection(t *testing.T) {
	t.Run("API Key header", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if apiKey := r.Header.Get("x-api-key"); apiKey != "my-secret-key" {
				t.Errorf("expected x-api-key header 'my-secret-key', got '%s'", apiKey)
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`[]`))
		}))
		defer ts.Close()

		c := NewClient(ts.URL, "my-secret-key", "", "", 2*time.Second)
		_, err := c.GetVehicles(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("Basic Auth header fallback", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			if !ok || user != "admin" || pass != "pass123" {
				t.Errorf("expected basic auth admin:pass123, got %s:%s (ok=%v)", user, pass, ok)
			}
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`[]`))
		}))
		defer ts.Close()

		c := NewClient(ts.URL, "", "admin", "pass123", 2*time.Second)
		_, err := c.GetVehicles(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestClient_GetVehicles(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/vehicles" {
			t.Errorf("expected /api/vehicles, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[
			{"id": 1, "year": 2021, "make": "Toyota", "model": "Corolla", "licensePlate": "34ABC123", "isElectric": false, "tags": ["Daily"]},
			{"id": 2, "year": 2023, "make": "Tesla", "model": "Model 3", "licensePlate": "34EV999", "isElectric": true}
		]`))
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "", "", "", 2*time.Second)
	vehicles, err := c.GetVehicles(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(vehicles) != 2 {
		t.Fatalf("expected 2 vehicles, got %d", len(vehicles))
	}
	if vehicles[0].Make != "Toyota" || vehicles[1].IsElectric != true {
		t.Errorf("vehicle parsing mismatch: %+v", vehicles)
	}
}

func TestClient_GetVehicleInfo(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/vehicle/info" {
			t.Errorf("expected /api/vehicle/info, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("vehicleId") != "1" {
			t.Errorf("expected vehicleId=1 query parameter, got %s", r.URL.Query().Get("vehicleId"))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{
			"vehicleData": {"id": 1, "year": 2021, "make": "Toyota", "model": "Corolla"},
			"lastReportedOdometer": 45000,
			"serviceRecordCount": 5,
			"serviceRecordCost": 1200.50,
			"urgentReminderCount": 1,
			"nextReminder": {"description": "Oil Change", "urgency": "Urgent"}
		}]`))
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "", "", "", 2*time.Second)
	info, err := c.GetVehicleInfo(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.VehicleData.Make != "Toyota" {
		t.Errorf("expected Toyota, got %s", info.VehicleData.Make)
	}
	if info.LastReportedOdometer != 45000 {
		t.Errorf("expected odometer 45000, got %d", info.LastReportedOdometer)
	}
	if info.NextReminder == nil || info.NextReminder.Description != "Oil Change" {
		t.Errorf("expected NextReminder Oil Change, got %+v", info.NextReminder)
	}
}

func TestClient_Errors(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/unauthorized" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if r.URL.Path == "/api/notfound" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal database failure"))
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "", "", "", 2*time.Second)

	reqAuth, _ := c.newRequest(context.Background(), "GET", "/api/unauthorized", nil)
	err := c.doRequest(reqAuth, nil)
	if apiErr, ok := err.(*APIError); !ok || apiErr.Code != "AUTH_ERROR" {
		t.Errorf("expected AUTH_ERROR, got %v", err)
	}

	reqNotFound, _ := c.newRequest(context.Background(), "GET", "/api/notfound", nil)
	err = c.doRequest(reqNotFound, nil)
	if apiErr, ok := err.(*APIError); !ok || apiErr.Code != "NOT_FOUND" {
		t.Errorf("expected NOT_FOUND, got %v", err)
	}
}

func TestClient_GetGasRecords(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/vehicle/gasrecords" {
			t.Errorf("expected path /api/vehicle/gasrecords, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[
			{"id": 40, "vehicleId": 1, "date": "2025-09-22", "odometer": 56003, "fuelConsumed": 41.43, "cost": 2200, "files": []}
		]`))
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "", "", "", 2*time.Second)
	records, err := c.GetGasRecords(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if fmt.Sprintf("%v", records[0].ID) != "40" {
		t.Errorf("expected record ID 40, got %v", records[0].ID)
	}
}

func TestClient_DownloadFile(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/documents/test.pdf" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/pdf")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("%PDF-1.4 test document content"))
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "", "", "", 2*time.Second)

	// Successful download
	rc, contentType, length, err := c.DownloadFile(context.Background(), "/documents/test.pdf")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer rc.Close()

	if contentType != "application/pdf" {
		t.Errorf("expected content-type application/pdf, got %s", contentType)
	}
	if length != int64(len("%PDF-1.4 test document content")) {
		t.Errorf("expected length %d, got %d", len("%PDF-1.4 test document content"), length)
	}

	// 404 download
	_, _, _, err = c.DownloadFile(context.Background(), "/documents/nonexistent.pdf")
	if err == nil {
		t.Fatalf("expected error for 404 file, got nil")
	}
	if apiErr, ok := err.(*APIError); !ok || apiErr.Code != "NOT_FOUND" {
		t.Errorf("expected NOT_FOUND APIError, got %v", err)
	}
}

func TestClient_UploadDocuments(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/documents/upload" {
			t.Errorf("expected path /api/documents/upload, got %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected method POST, got %s", r.Method)
		}

		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			t.Fatalf("failed to parse multipart form: %v", err)
		}
		files := r.MultipartForm.File["documents"]
		if len(files) != 1 {
			t.Fatalf("expected 1 uploaded file under 'documents', got %d", len(files))
		}
		if files[0].Filename != "test-receipt.pdf" {
			t.Errorf("expected filename test-receipt.pdf, got %s", files[0].Filename)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[{"name": "test-receipt.pdf", "location": "/documents/uuid-1234.pdf", "isPending": false}]`))
	}))
	defer ts.Close()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test-receipt.pdf")
	if err := os.WriteFile(filePath, []byte("dummy file content"), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	c := NewClient(ts.URL, "test-key", "", "", 2*time.Second)
	res, err := c.UploadDocuments(context.Background(), []string{filePath})
	if err != nil {
		t.Fatalf("unexpected error uploading documents: %v", err)
	}

	if len(res) != 1 {
		t.Fatalf("expected 1 result, got %d", len(res))
	}
	if res[0].Name != "test-receipt.pdf" || res[0].Location != "/documents/uuid-1234.pdf" {
		t.Errorf("unexpected file descriptor: %+v", res[0])
	}
}

func TestClient_AddServiceRecord(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/vehicle/servicerecords/add" {
			t.Errorf("expected path /api/vehicle/servicerecords/add, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("vehicleId") != "3" {
			t.Errorf("expected vehicleId=3, got %s", r.URL.Query().Get("vehicleId"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		var payload ServiceRecordInput
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("failed to decode json body: %v", err)
		}
		if payload.Description != "Oil Change" || payload.Odometer != "50000" {
			t.Errorf("unexpected payload content: %+v", payload)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success": true, "message": "Service Record Added", "additionalData": {"recordId": 19}}`))
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "test-key", "", "", 2*time.Second)
	res, err := c.AddServiceRecord(context.Background(), 3, ServiceRecordInput{
		Date:        "2026-10-01",
		Odometer:    "50000",
		Description: "Oil Change",
		Cost:        "1500",
	})
	if err != nil {
		t.Fatalf("unexpected error adding service record: %v", err)
	}

	if !res.Success {
		t.Errorf("expected success true, got false")
	}
	if fmt.Sprintf("%v", res.RecordID) != "19" {
		t.Errorf("expected RecordID 19, got %v", res.RecordID)
	}
	if res.VehicleID != 3 {
		t.Errorf("expected VehicleID 3, got %d", res.VehicleID)
	}
}

func TestClient_AddGasRecord(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/vehicle/gasrecords/add" {
			t.Errorf("expected path /api/vehicle/gasrecords/add, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("vehicleId") != "2" {
			t.Errorf("expected vehicleId=2, got %s", r.URL.Query().Get("vehicleId"))
		}

		var payload GasRecordInput
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("failed to decode json body: %v", err)
		}
		if payload.FuelConsumed != "45.5" || payload.IsFillToFull != "true" || payload.MissedFuelUp != "false" {
			t.Errorf("unexpected payload: %+v", payload)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success": true, "message": "Gas Record Added", "additionalData": {"recordId": 42}}`))
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "test-key", "", "", 2*time.Second)
	res, err := c.AddGasRecord(context.Background(), 2, GasRecordInput{
		Date:         "2026-10-01",
		Odometer:     "60000",
		FuelConsumed: "45.5",
		Cost:         "2100",
		IsFillToFull: "true",
		MissedFuelUp: "false",
	})
	if err != nil {
		t.Fatalf("unexpected error adding gas record: %v", err)
	}

	if fmt.Sprintf("%v", res.RecordID) != "42" {
		t.Errorf("expected RecordID 42, got %v", res.RecordID)
	}
}

func TestClient_AddNoteRecord(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/vehicle/notes/add" {
			t.Errorf("expected path /api/vehicle/notes/add, got %s", r.URL.Path)
		}
		if r.URL.Query().Get("vehicleId") != "1" {
			t.Errorf("expected vehicleId=1, got %s", r.URL.Query().Get("vehicleId"))
		}

		var payload NoteRecordInput
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("failed to decode json body: %v", err)
		}
		if payload.Description != "Tire Rotation" || payload.NoteText != "Done at front shop" {
			t.Errorf("unexpected note payload: %+v", payload)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success": true, "message": "Note Added", "additionalData": {"recordId": 8}}`))
	}))
	defer ts.Close()

	c := NewClient(ts.URL, "test-key", "", "", 2*time.Second)
	res, err := c.AddNoteRecord(context.Background(), 1, NoteRecordInput{
		Description: "Tire Rotation",
		NoteText:    "Done at front shop",
	})
	if err != nil {
		t.Fatalf("unexpected error adding note record: %v", err)
	}

	if fmt.Sprintf("%v", res.RecordID) != "8" {
		t.Errorf("expected RecordID 8, got %v", res.RecordID)
	}
}


