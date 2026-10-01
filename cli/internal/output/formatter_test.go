package output

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/burakberk/lubelogger-cli/internal/client"
)

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestFormatVehicles_JSON(t *testing.T) {
	vehicles := []client.Vehicle{
		{
			ID:           1,
			Year:         2021,
			Make:         "Toyota",
			Model:        "Corolla",
			LicensePlate: "34ABC123",
		},
	}

	out := captureStdout(func() {
		_ = FormatVehicles(vehicles, true)
	})

	if !strings.Contains(out, `"licensePlate": "34ABC123"`) {
		t.Errorf("expected JSON output containing licensePlate, got: %s", out)
	}
	if !strings.Contains(out, `"make": "Toyota"`) {
		t.Errorf("expected JSON output containing make, got: %s", out)
	}
}

func TestFormatVehicles_Table(t *testing.T) {
	vehicles := []client.Vehicle{
		{
			ID:           1,
			Year:         2021,
			Make:         "Toyota",
			Model:        "Corolla",
			LicensePlate: "34ABC123",
		},
	}

	out := captureStdout(func() {
		_ = FormatVehicles(vehicles, false)
	})

	if !strings.Contains(out, "ID") || !strings.Contains(out, "MAKE") || !strings.Contains(out, "MODEL") {
		t.Errorf("expected table header with ID, MAKE, MODEL, got: %s", out)
	}
	if !strings.Contains(out, "Toyota") || !strings.Contains(out, "Corolla") {
		t.Errorf("expected table row with Toyota Corolla, got: %s", out)
	}
}

func TestFormatVehicleInfo_JSON(t *testing.T) {
	info := &client.VehicleInfo{
		VehicleData: client.Vehicle{
			ID:    2,
			Year:  2023,
			Make:  "Tesla",
			Model: "Model Y",
		},
		LastReportedOdometer: 15400,
		ServiceRecordCount:   2,
		ServiceRecordCost:    300.0,
	}

	out := captureStdout(func() {
		_ = FormatVehicleInfo(info, true)
	})

	if !strings.Contains(out, `"lastReportedOdometer": 15400`) {
		t.Errorf("expected JSON output containing lastReportedOdometer, got: %s", out)
	}
	if !strings.Contains(out, `"model": "Model Y"`) {
		t.Errorf("expected JSON output containing Model Y, got: %s", out)
	}
}

func TestFormatMutationResult(t *testing.T) {
	res := &client.MutationResult{
		Success:     true,
		RecordType:  "service",
		RecordID:    int64(19),
		VehicleID:   3,
		Date:        "2026-10-01",
		Odometer:    "70500",
		Cost:        "4500",
		Description: "Periyodik Bakım",
		Files: []client.UploadedFile{
			{Name: "invoice.pdf", Location: "/documents/abc.pdf"},
		},
	}

	t.Run("JSON output", func(t *testing.T) {
		out := captureStdout(func() {
			_ = FormatMutationResult(res, true)
		})
		if !strings.Contains(out, `"recordType": "service"`) || !strings.Contains(out, `"recordId": 19`) {
			t.Errorf("unexpected json: %s", out)
		}
	})

	t.Run("Text output", func(t *testing.T) {
		out := captureStdout(func() {
			_ = FormatMutationResult(res, false)
		})
		if !strings.Contains(out, "Created service record #19 for vehicle #3") {
			t.Errorf("unexpected text output: %s", out)
		}
		if !strings.Contains(out, "70500") || !strings.Contains(out, "invoice.pdf") {
			t.Errorf("missing details in text output: %s", out)
		}
	})
}

func TestFormatUploadedFiles(t *testing.T) {
	files := []client.UploadedFile{
		{Name: "doc1.pdf", Location: "/documents/123.pdf"},
		{Name: "doc2.jpg", Location: "/documents/456.jpg"},
	}

	t.Run("JSON output", func(t *testing.T) {
		out := captureStdout(func() {
			_ = FormatUploadedFiles(files, true)
		})
		if !strings.Contains(out, "/documents/123.pdf") {
			t.Errorf("unexpected json: %s", out)
		}
	})

	t.Run("Text output", func(t *testing.T) {
		out := captureStdout(func() {
			_ = FormatUploadedFiles(files, false)
		})
		if !strings.Contains(out, "Uploaded 2 file(s)") || !strings.Contains(out, "doc1.pdf -> /documents/123.pdf") {
			t.Errorf("unexpected text output: %s", out)
		}
	})
}

