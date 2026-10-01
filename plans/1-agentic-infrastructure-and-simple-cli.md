# Plan 1: Agentic Infrastructure & Minimalist LubeLogger CLI

## 1. Overview & Objectives
Establish the foundation of the LubeLogger CLI (`lubelog`) and agentic workflow. The CLI will serve as an **Agent-First, MCP-equivalent interface** that an AI Agent or developer can invoke to query vehicles and vehicle maintenance states with deterministic JSON output, sub-10ms latency, and structured error reporting.

---

## 2. Key Architecture Decisions
1. **Language & Runtime:** Pure Go (1.26+) using only standard library packages (`net/http`, `encoding/json`, `flag`, `text/tabwriter`) for a single static binary with zero external dependencies and sub-10ms startup.
2. **Configuration Precedence:** Flags (`--url`, `--api-key`, `--json`) > Env Vars (`LUBELOGGER_URL`, `LUBELOGGER_API_KEY`, `LUBELOGGER_OUTPUT_FORMAT`) > Defaults.
3. **Agent-First Ergonomics:** Pure JSON mode via `--json` or redirected stdout, structured error envelopes on non-zero exit, no interactive stdin blocking.
4. **Scope (Phase 1):** `status` (ping/health/whoami), `vehicle list`, and `vehicle get <id>`.

---

## 3. Checklist & Milestones

- [x] Initialize Git repository and `.gitignore` in `/Users/berk/Repos/lubelogger`.
- [x] Author root `AGENTS.md` specifying universal principles, agent lifecycle, and CLI standards.
- [x] Create agent roles catalog in `docs/agents/roles/`:
  - [x] `0-pragmatic-product-owner.md`
  - [x] `1-agentic-cli-architect.md`
  - [x] `2-lubelogger-api-specialist.md`
- [x] Implement Go CLI in `cli/`:
  - [x] `go.mod` module setup (`github.com/burakberk/lubelogger-cli`)
  - [x] `internal/client`: HTTP client with API Key & Basic Auth, invariant headers, DTO models.
  - [x] `internal/config`: Env var and flag resolver.
  - [x] `internal/output`: JSON and tabular text formatters, structured error envelopes and exit codes.
  - [x] `internal/commands`: `status`, `vehicle list`, `vehicle get` subcommands.
  - [x] `cmd/lubelog/main.go`: Entrypoint and command router with detailed `--help`.
- [x] Author Agent Skill definition in `skills/lubelog-cli/SKILL.md`.
- [x] Add automated unit tests for client and output formatters (`go test ./...`).
- [x] Build static binary and verify command execution.
- [ ] Commit all changes with Conventional Commits.
