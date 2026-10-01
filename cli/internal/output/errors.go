package output

import (
	"encoding/json"
	"fmt"
	"os"
)

const (
	ExitSuccess      = 0
	ExitGeneralError = 1
	ExitAuthError    = 2
	ExitNotFound     = 3
	ExitNetworkError = 4
	ExitUsageError   = 5
)

// ErrorEnvelope is the structured error format emitted to stderr when in JSON mode.
type ErrorEnvelope struct {
	Error      bool   `json:"error"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode,omitempty"`
}

// ExitWithError prints the error in the requested format and terminates the process with the given exit code.
func ExitWithError(code string, msg string, statusCode int, exitCode int, asJSON bool) {
	if asJSON {
		env := ErrorEnvelope{
			Error:      true,
			Code:       code,
			Message:    msg,
			StatusCode: statusCode,
		}
		bytes, _ := json.MarshalIndent(env, "", "  ")
		fmt.Fprintln(os.Stderr, string(bytes))
	} else {
		if statusCode > 0 {
			fmt.Fprintf(os.Stderr, "Error [%s] (status %d): %s\n", code, statusCode, msg)
		} else {
			fmt.Fprintf(os.Stderr, "Error [%s]: %s\n", code, msg)
		}
	}
	os.Exit(exitCode)
}
