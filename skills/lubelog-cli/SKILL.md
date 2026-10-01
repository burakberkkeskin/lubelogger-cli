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

### 4. `lubelog service list <vehicleId> [--limit <n>]`
Retrieves chronological service and maintenance records for a vehicle (sorted newest first). Use `--limit <n>` to retrieve only the last N records.

```bash
lubelog service list 3 --limit 3 --json
```

**JSON Output Example:**
```json
[
  {
    "id": 18,
    "vehicleId": 3,
    "date": "2026-09-19",
    "odometer": 69779,
    "description": "Akü değişimi",
    "cost": 4000,
    "notes": "Üstündeki orijinal akü, Tudor 60A akü ile değişti.",
    "tags": ""
  },
  {
    "id": 17,
    "vehicleId": 3,
    "date": "2026-09-01",
    "odometer": 69528,
    "description": "70 bin ağır bakımı",
    "cost": 16000,
    "notes": "- Triger kayışı (orjinal)\n- V kayışı (orijinal)\n- Devirdaim (orijinal)\n- Termostat\n- Gergi (Orijinal)\n- Antifriz\n- Saf su",
    "tags": ""
  },
  {
    "id": 14,
    "vehicleId": 3,
    "date": "2026-04-22",
    "odometer": 64500,
    "description": "Periyodik bakım: Yağ, hava filtre, yağ filtre, polen filtre",
    "cost": 5000,
    "tags": ""
  }
]
```

---

### 5. `lubelog note list <vehicleId>`
Retrieves notes, reminders, expert reports, and attached files/documents for a vehicle.

```bash
lubelog note list 3 --json
```

**JSON Output Example:**
```json
[
  {
    "id": 1,
    "vehicleId": 3,
    "description": "Vavacars Tarcanlar Ekspertiz",
    "noteText": "Aracın Vavacars'dan satın alınırkenki 6.000 KM Tarcanlar ekspertiz raporu.",
    "pinned": false,
    "files": [
      {
        "name": "tarcanlar-ekspertiz.pdf",
        "location": "/documents/7a22907b-03bb-4672-ace8-848607bc3d28.pdf"
      }
    ]
  }
]
```

---

## Autonomous Agent Execution Protocols (Blank-Context Guide)

When an AI Agent is asked a question regarding vehicles or maintenance, follow this exact 3-step protocol without inspecting the codebase:

1. **Step 1: Discover Vehicles & Resolve ID**
   * Run `lubelog vehicle list --json`
   * Match the user's vehicle name, plate, or model (e.g. "Egea" -> `id: 3`).
2. **Step 2: Fetch Target Information**
   * If the user asks for vehicle overview/mileage: `lubelog vehicle get <id> --json`
   * If the user asks for service records: `lubelog service list <id> [--limit <n>] --json`
   * If the user asks for notes, documents, or expert reports: `lubelog note list <id> --json`
3. **Step 3: Parse & Present**
   * Parse the pure JSON array. Format the dates, costs, odometers, notes, and file names cleanly for the user.

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
