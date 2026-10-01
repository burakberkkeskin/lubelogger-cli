# LubeLogger Agentic Workspace & CLI (`lubelog`)

An ultra-fast, minimalist, **Agent-First** command-line interface and tool integration layer for the [LubeLogger](https://github.com/hargata/lubelog) vehicle maintenance tracking platform.

## Architecture & Design Highlights

* **Agent-First Ergonomics (MCP Equivalent):** Designed specifically to be invoked by AI agents (Antigravity, Claude Code, Cursor, etc.) and sysadmins. Features deterministic JSON output schemas, zero interactive prompts, and standardized machine-readable exit codes.
* **Sub-10ms Execution:** Compiled static Go binary (`lubelog`) with zero runtime dependencies.
* **Zero Secrets:** Fully configurable via environment variables (`LUBELOGGER_URL`, `LUBELOGGER_API_KEY`) or CLI flags.

---

## Repository Structure

```
├── AGENTS.md                  # Monorepo and Agent operating guidelines
├── README.md                  # Workspace overview and quickstart
├── cli/                       # Pure Go CLI application
│   ├── cmd/lubelog/           # Main entrypoint and argument router
│   ├── internal/client/       # LubeLogger REST API client & DTOs
│   ├── internal/commands/     # Subcommand implementations (status, vehicle)
│   ├── internal/config/       # Configuration loader & flag/env resolver
│   └── internal/output/       # JSON / Tabular formatters and error envelopes
├── docs/agents/roles/         # Autonomous agent role specifications
│   ├── 0-pragmatic-product-owner.md
│   ├── 1-agentic-cli-architect.md
│   └── 2-lubelogger-api-specialist.md
├── plans/                     # Living task planning records
│   └── 1-agentic-infrastructure-and-simple-cli.md
└── skills/                    # Agent skill definitions
    └── lubelog-cli/SKILL.md   # Skill guide for LLM agents invoking lubelog
```

---

## Quickstart

### 1. Build the Binary
```bash
cd cli
go build -o bin/lubelog ./cmd/lubelog
```

### 2. Configure Environment
```bash
export LUBELOGGER_URL="https://lubelogger.burakberk.dev"  # or http://localhost:5000
export LUBELOGGER_API_KEY="<YOUR_API_KEY>"
```

### 3. Usage Examples
```bash
# Health & Status check
./cli/bin/lubelog status --json

# List all vehicles
./cli/bin/lubelog vehicle list --json

# Get details & aggregates for a vehicle
./cli/bin/lubelog vehicle get 1 --json

# List service/maintenance records (limit to last 3)
./cli/bin/lubelog service list 1 --limit 3 --json

# List notes, expert reports and attachments
./cli/bin/lubelog note list 1 --json

# Human-readable table output
./cli/bin/lubelog vehicle list
./cli/bin/lubelog service list 1 --limit 3
./cli/bin/lubelog note list 1
```

---

## Standardized Exit Codes

* `0`: Success
* `1`: General / unexpected failure
* `2`: Configuration or Authentication error (401/403)
* `3`: Resource Not Found (404)
* `4`: Network / Connectivity error
* `5`: Validation / Usage error
