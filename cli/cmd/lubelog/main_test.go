package main

import (
	"reflect"
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name         string
		input        []string
		wantFlags    []string
		wantPosArgs  []string
	}{
		{
			name:        "flags at end",
			input:       []string{"vehicle", "list", "--json"},
			wantFlags:   []string{"--json"},
			wantPosArgs: []string{"vehicle", "list"},
		},
		{
			name:        "flags at beginning",
			input:       []string{"--json", "vehicle", "list"},
			wantFlags:   []string{"--json"},
			wantPosArgs: []string{"vehicle", "list"},
		},
		{
			name:        "value flag with space",
			input:       []string{"--url", "http://localhost:5000", "vehicle", "get", "1", "--json"},
			wantFlags:   []string{"--url", "http://localhost:5000", "--json"},
			wantPosArgs: []string{"vehicle", "get", "1"},
		},
		{
			name:        "value flag with equals",
			input:       []string{"--url=http://localhost:5000", "status"},
			wantFlags:   []string{"--url=http://localhost:5000"},
			wantPosArgs: []string{"status"},
		},
		{
			name:        "mutation flags with multiple files",
			input:       []string{"service", "add", "1", "--desc", "Oil Change", "--odometer", "50000", "--cost", "1200", "--file", "receipt.pdf", "--file", "photo.jpg", "--json"},
			wantFlags:   []string{"--desc", "Oil Change", "--odometer", "50000", "--cost", "1200", "--file", "receipt.pdf", "--file", "photo.jpg", "--json"},
			wantPosArgs: []string{"service", "add", "1"},
		},
		{
			name:        "fuel add with flags and boolean",
			input:       []string{"fuel", "add", "2", "--odometer", "60000", "--liters", "45", "--cost", "2000", "--full"},
			wantFlags:   []string{"--odometer", "60000", "--liters", "45", "--cost", "2000", "--full"},
			wantPosArgs: []string{"fuel", "add", "2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFlags, gotPos := parseArgs(tt.input)
			if !reflect.DeepEqual(gotFlags, tt.wantFlags) {
				t.Errorf("flags mismatch: got %v, want %v", gotFlags, tt.wantFlags)
			}
			if !reflect.DeepEqual(gotPos, tt.wantPosArgs) {
				t.Errorf("positional mismatch: got %v, want %v", gotPos, tt.wantPosArgs)
			}
		})
	}
}
