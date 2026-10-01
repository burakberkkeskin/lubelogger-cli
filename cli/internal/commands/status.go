package commands

import (
	"context"

	"github.com/burakberk/lubelogger-cli/internal/client"
	"github.com/burakberk/lubelogger-cli/internal/config"
	"github.com/burakberk/lubelogger-cli/internal/output"
)

// RunStatus checks connection to LubeLogger and queries server health and identity.
func RunStatus(ctx context.Context, cfg *config.Config, c *client.Client) {
	res := &client.StatusResponse{
		Status:        "healthy",
		URL:           cfg.URL,
		Authenticated: false,
	}

	health, errHealth := c.GetServerHealth(ctx)
	if errHealth == nil && health != nil {
		res.Status = health.Status
		res.ServerVersion = health.Version
	}

	// Query WhoAmI
	user, errWho := c.GetWhoAmI(ctx)
	if errWho == nil && user != nil {
		res.Authenticated = true
		res.User = user
	} else if errHealth != nil {
		// If health failed, check whoami error to diagnose server status
		if apiErr, ok := errWho.(*client.APIError); ok {
			if apiErr.Code == "AUTH_ERROR" {
				// Server is reachable, but requires authentication
				res.Status = "reachable (auth required)"
				if err := output.FormatStatus(res, cfg.JSONOutput); err != nil {
					output.ExitWithError("FORMAT_ERROR", err.Error(), 0, output.ExitGeneralError, cfg.JSONOutput)
				}
				return
			}
		}

		// If health was network error
		if apiErr, ok := errHealth.(*client.APIError); ok {
			switch apiErr.Code {
			case "NETWORK_ERROR":
				output.ExitWithError("NETWORK_ERROR", apiErr.Message, apiErr.StatusCode, output.ExitNetworkError, cfg.JSONOutput)
			default:
				output.ExitWithError(apiErr.Code, apiErr.Message, apiErr.StatusCode, output.ExitGeneralError, cfg.JSONOutput)
			}
		}
		output.ExitWithError("CONNECTIVITY_ERROR", errHealth.Error(), 0, output.ExitNetworkError, cfg.JSONOutput)
	}

	if err := output.FormatStatus(res, cfg.JSONOutput); err != nil {
		output.ExitWithError("FORMAT_ERROR", err.Error(), 0, output.ExitGeneralError, cfg.JSONOutput)
	}
}
