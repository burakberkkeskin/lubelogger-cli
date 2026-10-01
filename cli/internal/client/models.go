package client

// Vehicle represents the core vehicle record in LubeLogger.
type Vehicle struct {
	ID            int      `json:"id"`
	Year          int      `json:"year"`
	Make          string   `json:"make"`
	Model         string   `json:"model"`
	LicensePlate  string   `json:"licensePlate"`
	PurchaseDate  string   `json:"purchaseDate,omitempty"`
	SoldDate      string   `json:"soldDate,omitempty"`
	PurchasePrice float64  `json:"purchasePrice,omitempty"`
	SoldPrice     float64  `json:"soldPrice,omitempty"`
	IsElectric    bool     `json:"isElectric"`
	IsDiesel      bool     `json:"isDiesel"`
	UseHours      bool     `json:"useHours"`
	Tags          []string `json:"tags,omitempty"`
	ImageLocation string   `json:"imageLocation,omitempty"`
}

// NextReminder represents the earliest upcoming reminder on a vehicle.
type NextReminder struct {
	Description string `json:"description,omitempty"`
	Urgency     string `json:"urgency,omitempty"`
	Metric      string `json:"metric,omitempty"`
}

// UploadedFile represents a file attachment (invoice, report, image) on any record.
type UploadedFile struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

// NoteFile is an alias for UploadedFile.
type NoteFile = UploadedFile

// ServiceRecord represents a single service/maintenance event for a vehicle.
type ServiceRecord struct {
	ID          interface{}    `json:"id"`
	VehicleID   interface{}    `json:"vehicleId"`
	Date        string         `json:"date"`
	Odometer    interface{}    `json:"odometer"`
	Description string         `json:"description"`
	Cost        interface{}    `json:"cost"`
	Notes       *string        `json:"notes,omitempty"`
	Tags        interface{}    `json:"tags,omitempty"`
	Files       []UploadedFile `json:"files,omitempty"`
}

// GasRecord represents a fuel-up entry for a vehicle.
type GasRecord struct {
	ID           interface{}    `json:"id"`
	VehicleID    interface{}    `json:"vehicleId"`
	Date         string         `json:"date"`
	Odometer     interface{}    `json:"odometer"`
	FuelConsumed interface{}    `json:"fuelConsumed"`
	Cost         interface{}    `json:"cost"`
	FuelEconomy  interface{}    `json:"fuelEconomy,omitempty"`
	IsFillToFull interface{}    `json:"isFillToFull,omitempty"`
	MissedFuelUp interface{}    `json:"missedFuelUp,omitempty"`
	Notes        *string        `json:"notes,omitempty"`
	Tags         interface{}    `json:"tags,omitempty"`
	Files        []UploadedFile `json:"files,omitempty"`
}

// NoteRecord represents a freeform note or document entry for a vehicle.
type NoteRecord struct {
	ID          interface{}    `json:"id"`
	VehicleID   interface{}    `json:"vehicleId"`
	Description string         `json:"description"`
	NoteText    string         `json:"noteText"`
	Pinned      interface{}    `json:"pinned,omitempty"`
	Tags        interface{}    `json:"tags,omitempty"`
	Files       []UploadedFile `json:"files,omitempty"`
}

// DownloadedFile represents the outcome of a downloaded file attachment.
type DownloadedFile struct {
	Name        string `json:"name"`
	Location    string `json:"location"`
	SavedTo     string `json:"savedTo"`
	Bytes       int64  `json:"bytes"`
	ContentType string `json:"contentType,omitempty"`
}

