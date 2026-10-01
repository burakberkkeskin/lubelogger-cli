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
	"github.com/burakberk/lubelogger-cli/internal/version"
)

func printVersion(asJSON bool) {
	info := version.Get()
	if asJSON {
		_ = output.PrintJSON(info)
	} else {
		fmt.Println(info.String())
	}
}

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringSlice) Set(val string) error {
	*s = append(*s, val)
	return nil
}

func printUsage() {
	fmt.Println(`LubeLogger CLI (lubelog) - Agent-First Vehicle Maintenance Query & Record Tool

USAGE:
  lubelog [flags] <command> [arguments...]

COMMANDS:
  status                                 Check LubeLogger server health and authenticated identity
  vehicle list                           List all accessible vehicles (alias: vehicles)
  vehicle get <id>                       Get detailed info and aggregates for a vehicle (alias: vehicle <id>)
  service list <vehicleId>               List service/maintenance records for a vehicle (alias: service <id>)
  service add <vehicleId>                Add a new service record (supports --file attachments)
  service download <vehicleId> <recId>   Download attachments for a service record
  fuel list <vehicleId>                  List fuel/gas entries for a vehicle (alias: gas)
  fuel add <vehicleId>                   Add a new fuel/gas record (alias: gas add)
  fuel download <vehicleId> <recId>      Download attachments for a fuel record
  note list <vehicleId>                  List notes and document attachments for a vehicle (alias: note <id>)
  note add <vehicleId>                   Add a new note to a vehicle (supports --file attachments)
  note download <vehicleId> <noteId>     Download attachments for a note record
  file upload <path>...                  Upload local file(s) directly to server (alias: upload)
  file download <location>               Download file attachment directly by path (alias: download <loc>)
  version                                Display CLI version and build metadata

RECORD & MUTATION FLAGS:
  --desc, --description <str>  Service description or note title
  --odometer, --mileage <int>  Odometer / mileage reading (km or mi)
  --cost, --price <num>        Total cost / amount
  --date <YYYY-MM-DD>          Record date (default: today)
  --liters, --fuel <num>       Fuel consumed (for fuel add)
  --full                       Fill-up to full (default: true)
  --partial, --not-full        Partial fill-up (sets full=false)
  --missed                     Missed fuel-up (default: false)
  --title <str>                Title for note
  --text, --body <str>         Note body text
  --pinned                     Pin note to top of vehicle dashboard
  --notes <str>                Additional record notes
  --tags <str>                 Space-separated tags
  --file, -f <path>            Local file to attach (can be repeated)

GLOBAL FLAGS:
  --url <string>        Base URL of LubeLogger (default: $LUBELOGGER_URL or http://localhost:5000)
  --api-key <string>    API key for authentication (default: $LUBELOGGER_API_KEY)
  --username <string>   Username for Basic Auth fallback (default: $LUBELOGGER_USERNAME)
  --password <string>   Password for Basic Auth fallback (default: $LUBELOGGER_PASSWORD)
  --timeout <string>    Request timeout duration (default: $LUBELOGGER_TIMEOUT or 10s)
  --limit <int>         Limit number of records returned (e.g. --limit 3)
  --output, -o <path>   Output destination file or directory for downloads
  --dir <path>          Output directory for downloaded files
  --json                Output response in pure JSON format (default: non-TTY)
  -v, --version         Display CLI version
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
  lubelog service add 1 --desc "Oil Change" --odometer 55000 --cost 1200 --file ./receipt.pdf
  lubelog fuel add 1 --odometer 55300 --liters 42.5 --cost 1850 --full
  lubelog note add 1 --title "Tire Rotation" --text "Done at front shop" --pinned
  lubelog file upload ./invoice.pdf ./photo.jpg --json
  lubelog service download 1 14 --output ./downloads
  lubelog file download /documents/abc-123.pdf --output ./report.pdf`)
}

