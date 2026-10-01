# AGENTS.md - LubeLogger CLI & Agentic Workspace Guidelines

> **Context**: This repository contains the **LubeLogger CLI (`lubelog`)** and the associated Agentic Infrastructure.
> **Audience**: AI Agents (Antigravity, Claude Code, Cursor, Copilot) and human infrastructure engineers.

---

## 1. Project Overview & Role

**LubeLogger CLI (`lubelog`)** is an ultra-fast, minimalist, **Agent-First** command-line interface and tool integration layer for the LubeLogger vehicle service and maintenance tracking platform.

*   **Primary User Persona:** AI Agents (Antigravity, Claude Code, Cursor, etc.) interacting with LubeLogger programmatically (functioning as a deterministic, lightweight alternative to MCP servers), alongside human sysadmins running terminal queries.
*   **Target Backend:** LubeLogger REST API (e.g. `https://lubelogger.burakberk.dev` or local `http://localhost:5000`).
*   **Core Invariants:**
    1.  **Agent-First Ergonomics:** Non-interactive execution, deterministic JSON payloads, and machine-readable exit codes.
    2.  **Sub-10ms Latency & Zero Dependencies:** Compiled single static Go binary (`lubelog`) requiring no external runtime daemons, interpreters, or virtualenvs.
    3.  **KISS & YAGNI:** Focused strictly on high-frequency, read/write maintenance actions (starting with vehicle listing and vehicle details). No premature OIDC or complex administrative bloat.

### Monorepo Structure

*   `cli/`: Go 1.26+ CLI codebase (`cmd/`, `internal/client`, `internal/commands`, `internal/config`, `internal/output`).
*   `plans/`: Central task planning, execution tracking, and sequential roadmap records (`1-xxx.md`).
*   `docs/agents/roles/`: Autonomous agent role specifications and guardrails.
*   `docs/architecture/`: Technical design documents, API contracts, and integration topologies.
*   `skills/`: Agent skill definitions (`skills/lubelog-cli/SKILL.md`) enabling direct tool discovery and invocation by AI agents.
*   `lubelog/`: Upstream LubeLogger source code clone (retained as local reference, gitignored).

### Agent Discovery & Blank-Context Protocol (Zero Source-Code Reading)

Incoming agents in a fresh context MUST NOT scan source code. Use the CLI commands via `skills/lubelog-cli/SKILL.md`:

| User Goal | Autonomous Agent Action | Example Command |
| :--- | :--- | :--- |
| **Check server & identity** | Verify connection | `lubelog status --json` |
| **Find vehicle ID** | List vehicles by plate / name | `lubelog vehicle list --json` |
| **Inspect vehicle summary** | Get odometer, reminder & costs | `lubelog vehicle get <id> --json` |
| **Fetch service records** | Get maintenance history / last N records | `lubelog service list <id> --limit <n> --json` |
| **Fetch notes & documents** | Get expert reports, notes & files | `lubelog note list <id> --json` |

---

## 2. Universal Architectural Principles

1.  **Strict Agent-First Standards (MCP-Equivalent CLI):**
    *   **Deterministic Machine Output:** Commands MUST support pure, unadorned JSON output (`--json` flag or default when output is redirected/piped).
    *   **Zero Interactive Prompts:** Commands MUST NOT hang or block on `stdin` prompts. If mandatory arguments or credentials are missing, fail immediately with an informative error message and an appropriate exit code.
    *   **Standardized Exit Codes:**
        *   `0`: Success
        *   `1`: General / unexpected failure
        *   `2`: Configuration or Authentication error (missing URL, invalid API key, 401/403)
        *   `3`: Resource Not Found (e.g., requested Vehicle ID does not exist, 404)
        *   `4`: Network / Connectivity error (server down, DNS failure, timeout)
        *   `5`: Validation / Usage error (missing or malformed arguments)
    *   **JSON Error Envelopes:** When JSON output is requested, errors to `stderr` or `stdout` MUST follow a structured schema:
        ```json
        {
          "error": true,
          "code": "NOT_FOUND",
          "message": "Vehicle with ID 10 not found",
          "statusCode": 404
        }
        ```

