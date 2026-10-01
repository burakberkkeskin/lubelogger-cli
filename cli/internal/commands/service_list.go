package commands

import (
	"context"
	"sort"
	"strings"
	"time"

	"github.com/burakberk/lubelogger-cli/internal/client"
	"github.com/burakberk/lubelogger-cli/internal/config"
	"github.com/burakberk/lubelogger-cli/internal/output"
)

func parseRecordDate(d string) time.Time {
	formats := []string{
		"2.1.2006",
		"02.01.2006",
		"2006-01-02",
		"1/2/2006",
		"01/02/2006",
	}
	clean := strings.TrimSpace(d)
	for _, f := range formats {
		if t, err := time.Parse(f, clean); err == nil {
			return t
		}
	}
	return time.Time{}
}

// RunServiceList queries and displays service records for a vehicle.
func RunServiceList(ctx context.Context, cfg *config.Config, c *client.Client, vehicleID int, limit int) {
	if vehicleID <= 0 {
		output.ExitWithError("INVALID_ARGUMENT", "vehicle ID must be a positive integer", 0, output.ExitUsageError, cfg.JSONOutput)
	}

	records, err := c.GetServiceRecords(ctx, vehicleID)
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok {
			switch apiErr.Code {
			case "AUTH_ERROR":
				output.ExitWithError("AUTH_ERROR", apiErr.Message, apiErr.StatusCode, output.ExitAuthError, cfg.JSONOutput)
			case "NOT_FOUND":
				output.ExitWithError("NOT_FOUND", apiErr.Message, apiErr.StatusCode, output.ExitNotFound, cfg.JSONOutput)
			case "NETWORK_ERROR":
				output.ExitWithError("NETWORK_ERROR", apiErr.Message, apiErr.StatusCode, output.ExitNetworkError, cfg.JSONOutput)
			default:
				output.ExitWithError(apiErr.Code, apiErr.Message, apiErr.StatusCode, output.ExitGeneralError, cfg.JSONOutput)
			}
		}
		output.ExitWithError("QUERY_ERROR", err.Error(), 0, output.ExitGeneralError, cfg.JSONOutput)
	}

	// Sort descending by date (newest first)
	sort.Slice(records, func(i, j int) bool {
		tI := parseRecordDate(records[i].Date)
		tJ := parseRecordDate(records[j].Date)
		if !tI.IsZero() && !tJ.IsZero() {
			return tI.After(tJ)
		}
		return parseRecordDate(records[i].Date).After(parseRecordDate(records[j].Date))
	})

	if limit > 0 && len(records) > limit {
		records = records[:limit]
	}

	if err := output.FormatServiceRecords(records, cfg.JSONOutput); err != nil {
		output.ExitWithError("FORMAT_ERROR", err.Error(), 0, output.ExitGeneralError, cfg.JSONOutput)
	}
}