func parseArgs(args []string) (flags []string, positional []string) {
	valueFlags := map[string]bool{
		"--url": true, "-url": true,
		"--api-key": true, "-api-key": true,
		"--username": true, "-username": true,
		"--password": true, "-password": true,
		"--timeout": true, "-timeout": true,
		"--limit": true, "-limit": true,
		"--output": true, "-output": true,
		"-o": true,
		"--dir": true, "-dir": true,
		"--date": true, "-date": true,
		"--desc": true, "-desc": true,
		"--description": true, "-description": true,
		"--odometer": true, "-odometer": true,
		"--mileage": true, "-mileage": true,
		"--km": true, "-km": true,
		"--cost": true, "-cost": true,
		"--price": true, "-price": true,
		"--liters": true, "-liters": true,
		"--fuel": true, "-fuel": true,
		"--gallons": true, "-gallons": true,
		"--title": true, "-title": true,
		"--text": true, "-text": true,
		"--body": true, "-body": true,
		"--notes": true, "-notes": true,
		"--tags": true, "-tags": true,
		"--file": true, "-file": true,
		"-f": true,
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
	var flagOutput, flagDir string
	var flagDate, flagDesc, flagOdometer, flagCost, flagLiters string
	var flagTitle, flagText, flagNotes, flagTags string
	var flagFull, flagPartial, flagMissed, flagPinned bool
	var flagFiles stringSlice
	var flagLimit int
	var flagJSON, flagHelp, flagVersion bool

	flagArgs, posArgs := parseArgs(os.Args[1:])

	// Parse flags
	fs := flag.NewFlagSet("lubelog", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.StringVar(&flagURL, "url", "", "Base URL of LubeLogger")
	fs.StringVar(&flagAPIKey, "api-key", "", "API key")
	fs.StringVar(&flagUser, "username", "", "Username")
	fs.StringVar(&flagPass, "password", "", "Password")
	fs.StringVar(&flagTimeout, "timeout", "", "Request timeout")
	fs.StringVar(&flagOutput, "output", "", "Output destination file or directory")
	fs.StringVar(&flagOutput, "o", "", "Output destination file or directory")
	fs.StringVar(&flagDir, "dir", "", "Output directory for downloads")
	fs.IntVar(&flagLimit, "limit", 0, "Limit number of records returned")
	fs.BoolVar(&flagJSON, "json", false, "Force JSON output")
	fs.BoolVar(&flagVersion, "version", false, "Show version")
	fs.BoolVar(&flagVersion, "v", false, "Show version")
	fs.BoolVar(&flagHelp, "help", false, "Show help")
	fs.BoolVar(&flagHelp, "h", false, "Show help")

	// Mutation & record flags
	fs.StringVar(&flagDate, "date", "", "Record date (YYYY-MM-DD)")
	fs.StringVar(&flagDesc, "desc", "", "Description")
	fs.StringVar(&flagDesc, "description", "", "Description")
	fs.StringVar(&flagOdometer, "odometer", "", "Odometer reading")
	fs.StringVar(&flagOdometer, "mileage", "", "Odometer reading")
	fs.StringVar(&flagOdometer, "km", "", "Odometer reading")
	fs.StringVar(&flagCost, "cost", "", "Cost")
	fs.StringVar(&flagCost, "price", "", "Cost")
	fs.StringVar(&flagLiters, "liters", "", "Fuel consumed")
	fs.StringVar(&flagLiters, "fuel", "", "Fuel consumed")
	fs.StringVar(&flagLiters, "gallons", "", "Fuel consumed")
	fs.StringVar(&flagTitle, "title", "", "Note title")
	fs.StringVar(&flagText, "text", "", "Note body text")
	fs.StringVar(&flagText, "body", "", "Note body text")
	fs.StringVar(&flagNotes, "notes", "", "Additional notes")
	fs.StringVar(&flagTags, "tags", "", "Space-separated tags")
	fs.BoolVar(&flagFull, "full", true, "Is fill-up to full")
	fs.BoolVar(&flagPartial, "partial", false, "Is partial fill-up")
	fs.BoolVar(&flagPartial, "not-full", false, "Is partial fill-up")
	fs.BoolVar(&flagMissed, "missed", false, "Missed fuel up")
	fs.BoolVar(&flagPinned, "pinned", false, "Pin note")
	fs.Var(&flagFiles, "file", "File to attach (can be repeated)")
	fs.Var(&flagFiles, "f", "File to attach (can be repeated)")

	if err := fs.Parse(flagArgs); err != nil {
		output.ExitWithError("INVALID_FLAGS", err.Error(), 0, output.ExitUsageError, false)
	}

	if flagVersion {
		printVersion(flagJSON)
		os.Exit(output.ExitSuccess)
	}

	args := posArgs

	if flagHelp || len(args) == 0 {
		printUsage()
		os.Exit(output.ExitSuccess)
	}

	if flagOutput == "" && flagDir != "" {
		flagOutput = flagDir
	}

	cfg := config.LoadConfig(flagURL, flagAPIKey, flagUser, flagPass, flagTimeout, flagJSON)
	apiClient := client.NewClient(cfg.URL, cfg.APIKey, cfg.Username, cfg.Password, cfg.Timeout)
	ctx := context.Background()

	cmd := strings.ToLower(args[0])

	switch cmd {
	case "version":
		printVersion(cfg.JSONOutput)

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
			output.ExitWithError("MISSING_SUBCOMMAND", "expected 'service list <vehicleId>', 'service add <vehicleId>', or 'service download <vehicleId> <recordId>'", 0, output.ExitUsageError, cfg.JSONOutput)
		}
		subCmd := strings.ToLower(args[1])
		if subCmd == "add" {
			if len(args) < 3 {
				output.ExitWithError("MISSING_ARGUMENT", "usage: 'lubelog service add <vehicleId> [flags]'", 0, output.ExitUsageError, cfg.JSONOutput)
			}
			vID, err := strconv.Atoi(args[2])
			if err != nil || vID <= 0 {
				output.ExitWithError("INVALID_ARGUMENT", "vehicle ID must be a positive integer", 0, output.ExitUsageError, cfg.JSONOutput)
			}
			commands.RunServiceAdd(ctx, cfg, apiClient, vID, flagDesc, flagOdometer, flagCost, flagDate, flagNotes, flagTags, flagFiles)
			return
		}
		if subCmd == "download" {
			if len(args) < 4 {
				output.ExitWithError("MISSING_ARGUMENT", "usage: 'lubelog service download <vehicleId> <recordId>'", 0, output.ExitUsageError, cfg.JSONOutput)
			}
			vID, err := strconv.Atoi(args[2])
			if err != nil || vID <= 0 {
				output.ExitWithError("INVALID_ARGUMENT", "vehicle ID must be a positive integer", 0, output.ExitUsageError, cfg.JSONOutput)
			}
			rID, err := strconv.Atoi(args[3])
			if err != nil || rID <= 0 {
				output.ExitWithError("INVALID_ARGUMENT", "record ID must be a positive integer", 0, output.ExitUsageError, cfg.JSONOutput)
			}
			commands.RunServiceDownload(ctx, cfg, apiClient, vID, rID, flagOutput)
			return
		}
		targetArg := args[1]
		if subCmd == "list" || subCmd == "ls" {
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

	case "fuels", "fuel", "gas", "gasrecords":
		if len(args) < 2 {
			output.ExitWithError("MISSING_SUBCOMMAND", "expected 'fuel list <vehicleId>', 'fuel add <vehicleId>', or 'fuel download <vehicleId> <recordId>'", 0, output.ExitUsageError, cfg.JSONOutput)
		}
		subCmd := strings.ToLower(args[1])
		if subCmd == "add" {
			if len(args) < 3 {
				output.ExitWithError("MISSING_ARGUMENT", "usage: 'lubelog fuel add <vehicleId> [flags]'", 0, output.ExitUsageError, cfg.JSONOutput)
			}
			vID, err := strconv.Atoi(args[2])
			if err != nil || vID <= 0 {
				output.ExitWithError("INVALID_ARGUMENT", "vehicle ID must be a positive integer", 0, output.ExitUsageError, cfg.JSONOutput)
			}
			full := flagFull
			if flagPartial {
				full = false
			}
			commands.RunFuelAdd(ctx, cfg, apiClient, vID, flagOdometer, flagLiters, flagCost, flagDate, full, flagMissed, flagNotes, flagTags, flagFiles)
			return
		}
		if subCmd == "download" {
			if len(args) < 4 {
				output.ExitWithError("MISSING_ARGUMENT", "usage: 'lubelog fuel download <vehicleId> <recordId>'", 0, output.ExitUsageError, cfg.JSONOutput)
			}
			vID, err := strconv.Atoi(args[2])
			if err != nil || vID <= 0 {
				output.ExitWithError("INVALID_ARGUMENT", "vehicle ID must be a positive integer", 0, output.ExitUsageError, cfg.JSONOutput)
			}
			rID, err := strconv.Atoi(args[3])
			if err != nil || rID <= 0 {
				output.ExitWithError("INVALID_ARGUMENT", "record ID must be a positive integer", 0, output.ExitUsageError, cfg.JSONOutput)
			}
			commands.RunFuelDownload(ctx, cfg, apiClient, vID, rID, flagOutput)
			return
		}
		targetArg := args[1]
		if subCmd == "list" || subCmd == "ls" {
			if len(args) < 3 {
				output.ExitWithError("MISSING_ARGUMENT", "missing required vehicle ID: 'lubelog fuel list <vehicleId>'", 0, output.ExitUsageError, cfg.JSONOutput)
			}
			targetArg = args[2]
		}
		vehicleID, err := strconv.Atoi(targetArg)
		if err != nil || vehicleID <= 0 {
			output.ExitWithError("INVALID_ARGUMENT", "vehicle ID must be a positive integer", 0, output.ExitUsageError, cfg.JSONOutput)
		}
		commands.RunFuelList(ctx, cfg, apiClient, vehicleID, flagLimit)

	case "notes", "note":
		if len(args) < 2 {
			output.ExitWithError("MISSING_SUBCOMMAND", "expected 'note list <vehicleId>', 'note add <vehicleId>', or 'note download <vehicleId> <recordId>'", 0, output.ExitUsageError, cfg.JSONOutput)
		}
		subCmd := strings.ToLower(args[1])
		if subCmd == "add" {
			if len(args) < 3 {
				output.ExitWithError("MISSING_ARGUMENT", "usage: 'lubelog note add <vehicleId> [flags]'", 0, output.ExitUsageError, cfg.JSONOutput)
			}
			vID, err := strconv.Atoi(args[2])
			if err != nil || vID <= 0 {
				output.ExitWithError("INVALID_ARGUMENT", "vehicle ID must be a positive integer", 0, output.ExitUsageError, cfg.JSONOutput)
			}
			title := flagTitle
			if title == "" && flagDesc != "" {
				title = flagDesc
			}
			text := flagText
			if text == "" && flagNotes != "" {
				text = flagNotes
			}
			commands.RunNoteAdd(ctx, cfg, apiClient, vID, title, text, flagPinned, flagTags, flagFiles)
			return
		}
		if subCmd == "download" {
			if len(args) < 4 {
				output.ExitWithError("MISSING_ARGUMENT", "usage: 'lubelog note download <vehicleId> <recordId>'", 0, output.ExitUsageError, cfg.JSONOutput)
			}
			vID, err := strconv.Atoi(args[2])
			if err != nil || vID <= 0 {
				output.ExitWithError("INVALID_ARGUMENT", "vehicle ID must be a positive integer", 0, output.ExitUsageError, cfg.JSONOutput)
			}
			rID, err := strconv.Atoi(args[3])
			if err != nil || rID <= 0 {
				output.ExitWithError("INVALID_ARGUMENT", "record ID must be a positive integer", 0, output.ExitUsageError, cfg.JSONOutput)
			}
			commands.RunNoteDownload(ctx, cfg, apiClient, vID, rID, flagOutput)
			return
		}
		targetArg := args[1]
		if subCmd == "list" || subCmd == "ls" {
			if len(args) < 3 {
				output.ExitWithError("MISSING_ARGUMENT", "missing required vehicle ID: 'lubelog note list <vehicleId>'", 0, output.ExitUsageError, cfg.JSONOutput)
			}
			targetArg = args[2]
		}
		vehicleID, err := strconv.Atoi(targetArg)
		if err != nil || vehicleID <= 0 {
			output.ExitWithError("INVALID_ARGUMENT", "vehicle ID must be a positive integer", 0, output.ExitUsageError, cfg.JSONOutput)
		}
		commands.RunNoteList(ctx, cfg, apiClient, vehicleID)

	case "file", "document", "attachment":
		if len(args) < 2 {
			output.ExitWithError("MISSING_SUBCOMMAND", "expected 'file upload <path>...' or 'file download <location>'", 0, output.ExitUsageError, cfg.JSONOutput)
		}
		subCmd := strings.ToLower(args[1])
		if subCmd == "upload" || subCmd == "up" {
			var filesToUpload []string
			if len(args) > 2 {
				filesToUpload = append(filesToUpload, args[2:]...)
			}
			filesToUpload = append(filesToUpload, flagFiles...)
			commands.RunFileUpload(ctx, cfg, apiClient, filesToUpload)
			return
		}
		if subCmd == "download" || subCmd == "get" {
			if len(args) < 3 {
				output.ExitWithError("MISSING_ARGUMENT", "missing required file location: 'lubelog file download <location>'", 0, output.ExitUsageError, cfg.JSONOutput)
			}
			commands.RunFileDownload(ctx, cfg, apiClient, args[2], flagOutput)
		} else {
			commands.RunFileDownload(ctx, cfg, apiClient, args[1], flagOutput)
		}

	case "upload":
		var filesToUpload []string
		if len(args) > 1 {
			filesToUpload = append(filesToUpload, args[1:]...)
		}
		filesToUpload = append(filesToUpload, flagFiles...)
		commands.RunFileUpload(ctx, cfg, apiClient, filesToUpload)

	case "download":
		if len(args) < 2 {
			output.ExitWithError("MISSING_ARGUMENT", "missing file location: 'lubelog download <location>'", 0, output.ExitUsageError, cfg.JSONOutput)
		}
		commands.RunFileDownload(ctx, cfg, apiClient, args[1], flagOutput)

	case "help":
		printUsage()

	default:
		output.ExitWithError("UNKNOWN_COMMAND", fmt.Sprintf("unknown command '%s'. Run 'lubelog --help' for usage.", cmd), 0, output.ExitUsageError, cfg.JSONOutput)
	}
}

