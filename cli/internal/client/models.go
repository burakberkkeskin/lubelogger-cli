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

// ServiceRecord represents a single service/maintenance event for a vehicle.
type ServiceRecord struct {
	ID          interface{} `json:"id"`
	VehicleID   interface{} `json:"vehicleId"`
	Date        string      `json:"date"`
	Odometer    interface{} `json:"odometer"`
	Description string      `json:"description"`
	Cost        interface{} `json:"cost"`
	Notes       *string     `json:"notes,omitempty"`
	Tags        interface{} `json:"tags,omitempty"`
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