2.  **Configuration Precedence:**
    Configuration parameters MUST resolve in this strict order:
    1.  CLI command-line flags (`--url`, `--api-key`, `--output`, `--timeout`)
    2.  Environment variables (`LUBELOGGER_URL`, `LUBELOGGER_API_KEY`, `LUBELOGGER_OUTPUT_FORMAT`, `LUBELOGGER_TIMEOUT`)
    3.  Local configuration file (if present: `~/.config/lubelogger/config.json` or `.lubelogger.env`)
    4.  Built-in defaults (`http://localhost:5000`, `timeout: 10s`, format: `table` if TTY, `json` if piped)

3.  **Zero Secrets Policy (Non-Negotiable):**
    *   NEVER commit API keys, passwords, bearer tokens, or unencrypted `.env` files.
    *   Always refer to credentials via environment variables (`LUBELOGGER_API_KEY`).
    *   When documenting examples, sanitize all credentials using `<REDACTED>` or placeholder tokens.

4.  **Language Standard (STRICT):**
    *   **Codebase (100% English):** All source code, struct/variable/function names, code comments, commit messages, command descriptions, and error messages MUST be written strictly in **English**.

5.  **Critical Thinking & Pragmatic Scope Filter (Role 0):**
    *   Every feature proposal or refactoring MUST pass the **Pragmatic Product Owner (Patron)** filter.
    *   Push back aggressively against feature bloat, cosmetic rewrites, or unnecessary administrative features.

---

## 3. Planning & Task Tracking (Plans Directory)

*   **Mandatory Plan Files:** All tasks must have a dedicated plan file in `./plans/` following sequential naming: `1-agentic-infrastructure-and-simple-cli.md`, `2-service-records.md`, etc.
*   **Plan First:** Prioritize planning complex tasks thoroughly before writing code. A plan must specify goals, affected files, architecture decisions, task checklists, and verification methods.
*   **Live Progress Updates:** Agents must update the plan file in real time as items and checklist milestones are completed.

---

## 4. AI Agent Roles Directory (`docs/agents/roles/`)

Whenever an agent operates in a specific capacity, it **MUST** view and follow its corresponding role file:

*   [**0-pragmatic-product-owner.md**](docs/agents/roles/0-pragmatic-product-owner.md): **Mandatory First Gatekeeper.** Protects runway, enforces KISS/YAGNI, and cuts scope bloat.
*   [**1-agentic-cli-architect.md**](docs/agents/roles/1-agentic-cli-architect.md): Enforces CLI ergonomics, JSON schemas, exit codes, and non-interactive workflows.
*   [**2-lubelogger-api-specialist.md**](docs/agents/roles/2-lubelogger-api-specialist.md): Authoritative expert on LubeLogger REST API endpoints, DTO contracts, and authentication protocols.

---

## 5. The 5-Phase Agentic Lifecycle

Any agent working on this repository MUST operate within this unified 5-phase pipeline:

```
[Phase 1: Inception & Scope Trimming] ──> [Phase 2: API Contract & Client]
                                                      │
[Phase 5: Release & Git Commit]     <──  [Phase 4: Agent Skill & Verification] <── [Phase 3: CLI Subcommand & Formatter]
```

1.  **Phase 1 — Inception & Scope Trimming:** Evaluate requirement against KISS/YAGNI; create or update sequential plan in `plans/`.
2.  **Phase 2 — API Contract & Client:** Define domain models, request/response DTOs, and HTTP transport in `internal/client/`.
3.  **Phase 3 — CLI Subcommand & Formatter:** Implement subcommands, flag parsing, table/JSON formatting, and exit codes in `internal/commands/` and `internal/output/`.
4.  **Phase 4 — Agent Skill & Verification:** Update `skills/lubelog-cli/SKILL.md`; execute automated unit tests and end-to-end command verification.
5.  **Phase 5 — Release & Mandatory Git Commit:** Stage logical units and commit adhering to Conventional Commits format.

---

## 6. Task Completion & Git Commit Mandate

*   **Atomic & Mandatory Commits:** Whenever a task or plan milestone is completed and verified, agents MUST automatically stage and commit all relevant changes.
*   **Commit Message Convention:** Commit messages MUST be written strictly in **English** using Conventional Commits format (e.g. `feat(cli): implement vehicle list and get subcommands`, `docs(agents): define agentic roles and skill specifications`).
*   **Clean History:** Never leave finished work in a dirty working directory.
