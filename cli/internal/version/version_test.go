package version

import (
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	info := Get()
	if info.Version == "" {
		t.Errorf("expected non-empty Version")
	}
	if info.GoVersion == "" {
		t.Errorf("expected non-empty GoVersion")
	}
	if info.OS == "" {
		t.Errorf("expected non-empty OS")
	}
	if info.Arch == "" {
		t.Errorf("expected non-empty Arch")
	}

	str := info.String()
	if !strings.Contains(str, "lubelog") {
		t.Errorf("expected String() to contain 'lubelog', got: %s", str)
	}
}