// VehicleInfo represents aggregated statistics and state for a vehicle.
type VehicleInfo struct {
	VehicleData             Vehicle       `json:"vehicleData"`
	LastReportedOdometer    int           `json:"lastReportedOdometer"`
	ServiceRecordCount      int           `json:"serviceRecordCount"`
	ServiceRecordCost       float64       `json:"serviceRecordCost"`
	RepairRecordCount       int           `json:"repairRecordCount"`
	RepairRecordCost        float64       `json:"repairRecordCost"`
	UpgradeRecordCount      int           `json:"upgradeRecordCount"`
	UpgradeRecordCost       float64       `json:"upgradeRecordCost"`
	GasRecordCount          int           `json:"gasRecordCount"`
	GasRecordCost           float64       `json:"gasRecordCost"`
	TaxRecordCount          int           `json:"taxRecordCount"`
	TaxRecordCost           float64       `json:"taxRecordCost"`
	VeryUrgentReminderCount int           `json:"veryUrgentReminderCount"`
	UrgentReminderCount     int           `json:"urgentReminderCount"`
	NotUrgentReminderCount  int           `json:"notUrgentReminderCount"`
	PastDueReminderCount    int           `json:"pastDueReminderCount"`
	NextReminder            *NextReminder `json:"nextReminder,omitempty"`
	PlanRecordBackLogCount  int           `json:"planRecordBackLogCount"`
	PlanRecordInProgressCount int         `json:"planRecordInProgressCount"`
	PlanRecordDoneCount     int           `json:"planRecordDoneCount"`
}

// UserExportModel represents identity information from /api/whoami.
type UserExportModel struct {
	Username     string      `json:"username"`
	EmailAddress string      `json:"emailAddress"`
	IsAdmin      interface{} `json:"isAdmin"`
	IsRoot       interface{} `json:"isRoot"`
}

// DatabaseHealth represents health probe status for the underlying DB.
type DatabaseHealth struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// ServerHealth represents the payload returned by /health.
type ServerHealth struct {
	Version       string           `json:"version"`
	Status        string           `json:"status"`
	TotalDuration string           `json:"totalDuration,omitempty"`
	Checks        []DatabaseHealth `json:"checks"`
}

// StatusResponse is the composite response returned by `lubelog status`.
type StatusResponse struct {
	Status        string           `json:"status"`
	ServerVersion string           `json:"serverVersion,omitempty"`
	URL           string           `json:"url"`
	Authenticated bool             `json:"authenticated"`
	User          *UserExportModel `json:"user,omitempty"`
}

// ServiceRecordInput represents the payload to create a new service record.
type ServiceRecordInput struct {
	Date        string         `json:"date"`
	Odometer    string         `json:"odometer"`
	Description string         `json:"description"`
	Cost        string         `json:"cost"`
	Notes       string         `json:"notes,omitempty"`
	Tags        string         `json:"tags,omitempty"`
	Files       []UploadedFile `json:"files,omitempty"`
}

// GasRecordInput represents the payload to create a new fuel/gas record.
type GasRecordInput struct {
	Date         string         `json:"date"`
	Odometer     string         `json:"odometer"`
	FuelConsumed string         `json:"fuelConsumed"`
	Cost         string         `json:"cost"`
	IsFillToFull string         `json:"isFillToFull"`
	MissedFuelUp string         `json:"missedFuelUp"`
	Notes        string         `json:"notes,omitempty"`
	Tags         string         `json:"tags,omitempty"`
	Files        []UploadedFile `json:"files,omitempty"`
}

// NoteRecordInput represents the payload to create a new note record.
type NoteRecordInput struct {
	Description string         `json:"description"`
	NoteText    string         `json:"noteText"`
	Pinned      string         `json:"pinned,omitempty"`
	Tags        string         `json:"tags,omitempty"`
	Files       []UploadedFile `json:"files,omitempty"`
}

// OperationResponse represents the generic response envelope from LubeLogger mutation APIs.
type OperationResponse struct {
	Success        bool                   `json:"success"`
	Message        string                 `json:"message"`
	AdditionalData map[string]interface{} `json:"additionalData,omitempty"`
}

// MutationResult is the standardized CLI output returned when a record is created.
type MutationResult struct {
	Success      bool           `json:"success"`
	RecordType   string         `json:"recordType"`
	RecordID     interface{}    `json:"recordId"`
	VehicleID    int            `json:"vehicleId"`
	Date         string         `json:"date,omitempty"`
	Odometer     interface{}    `json:"odometer,omitempty"`
	Cost         interface{}    `json:"cost,omitempty"`
	Description  string         `json:"description,omitempty"`
	FuelConsumed interface{}    `json:"fuelConsumed,omitempty"`
	Files        []UploadedFile `json:"files,omitempty"`
}

