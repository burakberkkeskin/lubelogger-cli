package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/burakberk/lubelogger-cli/internal/client"
)

// SaveDownloadedStream downloads a file from the API client and writes it to targetFile.
func SaveDownloadedStream(ctx context.Context, apiClient *client.Client, location string, targetFile string) (*client.DownloadedFile, error) {
	body, contentType, _, err := apiClient.DownloadFile(ctx, location)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	// Ensure parent directory exists
	dir := filepath.Dir(targetFile)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory '%s': %w", dir, err)
		}
	}

	out, err := os.Create(targetFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create file '%s': %w", targetFile, err)
	}
	defer out.Close()

	written, err := io.Copy(out, body)
	if err != nil {
		return nil, fmt.Errorf("failed to write file '%s': %w", targetFile, err)
	}

	return &client.DownloadedFile{
		Name:        filepath.Base(targetFile),
		Location:    location,
		SavedTo:     targetFile,
		Bytes:       written,
		ContentType: contentType,
	}, nil
}

// SanitizeFilename cleans an attachment filename for safe local disk writing.
func SanitizeFilename(name string) string {
	clean := filepath.Base(strings.TrimSpace(name))
	clean = strings.ReplaceAll(clean, "..", "")
	clean = strings.ReplaceAll(clean, "/", "_")
	clean = strings.ReplaceAll(clean, "\\", "_")
	if clean == "" || clean == "." {
		clean = "attachment"
	}
	return clean
}

// DownloadRecordFiles downloads a slice of files to destDir.
func DownloadRecordFiles(ctx context.Context, apiClient *client.Client, files []client.UploadedFile, destDir string) ([]client.DownloadedFile, error) {
	if destDir == "" {
		destDir = "."
	}
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory '%s': %w", destDir, err)
	}

	var downloaded []client.DownloadedFile
	seenNames := make(map[string]int)

	for _, f := range files {
		fileName := SanitizeFilename(f.Name)
		if fileName == "attachment" && f.Location != "" {
			fileName = SanitizeFilename(filepath.Base(f.Location))
		}

		// Prevent duplicate names in same batch
		seenNames[fileName]++
		if seenNames[fileName] > 1 {
			ext := filepath.Ext(fileName)
			raw := strings.TrimSuffix(fileName, ext)
			fileName = fmt.Sprintf("%s_%d%s", raw, seenNames[fileName]-1, ext)
		}

		targetPath := filepath.Join(destDir, fileName)
		res, err := SaveDownloadedStream(ctx, apiClient, f.Location, targetPath)
		if err != nil {
			return downloaded, fmt.Errorf("failed to download '%s' from '%s': %w", f.Name, f.Location, err)
		}
		// Preserve original display name if available
		if f.Name != "" {
			res.Name = f.Name
		}
		downloaded = append(downloaded, *res)
	}

	return downloaded, nil
}

// ValidateLocalFiles checks that all provided file paths exist, are accessible, and are regular files.
func ValidateLocalFiles(filePaths []string) error {
	for _, fp := range filePaths {
		stat, err := os.Stat(fp)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("file not found: '%s'", fp)
			}
			return fmt.Errorf("cannot access file '%s': %w", fp, err)
		}
		if stat.IsDir() {
			return fmt.Errorf("path '%s' is a directory, not a regular file", fp)
		}
	}
	return nil
}
