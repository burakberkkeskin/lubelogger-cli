package commands

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/burakberk/lubelogger-cli/internal/client"
	"github.com/burakberk/lubelogger-cli/internal/config"
	"github.com/burakberk/lubelogger-cli/internal/output"
)

// RunServiceAdd creates a new service record for a vehicle, optionally uploading attachments.
func RunServiceAdd(
	ctx context.Context,
	cfg *config.Config,
	c *client.Client,
	vehicleID int,
	desc string,
	odo string,
	cost string,
	date string,
	notes string,
	tags string,
	files []string,
) {
	if vehicleID <= 0 {
		output.ExitWithError("INVALID_ARGUMENT", "vehicle ID must be a positive integer", 0, output.ExitUsageError, cfg.JSONOutput)
	}

	desc = strings.TrimSpace(desc)
	if desc == "" {
		output.ExitWithError("MISSING_ARGUMENT", "missing required service description: --desc <string>", 0, output.ExitUsageError, cfg.JSONOutput)
	}

	odo = strings.TrimSpace(odo)
	if odo == "" {
		output.ExitWithError("MISSING_ARGUMENT", "missing required odometer reading: --odometer <int>", 0, output.ExitUsageError, cfg.JSONOutput)
	}

	cost = strings.TrimSpace(cost)
	if cost == "" {
		output.ExitWithError("MISSING_ARGUMENT", "missing required cost: --cost <number>", 0, output.ExitUsageError, cfg.JSONOutput)
	}

	date = strings.TrimSpace(date)
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	// Pre-flight check on local attachment files
	if err := ValidateLocalFiles(files); err != nil {
		output.ExitWithError("INVALID_FILE", err.Error(), 0, output.ExitUsageError, cfg.JSONOutput)
	}

	// Upload documents first if any specified
	var uploadedFiles []client.UploadedFile
	if len(files) > 0 {
		var err error
		uploadedFiles, err = c.UploadDocuments(ctx, files)
		if err != nil {
			if apiErr, ok := err.(*client.APIError); ok {
				output.ExitWithError("UPLOAD_FAILED", apiErr.Message, apiErr.StatusCode, output.ExitNetworkError, cfg.JSONOutput)
			}
			output.ExitWithError("UPLOAD_FAILED", fmt.Sprintf("failed to upload documents: %v", err), 0, output.ExitNetworkError, cfg.JSONOutput)
		}
	}

	input := client.ServiceRecordInput{
		Date:        date,
		Odometer:    odo,
		Description: desc,
		Cost:        cost,
		Notes:       strings.TrimSpace(notes),
		Tags:        strings.TrimSpace(tags),
		Files:       uploadedFiles,
	}

	res, err := c.AddServiceRecord(ctx, vehicleID, input)
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok {
			switch apiErr.Code {
			case "AUTH_ERROR":
				output.ExitWithError("AUTH_ERROR", apiErr.Message, apiErr.StatusCode, output.ExitAuthError, cfg.JSONOutput)
			case "NOT_FOUND":
				output.ExitWithError("NOT_FOUND", apiErr.Message, apiErr.StatusCode, output.ExitNotFound, cfg.JSONOutput)
			default:
				output.ExitWithError(apiErr.Code, apiErr.Message, apiErr.StatusCode, output.ExitGeneralError, cfg.JSONOutput)
			}
		}
		output.ExitWithError("MUTATION_FAILED", err.Error(), 0, output.ExitGeneralError, cfg.JSONOutput)
	}

	if err := output.FormatMutationResult(res, cfg.JSONOutput); err != nil {
		output.ExitWithError("FORMAT_ERROR", err.Error(), 0, output.ExitGeneralError, cfg.JSONOutput)
	}
}
