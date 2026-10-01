package commands

import (
	"context"
	"errors"
	"fmt"

	"github.com/burakberk/lubelogger-cli/internal/client"
	"github.com/burakberk/lubelogger-cli/internal/config"
	"github.com/burakberk/lubelogger-cli/internal/output"
)

// RunServiceDownload downloads all files attached to a specific service record.
func RunServiceDownload(ctx context.Context, cfg *config.Config, apiClient *client.Client, vehicleID int, recordID int, destDir string) {
	records, err := apiClient.GetServiceRecords(ctx, vehicleID)
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

	var targetRecord *client.ServiceRecord
	for i := range records {
		recIDStr := fmt.Sprintf("%v", records[i].ID)
		targetIDStr := fmt.Sprintf("%d", recordID)
		if recIDStr == targetIDStr {
			targetRecord = &records[i]
			break
		}
	}

	if targetRecord == nil {
		output.ExitWithError("RECORD_NOT_FOUND", fmt.Sprintf("service record ID %d not found for vehicle %d", recordID, vehicleID), 404, output.ExitNotFoundError, cfg.JSONOutput)
	}

	if len(targetRecord.Files) == 0 {
		if cfg.JSONOutput {
			_ = output.PrintJSON([]client.DownloadedFile{})
		} else {
			fmt.Printf("No attachment files found for service record #%d.\n", recordID)
		}
		return
	}

	downloaded, err := DownloadRecordFiles(ctx, apiClient, targetRecord.Files, destDir)
	if err != nil {
		output.ExitWithError("DOWNLOAD_FAILED", err.Error(), 0, output.ExitGeneralError, cfg.JSONOutput)
	}

	_ = output.FormatDownloadedFiles(downloaded, cfg.JSONOutput)
}
