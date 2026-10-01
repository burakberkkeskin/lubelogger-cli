package commands

import (
	"context"
	"errors"

	"github.com/burakberk/lubelogger-cli/internal/client"
	"github.com/burakberk/lubelogger-cli/internal/config"
	"github.com/burakberk/lubelogger-cli/internal/output"
)

// RunFuelList handles the 'fuel list <vehicleId>' command.
func RunFuelList(ctx context.Context, cfg *config.Config, apiClient *client.Client, vehicleID int, limit int) {
	records, err := apiClient.GetGasRecords(ctx, vehicleID)
	if err != nil {
		var apiErr *client.APIError
		if errors.As(err, &apiErr) {
			exitCode := output.ExitGeneralError
			if apiErr.StatusCode == 401 || apiErr.StatusCode == 403 {
				exitCode = output.ExitAuthError
			} else if apiErr.StatusCode == 404 {
				exitCode = output.ExitNotFoundError
			}
			output.ExitWithError(apiErr.Code, apiErr.Message, apiErr.StatusCode, exitCode, cfg.JSONOutput)
		}
		output.ExitWithError("API_ERROR", err.Error(), 0, output.ExitGeneralError, cfg.JSONOutput)
	}

	if limit > 0 && len(records) > limit {
		records = records[:limit]
	}

	if err := output.FormatGasRecords(records, cfg.JSONOutput); err != nil {
		output.ExitWithError("FORMAT_ERROR", err.Error(), 0, output.ExitGeneralError, cfg.JSONOutput)
	}
}
