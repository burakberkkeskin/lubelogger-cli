package commands

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/burakberk/lubelogger-cli/internal/client"
	"github.com/burakberk/lubelogger-cli/internal/config"
	"github.com/burakberk/lubelogger-cli/internal/output"
)

// RunNoteAdd creates a new note for a vehicle, optionally uploading attachments.
func RunNoteAdd(
	ctx context.Context,
	cfg *config.Config,
	c *client.Client,
	vehicleID int,
	title string,
	text string,
	pinned bool,
	tags string,
	files []string,
) {
	if vehicleID <= 0 {
		output.ExitWithError("INVALID_ARGUMENT", "vehicle ID must be a positive integer", 0, output.ExitUsageError, cfg.JSONOutput)
	}

	title = strings.TrimSpace(title)
	if title == "" {
		output.ExitWithError("MISSING_ARGUMENT", "missing required note title: --title <string>", 0, output.ExitUsageError, cfg.JSONOutput)
	}

	text = strings.TrimSpace(text)
	if text == "" {
		text = title
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

	input := client.NoteRecordInput{
		Description: title,
		NoteText:    text,
		Pinned:      strconv.FormatBool(pinned),
		Tags:        strings.TrimSpace(tags),
		Files:       uploadedFiles,
	}

	res, err := c.AddNoteRecord(ctx, vehicleID, input)
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
