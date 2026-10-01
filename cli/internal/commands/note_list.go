package commands

import (
	"context"

	"github.com/burakberk/lubelogger-cli/internal/client"
	"github.com/burakberk/lubelogger-cli/internal/config"
	"github.com/burakberk/lubelogger-cli/internal/output"
)

// RunNoteList queries and displays notes and attached documents for a vehicle.
func RunNoteList(ctx context.Context, cfg *config.Config, c *client.Client, vehicleID int) {
	if vehicleID <= 0 {
		output.ExitWithError("INVALID_ARGUMENT", "vehicle ID must be a positive integer", 0, output.ExitUsageError, cfg.JSONOutput)
	}

	notes, err := c.GetNotes(ctx, vehicleID)
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

	if err := output.FormatNotes(notes, cfg.JSONOutput); err != nil {
		output.ExitWithError("FORMAT_ERROR", err.Error(), 0, output.ExitGeneralError, cfg.JSONOutput)
	}
}
