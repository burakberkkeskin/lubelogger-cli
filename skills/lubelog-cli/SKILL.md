---
name: lubelog-cli
description: CLI tool and query interface for LubeLogger vehicle maintenance platform. Use to list vehicles, inspect vehicle details, odometer records, service history, and health.
---

# LubeLogger CLI (`lubelog`)

The `lubelog` CLI is an agentic tool designed for fast, programmatic access to vehicle maintenance records stored in LubeLogger.

## Installation / Binary Location

The binary is compiled from `/Users/berk/Repos/lubelogger/cli` and can be built or executed directly:

```bash
# Direct build
cd /Users/berk/Repos/lubelogger/cli && go build -o bin/lubelog ./cmd/lubelog

# Invocation
/Users/berk/Repos/lubelogger/cli/bin/lubelog <command> [options]
```

## Environment Configuration

Configure the target server and credentials via environment variables or flags:

| Variable | Description | Default | Flag Equivalent |
| :--- | :--- | :--- | :--- |
| `LUBELOGGER_URL` | Base URL of LubeLogger instance | `http://localhost:5000` | `--url <url>` |
| `LUBELOGGER_API_KEY` | API Key for authentication | (None) | `--api-key <key>` |
| `LUBELOGGER_USERNAME`| Username (Basic Auth fallback) | (None) | `--username <user>`|
| `LUBELOGGER_PASSWORD`| Password (Basic Auth fallback) | (None) | `--password <pass>`|
| `LUBELOGGER_TIMEOUT` | Request timeout duration | `10s` | `--timeout <duration>` |
| `LUBELOGGER_OUTPUT_FORMAT` | Output format (`table` or `json`) | `json` (non-TTY) | `--json` / `--format <fmt>` |

---

## Commands Reference

### 1. `lubelog status`
Checks connection, server health, and authenticated user identity.

```bash
lubelog status --json
```

**JSON Output Example:**
```json
{
  "status": "healthy",
  "version": "1.4.2",
  "authenticated": true,
  "user": {
    "username": "admin",
    "isAdmin": true,
    "isRoot": true
  }
}
```

---

### 2. `lubelog vehicle list`
Lists all vehicles accessible to the authenticated account.

```bash
lubelog vehicle list --json
```

**JSON Output Example:**
```json
[
  {
    "id": 1,
    "year": 2021,
    "make": "Toyota",
    "model": "Corolla",
    "licensePlate": "34ABC123",
    "isElectric": false,
    "isDiesel": false,
    "tags": ["Daily"]
  }
]
```

---

### 3. `lubelog vehicle get <vehicleId>`
Retrieves comprehensive details, latest odometer reading, record counts, and urgent reminder metrics for a specific vehicle.

```bash
lubelog vehicle get 1 --json
```

**JSON Output Example:**
```json
{
  "vehicleData": {
    "id": 1,
    "year": 2021,
    "make": "Toyota",
    "model": "Corolla",
    "licensePlate": "34ABC123",
    "isElectric": false,
    "isDiesel": false,
    "tags": ["Daily"]
  },
  "lastReportedOdometer": 45200,
  "serviceRecordCount": 6,
  "serviceRecordCost": 1250.0,
  "repairRecordCount": 0,
  "repairRecordCost": 0.0,
  "gasRecordCount": 28,
  "gasRecordCost": 3450.0,
  "urgentReminderCount": 1,
  "pastDueReminderCount": 0
}
```

---

## Standardized Exit Codes

When invoking `lubelog`, check the process exit code for deterministic error handling:

| Exit Code | Meaning | Agent Action |
| :--- | :--- | :--- |
| `0` | Success | Parse stdout JSON payload |
| `1` | General error | Inspect `stderr` for message |
| `2` | Auth / Config error | Check `LUBELOGGER_URL` and `LUBELOGGER_API_KEY` |
| `3` | Resource Not Found | Target vehicle ID does not exist; verify ID |
| `4` | Network / Timeout | Server unreachable; check connectivity |
| `5` | Validation / Usage | Missing required argument (e.g. missing vehicle ID) |

When `--json` is supplied and an error occurs, `stderr` will contain:
```json
{
  "error": true,
  "code": "NOT_FOUND",
  "message": "Vehicle with ID 9999 was not found",
  "statusCode": 404
}
```
