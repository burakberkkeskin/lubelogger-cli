package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/burakberk/lubelogger-cli/internal/client"
)

// PrintJSON serializes an object into indented JSON on stdout.
func PrintJSON(v interface{}) error {
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(bytes))
	return nil
}

// FormatStatus prints the status response either as JSON or formatted text.
func FormatStatus(s *client.StatusResponse, asJSON bool) error {
	if asJSON {
		return PrintJSON(s)
	}

	fmt.Println("=== LubeLogger Status ===")
	fmt.Printf("Health:         %s\n", strings.ToUpper(s.Status))
	if s.ServerVersion != "" {
		fmt.Printf("Server Version: %s\n", s.ServerVersion)
	}
	fmt.Printf("Server URL:     %s\n", s.URL)
	if s.Authenticated {
		userStr := "Unknown"
		if s.User != nil && s.User.Username != "" {
			userStr = s.User.Username
		}
		fmt.Printf("Authentication: YES (User: %s)\n", userStr)
	} else {
		fmt.Println("Authentication: NO / Anonymous")
	}
	return nil
}

// FormatVehicles prints the vehicle list either as JSON or a formatted table.
func FormatVehicles(vehicles []client.Vehicle, asJSON bool) error {
	if asJSON {
		return PrintJSON(vehicles)
	}

	if len(vehicles) == 0 {
		fmt.Println("No vehicles found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tYEAR\tMAKE\tMODEL\tLICENSE PLATE\tTYPE\tTAGS")
	for _, v := range vehicles {
		fuelType := "Gas/Petrol"
		if v.IsElectric {
			fuelType = "Electric"
		} else if v.IsDiesel {
			fuelType = "Diesel"
		}

		tags := "-"
		if len(v.Tags) > 0 {
			tags = strings.Join(v.Tags, ", ")
		}

		plate := v.LicensePlate
		if plate == "" {
			plate = "-"
		}

		fmt.Fprintf(w, "%d\t%d\t%s\t%s\t%s\t%s\t%s\n",
			v.ID, v.Year, v.Make, v.Model, plate, fuelType, tags)
	}
	return w.Flush()
}

// FormatVehicleInfo prints single vehicle details either as JSON or a summary card.
func FormatVehicleInfo(info *client.VehicleInfo, asJSON bool) error {
	if asJSON {
		return PrintJSON(info)
	}

	v := info.VehicleData
	fmt.Printf("=== Vehicle #%d: %d %s %s ===\n", v.ID, v.Year, v.Make, v.Model)
	if v.LicensePlate != "" {
		fmt.Printf("License Plate:    %s\n", v.LicensePlate)
	}
	if len(v.Tags) > 0 {
		fmt.Printf("Tags:             %s\n", strings.Join(v.Tags, ", "))
	}
	if info.LastReportedOdometer > 0 {
		fmt.Printf("Odometer Reading: %d\n", info.LastReportedOdometer)
	}

	fmt.Println("\n--- Service & Cost Summary ---")
	fmt.Printf("Service Records:  %d  (Total Cost: $%.2f)\n", info.ServiceRecordCount, info.ServiceRecordCost)
	fmt.Printf("Repair Records:   %d  (Total Cost: $%.2f)\n", info.RepairRecordCount, info.RepairRecordCost)
	fmt.Printf("Upgrade Records:  %d  (Total Cost: $%.2f)\n", info.UpgradeRecordCount, info.UpgradeRecordCost)
	fmt.Printf("Fuel/Gas Records: %d  (Total Cost: $%.2f)\n", info.GasRecordCount, info.GasRecordCost)

	fmt.Println("\n--- Reminders & Status ---")
	fmt.Printf("Urgent Reminders:     %d\n", info.UrgentReminderCount)
	fmt.Printf("Past Due Reminders:   %d\n", info.PastDueReminderCount)
	if info.NextReminder != nil && info.NextReminder.Description != "" {
		fmt.Printf("Next Reminder:        %s [%s]\n", info.NextReminder.Description, info.NextReminder.Urgency)
	}

	return nil
}

// FormatServiceRecords prints a list of service records as JSON or a table.
func FormatServiceRecords(records []client.ServiceRecord, asJSON bool) error {
	if asJSON {
		return PrintJSON(records)
	}

	if len(records) == 0 {
		fmt.Println("No service records found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "ID\tDATE\tODOMETER\tCOST\tDESCRIPTION\tNOTES")
	for _, r := range records {
		notes := "-"
		if r.Notes != nil && *r.Notes != "" {
			cleanNote := strings.ReplaceAll(*r.Notes, "\n", " ")
			if len(cleanNote) > 40 {
				cleanNote = cleanNote[:37] + "..."
			}
			notes = cleanNote
		}
		fmt.Fprintf(w, "%v\t%s\t%v\t%v\t%s\t%s\n",
			r.ID, r.Date, r.Odometer, r.Cost, r.Description, notes)
	}
	return w.Flush()
}

// FormatNotes prints a list of vehicle notes as JSON or formatted text.
func FormatNotes(notes []client.NoteRecord, asJSON bool) error {
	if asJSON {
		return PrintJSON(notes)
	}

	if len(notes) == 0 {
		fmt.Println("No notes found for this vehicle.")
		return nil
	}

	fmt.Printf("=== Vehicle Notes (%d entries) ===\n\n", len(notes))
	for _, n := range notes {
		fmt.Printf("[%v] %s\n", n.ID, n.Description)
		if n.NoteText != "" && n.NoteText != n.Description {
			fmt.Printf("    Note:  %s\n", strings.ReplaceAll(n.NoteText, "\n", "\n    "))
		}
		if len(n.Files) > 0 {
			var fileNames []string
			for _, f := range n.Files {
				fileNames = append(fileNames, f.Name)
			}
			fmt.Printf("    Files: %s\n", strings.Join(fileNames, ", "))
		}
		fmt.Println()
	}
	return nil
}
