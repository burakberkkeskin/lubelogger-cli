package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/burakberk/lubelogger-cli/internal/client"
	"github.com/burakberk/lubelogger-cli/internal/commands"
	"github.com/burakberk/lubelogger-cli/internal/config"
	"github.com/burakberk/lubelogger-cli/internal/output"
)

const version = "0.1.0"

func printUsage() {
	fmt.Println(`LubeLogger CLI (lubelog) - Agent-First Vehicle Maintenance Query Tool

USAGE:
  lubelog [flags] <command> [arguments...]

COMMANDS:
  status                     Check LubeLogger server health and authenticated identity
  vehicle list               List all accessible vehicles (alias: vehicles)
  vehicle get <id>           Get detailed info and aggregates for a vehicle (alias: vehicle <id>)
  service list <vehicleId>   List service/maintenance records for a vehicle (alias: service <id>)
  version                    Display CLI version

GLOBAL FLAGS:
  --url <string>        Base URL of LubeLogger (default: $LUBELOGGER_URL or http://localhost:5000)
  --api-key <string>    API key for authentication (default: $LUBELOGGER_API_KEY)
  --username <string>   Username for Basic Auth fallback (default: $LUBELOGGER_USERNAME)
  --password <string>   Password for Basic Auth fallback (default: $LUBELOGGER_PASSWORD)
  --timeout <string>    Request timeout duration (default: $LUBELOGGER_TIMEOUT or 10s)
  --limit <int>         Limit number of records returned (e.g. --limit 3)
  --json                Output response in pure JSON format (default: non-TTY)
  -h, --help            Show this help message

ENVIRONMENT VARIABLES:
  LUBELOGGER_URL            Target server base URL
  LUBELOGGER_API_KEY        API key for authentication
  LUBELOGGER_USERNAME       Basic auth username
  LUBELOGGER_PASSWORD       Basic auth password
  LUBELOGGER_TIMEOUT        Request timeout (e.g. 10s)
  LUBELOGGER_OUTPUT_FORMAT  Default format ("json" or "table")

EXAMPLES:
  lubelog status --json
  lubelog vehicle list --json
  lubelog vehicle get 1 --json
  lubelog service list 1 --limit 3 --json
  lubelog --url https://lubelogger.burakberk.dev --api-key <TOKEN> vehicle list`)
}

func parseArgs(args []string) (flags []string, positional []string) {
	valueFlags := map[string]bool{
		"--url": true, "-url": true,
		"--api-key": true, "-api-key": true,
		"--username": true, "-username": true,
		"--password": true, "-password": true,
		"--timeout": true, "-timeout": true,
		"--limit": true, "-limit": true,
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
			if !strings.Contains(arg, "=") && valueFlags[arg] && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				i++
				flags = append(flags, args[i])
			}
		} else {
			positional = append(positional, arg)
		}
	}
	return flags, positional
}

func main() {
	var flagURL, flagAPIKey, flagUser, flagPass, flagTimeout string
	var flagLimit int
	var flagJSON, flagHelp bool

	flagArgs, posArgs := parseArgs(os.Args[1:])

	// Parse flags
	fs := flag.NewFlagSet("lubelog", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.StringVar(&flagURL, "url", "", "Base URL of LubeLogger")
	fs.StringVar(&flagAPIKey, "api-key", "", "API key")
	fs.StringVar(&flagUser, "username", "", "Username")
	fs.StringVar(&flagPass, "password", "", "Password")
	fs.StringVar(&flagTimeout, "timeout", "", "Request timeout")
	fs.IntVar(&flagLimit, "limit", 0, "Limit number of records returned")
	fs.BoolVar(&flagJSON, "json", false, "Force JSON output")
	fs.BoolVar(&flagHelp, "help", false, "Show help")
	fs.BoolVar(&flagHelp, "h", false, "Show help")

	if err := fs.Parse(flagArgs); err != nil {
		output.ExitWithError("INVALID_FLAGS", err.Error(), 0, output.ExitUsageError, false)
	}

	args := posArgs

	if flagHelp || len(args) == 0 {
		printUsage()
		os.Exit(output.ExitSuccess)
	}

	cfg := config.LoadConfig(flagURL, flagAPIKey, flagUser, flagPass, flagTimeout, flagJSON)
	apiClient := client.NewClient(cfg.URL, cfg.APIKey, cfg.Username, cfg.Password, cfg.Timeout)
	ctx := context.Background()

	cmd := strings.ToLower(args[0])

	switch cmd {
	case "version":
		if cfg.JSONOutput {
			_ = output.PrintJSON(map[string]string{"version": version})
		} else {
			fmt.Printf("lubelog version %s\n", version)
		}

	case "status", "ping", "health":
		commands.RunStatus(ctx, cfg, apiClient)

	case "vehicles":
		commands.RunVehicleList(ctx, cfg, apiClient)

	case "vehicle":
		if len(args) < 2 {
			output.ExitWithError("MISSING_SUBCOMMAND", "expected 'vehicle list' or 'vehicle get <id>'", 0, output.ExitUsageError, cfg.JSONOutput)
		}
		subCmd := strings.ToLower(args[1])
		switch subCmd {
		case "list", "ls":
			commands.RunVehicleList(ctx, cfg, apiClient)
		case "get", "show", "info":
			if len(args) < 3 {
				output.ExitWithError("MISSING_ARGUMENT", "missing required vehicle ID: 'lubelog vehicle get <id>'", 0, output.ExitUsageError, cfg.JSONOutput)
			}
			id, err := strconv.Atoi(args[2])
			if err != nil || id <= 0 {
				output.ExitWithError("INVALID_ARGUMENT", "vehicle ID must be a positive integer", 0, output.ExitUsageError, cfg.JSONOutput)
			}
			commands.RunVehicleGet(ctx, cfg, apiClient, id)
		default:
			// Allow direct vehicle ID shorthand: `lubelog vehicle 1`
			if id, err := strconv.Atoi(subCmd); err == nil && id > 0 {
				commands.RunVehicleGet(ctx, cfg, apiClient, id)
				return
			}
			output.ExitWithError("UNKNOWN_SUBCOMMAND", fmt.Sprintf("unknown vehicle subcommand '%s'", subCmd), 0, output.ExitUsageError, cfg.JSONOutput)
		}

	case "services", "service":
		if len(args) < 2 {
			output.ExitWithError("MISSING_SUBCOMMAND", "expected 'service list <vehicleId>' (e.g. 'lubelog service list 1 --limit 3')", 0, output.ExitUsageError, cfg.JSONOutput)
		}
		targetArg := args[1]
		if strings.ToLower(targetArg) == "list" || strings.ToLower(targetArg) == "ls" {
			if len(args) < 3 {
				output.ExitWithError("MISSING_ARGUMENT", "missing required vehicle ID: 'lubelog service list <vehicleId>'", 0, output.ExitUsageError, cfg.JSONOutput)
			}
			targetArg = args[2]
		}
		vehicleID, err := strconv.Atoi(targetArg)
		if err != nil || vehicleID <= 0 {
			output.ExitWithError("INVALID_ARGUMENT", "vehicle ID must be a positive integer", 0, output.ExitUsageError, cfg.JSONOutput)
		}
		commands.RunServiceList(ctx, cfg, apiClient, vehicleID, flagLimit)

	case "help":
		printUsage()

	default:
		output.ExitWithError("UNKNOWN_COMMAND", fmt.Sprintf("unknown command '%s'. Run 'lubelog --help' for usage.", cmd), 0, output.ExitUsageError, cfg.JSONOutput)
	}
}
