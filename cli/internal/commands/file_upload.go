package commands

import (
	"context"
	"fmt"

	"github.com/burakberk/lubelogger-cli/internal/client"
	"github.com/burakberk/lubelogger-cli/internal/config"
	"github.com/burakberk/lubelogger-cli/internal/output"
)

// RunFileUpload uploads one or more files directly to the LubeLogger server.
func RunFileUpload(ctx context.Context, cfg *config.Config, c *client.Client, filePaths []string) {
	if len(filePaths) == 0 {
		output.ExitWithError("MISSING_ARGUMENT", "usage: 'lubelog file upload <path>...'", 0, output.ExitUsageError, cfg.JSONOutput)
	}

	if err := ValidateLocalFiles(filePaths); err != nil {
		output.ExitWithError("INVALID_FILE", err.Error(), 0, output.ExitUsageError, cfg.JSONOutput)
	}

	uploaded, err := c.UploadDocuments(ctx, filePaths)
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok {
			output.ExitWithError("UPLOAD_FAILED", apiErr.Message, apiErr.StatusCode, output.ExitNetworkError, cfg.JSONOutput)
		}
		output.ExitWithError("UPLOAD_FAILED", fmt.Sprintf("failed to upload documents: %v", err), 0, output.ExitNetworkError, cfg.JSONOutput)
	}

	if err := output.FormatUploadedFiles(uploaded, cfg.JSONOutput); err != nil {
		output.ExitWithError("FORMAT_ERROR", err.Error(), 0, output.ExitGeneralError, cfg.JSONOutput)
	}
}
