package commands

import (
	"context"

	"github.com/burakberk/lubelogger-cli/internal/client"
	"github.com/burakberk/lubelogger-cli/internal/config"
	"github.com/burakberk/lubelogger-cli/internal/output"
)

// RunStatus checks connection to LubeLogger and queries server health and identity.
func RunStatus(ctx context.Context, cfg *config.Config, c *client.Client) {
	health, err := c.GetServerHealth(ctx)
	if err != nil {
		if apiErr, ok := err.(*client.APIError); ok {
			switch apiErr.Code {
			case "AUTH_ERROR":
				output.ExitWithError("AUTH_ERROR", apiErr.Message, apiErr.StatusCode, output.ExitAuthError, cfg.JSONOutput)
			case "NETWORK_ERROR":
				output.ExitWithError("NETWORK_ERROR", apiErr.Message, apiErr.StatusCode, output.ExitNetworkError, cfg.JSONOutput)
			default:
				output.ExitWithError(apiErr.Code, apiErr.Message, apiErr.StatusCode, output.ExitGeneralError, cfg.JSONOutput)
			}
		}
		output.ExitWithError("CONNECTIVITY_ERROR", err.Error(), 0, output.ExitNetworkError, cfg.JSONOutput)
	}

	res := &client.StatusResponse{
		Status:        health.Status,
		ServerVersion: health.Version,
		URL:           cfg.URL,
		Authenticated: false,
	}

	// Try WhoAmI if API Key or Basic Auth is configured
	if cfg.APIKey != "" || (cfg.Username != "" && cfg.Password != "") {
		user, err := c.GetWhoAmI(ctx)
		if err == nil && user != nil {
			res.Authenticated = true
			res.User = user
		}
	}

	if err := output.FormatStatus(res, cfg.JSONOutput); err != nil {
		output.ExitWithError("FORMAT_ERROR", err.Error(), 0, output.ExitGeneralError, cfg.JSONOutput)
	}
}
