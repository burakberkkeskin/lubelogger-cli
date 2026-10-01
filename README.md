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

## Installation

### One-Line Install (`curl`)
Single copy-paste command for macOS and Linux. Directly resolves your OS (`uname -s`) and CPU architecture (`uname -m`):

```bash
sudo curl -fsSL "https://github.com/burakberkkeskin/lubelogger-cli/releases/latest/download/lubelog-$(uname -s)-$(uname -m)" -o /usr/local/bin/lubelog && sudo chmod +x /usr/local/bin/lubelog
```

*(Or install without sudo to `~/.local/bin`):*
```bash
curl -fsSL "https://github.com/burakberkkeskin/lubelogger-cli/releases/latest/download/lubelog-$(uname -s)-$(uname -m)" -o ~/.local/bin/lubelog && chmod +x ~/.local/bin/lubelog
```

> **Private Repository Access**: If downloading with the GitHub CLI:
> ```bash
> gh release download -R burakberkkeskin/lubelogger-cli --clobber -p "lubelog-$(uname -s)-$(uname -m)" -O /tmp/lubelog && sudo mv /tmp/lubelog /usr/local/bin/lubelog && sudo chmod +x /usr/local/bin/lubelog
> ```

### Build from Source (Local)
```bash
make build   # Builds cli/bin/lubelog embedding current git tag and commit
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

# Add a service record (with optional attachment)
./cli/bin/lubelog service add 1 --desc "Oil Change" --odometer 55000 --cost 1200 --file ./receipt.pdf --json

# Download attachments for a service record
./cli/bin/lubelog service download 1 14 --output ./downloads

# List fuel / gas entries
./cli/bin/lubelog fuel list 1 --limit 5 --json

# Add a fuel / gas record
./cli/bin/lubelog fuel add 1 --odometer 55300 --liters 42.5 --cost 1850 --full --json

# Download attachments for a fuel entry
./cli/bin/lubelog fuel download 1 40 --output ./downloads

# List notes, expert reports and attachments
./cli/bin/lubelog note list 1 --json

# Add a note (with optional attachment)
./cli/bin/lubelog note add 1 --title "Tire Rotation" --text "Done at front shop" --pinned --json

# Download attachments for a note
./cli/bin/lubelog note download 1 1 --output ./downloads

# Direct file upload to server
./cli/bin/lubelog file upload ./invoice.pdf ./photo.jpg --json

# Direct file download by server path/URL
./cli/bin/lubelog file download /documents/abc-123.pdf --output ./report.pdf

# Human-readable table output
./cli/bin/lubelog vehicle list
./cli/bin/lubelog service list 1 --limit 3
./cli/bin/lubelog fuel list 1
./cli/bin/lubelog note list 1

# Version information
./cli/bin/lubelog --version
./cli/bin/lubelog version --json
```

---

## Standardized Exit Codes

* `0`: Success
* `1`: General / unexpected failure
* `2`: Configuration or Authentication error (401/403)
* `3`: Resource Not Found (404)
* `4`: Network / Connectivity error
* `5`: Validation / Usage error

---

## Release Pipeline & Cross-Compilation

The repository features an automated GitHub Actions pipeline (`.github/workflows/release.yml`) that triggers **exclusively on version tag pushes** (`v*`).

### Target Architectures
Each release automatically builds pure Go static binaries (`CGO_ENABLED=0`) for:
* **Linux (x86_64)**: `linux/amd64`
* **Linux (ARM64)**: `linux/arm64`
* **macOS (Apple Silicon)**: `darwin/arm64`

Artifacts are packaged as `.tar.gz` archives with SHA-256 checksums (`checksums.txt`) and attached to the GitHub Release.

### How to Cut a Release

1. **Ensure Working Tree is Clean & Tests Pass:**
   ```bash
   cd cli && go test -v ./...
   ```

2. **Create an Annotated Git Tag:**
   ```bash
   git tag -a v0.1.0 -m "Release v0.1.0"
   ```

3. **Push the Tag to GitHub:**
   ```bash
   git push origin v0.1.0
   ```

4. **Monitor GitHub Actions & Download Release:**
   - The workflow runs unit tests, injects version metadata (`Version`, `Commit`, `Date`) via `ldflags`, cross-compiles all binaries, packages them, and creates the GitHub Release automatically under `https://github.com/burakberkkeskin/lubelogger-cli/releases`.

