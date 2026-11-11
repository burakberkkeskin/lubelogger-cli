---
name: lubelog-cli
description: CLI tool and query interface for LubeLogger vehicle maintenance platform. Use to list vehicles, inspect vehicle details, odometer records, service history, and health.
---

# LubeLogger CLI (`lubelog`)

The `lubelog` CLI is an agentic tool designed for fast, programmatic access to vehicle maintenance records stored in LubeLogger.

## Installation / Binary Location

The CLI binary is compiled from `/Users/berk/Repos/lubelogger/cli` and symlinked globally to `~/.local/bin/lubelog`:

```bash
# Global invocation (in system PATH)
lubelog <command> [options]

# Direct invocation / build from repository
cd /Users/berk/Repos/lubelogger/cli && go build -o bin/lubelog ./cmd/lubelog
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
    "description": "Battery replacement",
    "cost": 4000,
    "notes": "Replaced OEM battery with Tudor 60Ah battery.",
    "tags": ""
  },
  {
    "id": 17,
    "vehicleId": 3,
    "date": "2026-09-01",
    "odometer": 69528,
    "description": "Major maintenance",
    "cost": 16000,
    "notes": "- Timing belt\n- V belt\n- Water pump\n- Thermostat\n- Tensioner\n- Coolant",
    "tags": ""
  },
  {
    "id": 14,
    "vehicleId": 3,
    "date": "2026-04-22",
    "odometer": 64500,
    "description": "Scheduled service: Oil, air filter, oil filter, cabin filter",
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
    "description": "Pre-purchase inspection",
    "noteText": "Pre-purchase vehicle condition and inspection report.",
    "pinned": false,
    "files": [
      {
        "name": "inspection-report.pdf",
        "location": "/documents/7a22907b-03bb-4672-ace8-848607bc3d28.pdf"
      }
    ]
  }
]
```

---

### 6. `lubelog fuel list <vehicleId>` (alias: `gas`)
Lists fuel-up logs and efficiency history for a vehicle.

```bash
lubelog fuel list 3 --limit 5 --json
```

**JSON Output Example:**
```json
[
  {
    "id": 40,
    "vehicleId": 3,
    "date": "2025-09-22",
    "odometer": 56003,
    "fuelConsumed": 41.43,
    "cost": 2200,
    "fuelEconomy": 0,
    "isFillToFull": true,
    "missedFuelUp": false,
    "tags": "shell"
  }
]
```

---

### 7. Creating Records & Uploading Documents

AI agents can record maintenance, fuel fill-ups, and notes, including attaching local invoices/photos. Local files are pre-flight validated and automatically uploaded to LubeLogger's document storage before linking to the record.

#### A. Add Service Record
```bash
lubelog service add <vehicleId> \
  --desc <string> \
  --odometer <int> \
  --cost <number> \
  [--date YYYY-MM-DD] \
  [--notes <string>] \
  [--tags <string>] \
  [--file <path>...] \
  --json
```

**JSON Output Example:**
```json
{
  "success": true,
  "recordType": "service",
  "recordId": 19,
  "vehicleId": 3,
  "date": "2026-10-01",
  "odometer": "70500",
  "cost": "4500",
  "description": "Scheduled Service",
  "files": [
    {
      "name": "invoice.pdf",
      "location": "/documents/3f2d9885-a361-4216-bdba-686bb814ac6b.pdf"
    }
  ]
}
```

#### B. Add Fuel / Gas Record
```bash
lubelog fuel add <vehicleId> \
  --odometer <int> \
  --liters <number> \
  --cost <number> \
  [--date YYYY-MM-DD] \
  [--full] \
  [--partial] \
  [--missed] \
  [--notes <string>] \
  [--tags <string>] \
  [--file <path>...] \
  --json
```
* Note: `--full` is default (`true`). Use `--partial` or `--not-full` for non-complete fill-ups.

#### C. Add Note Record
```bash
lubelog note add <vehicleId> \
  --title <string> \
  --text <string> \
  [--pinned] \
  [--tags <string>] \
  [--file <path>...] \
  --json
```

#### D. Standalone File Upload
Upload files directly to LubeLogger's document storage:
```bash
lubelog file upload <path>... --json
```

**JSON Output Example:**
```json
[
  {
    "name": "invoice.pdf",
    "location": "/documents/3f2d9885-a361-4216-bdba-686bb814ac6b.pdf"
  }
]
```

---

### 8. Downloading Attachments & Documents
Download attached files (invoices, receipts, inspection PDFs, photos) from service records, fuel entries, notes, or directly by file path:

```bash
# Download all attachments from a service record
lubelog service download <vehicleId> <recordId> [--output <dir>] --json

# Download all attachments from a fuel record
lubelog fuel download <vehicleId> <recordId> [--output <dir>] --json

# Download all attachments from a note
lubelog note download <vehicleId> <noteId> [--output <dir>] --json

# Direct download by server location / URL
lubelog file download /documents/7a22907b-03bb-4672-ace8-848607bc3d28.pdf --output ./inspection.pdf --json
```

**JSON Output Example:**
```json
[
  {
    "name": "inspection-report.pdf",
    "location": "/documents/7a22907b-03bb-4672-ace8-848607bc3d28.pdf",
    "savedTo": "./downloads/inspection-report.pdf",
    "bytes": 10674297,
    "contentType": "application/pdf"
  }
]
```

---

### 9. `lubelog version` (or `-v`, `--version`)
Displays CLI version, build timestamp, commit SHA, and Go/OS/Arch metadata.

```bash
lubelog version --json
```

**JSON Output Example:**
```json
{
  "version": "v0.1.0",
  "commit": "b7d530b",
  "date": "2026-10-01T07:40:09Z",
  "goVersion": "go1.23.0",
  "os": "darwin",
  "arch": "arm64"
}
```

---

## Autonomous Agent Execution Protocols (Blank-Context Guide)

When an AI Agent is asked to query or record vehicle maintenance, follow this exact protocol without inspecting the codebase:

1. **Step 1: Discover Vehicles & Resolve ID**
   * Run `lubelog vehicle list --json`
   * Match the user's vehicle name, plate, or model (e.g. "Egea" -> `id: 3`).
2. **Step 2: Fetch Target Information, Download Documents, or Add Records**
   * If the user asks for vehicle overview/mileage: `lubelog vehicle get <id> --json`
   * If the user asks for service records: `lubelog service list <id> [--limit <n>] --json`
   * If the user asks for fuel/gas history: `lubelog fuel list <id> [--limit <n>] --json`
   * If the user asks for notes, documents, or expert reports: `lubelog note list <id> --json`
   * If the user asks to record maintenance/service:
     `lubelog service add <id> --desc "..." --odometer <n> --cost <n> [--file <path>...] --json`
   * If the user asks to record fuel:
     `lubelog fuel add <id> --odometer <n> --liters <n> --cost <n> [--file <path>...] --json`
   * If the user asks to record a note:
     `lubelog note add <id> --title "..." --text "..." [--file <path>...] --json`
   * If the user asks to download an invoice, report, or file:
     - From service record: `lubelog service download <id> <serviceId> [--output <dir>] --json`
     - From fuel record: `lubelog fuel download <id> <fuelId> [--output <dir>] --json`
     - From note: `lubelog note download <id> <noteId> [--output <dir>] --json`
     - Direct by path/URL: `lubelog file download <location> [--output <path>] --json`
3. **Step 3: Parse & Present**
   * Parse the pure JSON array or mutation object. Format the dates, costs, odometers, notes, and file names cleanly for the user.

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
