# Plan 2: Record Mutations (Add Service, Fuel, Note) & File Uploads

## 1. Overview & Objectives
Expand the `lubelog` CLI from a read/download query tool into a fully bidirectional, **Agent-First maintenance management interface**. This enables AI agents (and human operators) to programmatically register new service records, fuel logs, and notes, including attaching local documents (PDFs, invoices, images) via LubeLogger's REST API.

---

## 2. Live API Verification & Capabilities
LubeLogger REST API **fully supports** document uploads and record mutations out of the box using API key authentication (`x-api-key`):

1. **Document Upload Endpoint:**
   - `POST /api/documents/upload`
   - Content-Type: `multipart/form-data` with form field `documents` (one or multiple files).
   - Behavior: Writes files into LubeLogger's `data/documents/` storage with UUID-based names and returns `UploadedFiles[]` (`name`, `location`).
   - *Live Verification:* Verified live on `https://lubelogger.example.com`. Response: `[{"name":"...","location":"/documents/<UUID>.<ext>","isPending":false}]`.

2. **Service Record Creation:**
   - `POST /api/vehicle/servicerecords/add`
   - Payload: `GenericRecordExportModel` (`date`, `odometer`, `description`, `cost`, `notes`, `tags`, `files`).
   - Response: `{"success": true, "message": "Service Record Added", "data": {"recordId": <ID>}}`.

3. **Fuel/Gas Record Creation:**
   - `POST /api/vehicle/gasrecords/add`
   - Payload: `GasRecordExportModel` (`date`, `odometer`, `fuelConsumed`, `cost`, `isFillToFull`, `missedFuelUp`, `notes`, `tags`, `files`).
   - Response: `{"success": true, "message": "Gas Record Added", "data": {"recordId": <ID>}}`.

4. **Note Record Creation:**
   - `POST /api/vehicle/notes/add`
   - Payload: `NoteRecordExportModel` (`description`, `noteText`, `pinned`, `tags`, `files`).
   - Response: `{"success": true, "message": "Note Added", "data": {"recordId": <ID>}}`.

---

## 3. Architecture & Execution Workflow

### Two-Stage Atomic Upload & Link Pipeline

When an agent or human runs a command with one or more `--file` flags (e.g. `lubelog service add 3 --desc "..." --file ./invoice.pdf --file ./photo.jpg`):

```
┌─────────────────────────────────────────────────────────────┐
│ Step 1: Pre-flight Local Validation                         │
│ • Verify existence and readability of all local file paths. │
│ • Reject early (Exit 5) before hitting remote network.      │
└──────────────────────────────┬──────────────────────────────┘
                               │ (OK)
                               ▼
┌─────────────────────────────────────────────────────────────┐
│ Step 2: Multi-part Document Upload                          │
│ • Stream local files to POST /api/documents/upload.         │
│ • Receive remote file descriptors: [{name, location}].      │
└──────────────────────────────┬──────────────────────────────┘
                               │ (File UUIDs acquired)
                               ▼
┌─────────────────────────────────────────────────────────────┐
│ Step 3: Record Registration Payload                         │
│ • Assemble JSON payload with fields + uploaded files.       │
│ • Post to /api/vehicle/{servicerecords|gasrecords|notes}/add│
└──────────────────────────────┬──────────────────────────────┘
                               │ (Success)
                               ▼
┌─────────────────────────────────────────────────────────────┐
│ Step 4: Deterministic Output Emission                       │
│ • Emit JSON result envelope with new record ID & file paths │
│ • Table/Human summary with confirmation.                    │
└─────────────────────────────────────────────────────────────┘
```

---

## 4. CLI Command Ergonomics

### A. Service Record Addition
```bash
lubelog service add <vehicleId> \
  --desc <string> \
  --odometer <int> \
  --cost <number> \
  [--date YYYY-MM-DD] \
  [--notes <string>] \
  [--tags <string>] \
  [--file <path>...] \
  [--json]
```
- Defaults: `--date` defaults to today (`time.Now().Format("2006-01-02")`).
- Multi-file support: `--file ./fatura.pdf --file ./servis_fisi.jpg`.

### B. Fuel Record Addition
```bash
lubelog fuel add <vehicleId> \
  --odometer <int> \
  --liters <number> \
  --cost <number> \
  [--date YYYY-MM-DD] \
  [--full] \
  [--missed] \
  [--notes <string>] \
  [--tags <string>] \
  [--file <path>...] \
  [--json]
```
- Defaults: `--full=true` (most fill-ups are full), `--missed=false`, `--date` defaults to today.

### C. Note Addition
```bash
lubelog note add <vehicleId> \
  --title <string> \
  --text <string> \
  [--pinned] \
  [--tags <string>] \
  [--file <path>...] \
  [--json]
```

### D. Standalone File Upload
```bash
lubelog file upload <path>... [--json]
```
- Direct upload returning LubeLogger document locations (`/documents/<UUID>.<ext>`).

---

## 5. Standardized JSON Output Contract

When `--json` is supplied, all `add` commands return a deterministic envelope:

```json
{
  "success": true,
  "recordType": "service",
  "recordId": 19,
  "vehicleId": 3,
  "date": "2026-10-01",
  "odometer": 70500,
  "cost": 4500,
  "description": "Periyodik Bakım",
  "files": [
    {
      "name": "fatura.pdf",
      "location": "/documents/3f2d9885-a361-4216-bdba-686bb814ac6b.pdf"
    }
  ]
}
```

---

## 6. Implementation Milestones

- [ ] **Milestone 1: Client API Layer**
  - Add `UploadDocuments(ctx context.Context, filePaths []string) ([]UploadedFile, error)` using `multipart/writer`.
  - Add `AddServiceRecord(ctx context.Context, vehicleID int, req ServiceRecordInput) (int, error)`.
  - Add `AddGasRecord(ctx context.Context, vehicleID int, req GasRecordInput) (int, error)`.
  - Add `AddNoteRecord(ctx context.Context, vehicleID int, req NoteRecordInput) (int, error)`.
  - Comprehensive unit tests with `httptest.Server` verifying multipart body & JSON payloads.
- [ ] **Milestone 2: Command Implementations**
  - `commands/service_add.go`
  - `commands/fuel_add.go`
  - `commands/note_add.go`
  - `commands/file_upload.go`
- [ ] **Milestone 3: CLI Argument Routing & Flag Parsing**
  - Multi-value `--file` flag accumulator.
  - Sane defaults (auto-date to today, positive validations).
- [ ] **Milestone 4: Verification & Live Testing**
  - Test live creation and document attachment against homelab instance `https://lubelogger.example.com`.
  - Confirm newly created record is immediately queryable via `lubelog service list 3` and file downloadable via `lubelog service download 3 <newId>`.
- [ ] **Milestone 5: Documentation & Release**
  - Update `README.md` and `skills/lubelog-cli/SKILL.md`.
  - Tag `v0.2.0` and trigger automated cross-platform release pipeline.
