# Plan 3: Global Agent Skill Integration & Universal Bridge

## 1. Overview & Objectives
Transform the LubeLogger CLI and vehicle maintenance operations into a first-class, **Universal Agent Skill (`lubelogger-control`)** hosted centrally in `agent-profile`. This enables any AI agent (Antigravity, Claude Code, Cursor) or human engineer on this machine to query vehicle health, maintenance history, fuel records, and upload/download invoices from any workspace or directory.

---

## 2. Key Architectural Decisions
1. **Centralized Global Skill SSoT:**
   - Authored in [`/Users/berk/Repos/agent-profile/skills/lubelogger-control/SKILL.md`](file:///Users/berk/Repos/agent-profile/skills/lubelogger-control/SKILL.md).
   - Automatically exposed to Antigravity / Gemini via global symlink `~/.gemini/config/skills/lubelogger-control/SKILL.md`.
   - Adheres to the exact progressive disclosure and SSoT principles established by `homeassistant-control` and `homelab-control`.
2. **Global Executable in PATH:**
   - Symlinked `/Users/berk/Repos/lubelogger/cli/bin/lubelog` to `/Users/berk/.local/bin/lubelog`.
   - Directly invokable as `lubelog` across all shells and environments without relative paths.
3. **Deterministic Test Isolation:**
   - Resolved environment leakage in `cli/internal/config/config_test.go` so tests cleanly isolate ambient `LUBELOGGER_*` environment variables.

---

## 3. Checklist & Milestones

- [x] Fix environment variable isolation in `cli/internal/config/config_test.go`.
- [x] Run full automated test suite (`go test -v ./...`) across all packages.
- [x] Recompile optimized static binary (`cli/bin/lubelog`).
- [x] Create global executable symlink in `~/.local/bin/lubelog`.
- [x] Verify global execution from outside the workspace (`cd /tmp && lubelog status --json`).
- [x] Author universal agent skill in `agent-profile/skills/lubelogger-control/SKILL.md`.
- [x] Update `agent-profile/README.md` cataloging the new skill.
- [x] Update local skill guide in `skills/lubelog-cli/SKILL.md`.
- [x] Commit atomic changes in both repositories following Conventional Commits.
