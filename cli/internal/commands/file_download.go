package commands

import (
	"context"
	"errors"
	"os"
	"path"
	"path/filepath"

	"github.com/burakberk/lubelogger-cli/internal/client"
	"github.com/burakberk/lubelogger-cli/internal/config"
	"github.com/burakberk/lubelogger-cli/internal/output"
)

// RunFileDownload handles downloading a specific file location directly.
func RunFileDownload(ctx context.Context, cfg *config.Config, apiClient *client.Client, location string, outputDest string) {
	targetPath := outputDest
	if targetPath == "" {
		base := path.Base(location)
		if base == "" || base == "/" || base == "." {
			base = "downloaded_file"
		}
		targetPath = base
	} else {
		if fi, err := os.Stat(targetPath); err == nil && fi.IsDir() {
			base := path.Base(location)
			if base == "" || base == "/" || base == "." {
				base = "downloaded_file"
			}
			targetPath = filepath.Join(targetPath, base)
		}
	}

	downloaded, err := SaveDownloadedStream(ctx, apiClient, location, targetPath)
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
		output.ExitWithError("DOWNLOAD_FAILED", err.Error(), 0, output.ExitGeneralError, cfg.JSONOutput)
	}

	_ = output.FormatDownloadedFiles([]client.DownloadedFile{*downloaded}, cfg.JSONOutput)
}
