package hyprland

import (
	"os"
	"path"
	"testing"
)

func Test_Apply(t *testing.T) {
	var (
		monitors   []Monitor
		workspaces []MonitorWorkspace
		err        error
	)

	monitors = []Monitor{
		{
			ID:         "1",
			Name:       "DP-1",
			Resolution: "3440x1440@165.00",
			Position:   "0x1080",
			Scale:      "1.00",
			Enabled:    true,
		},
		{
			ID:         "2",
			Name:       "DP-2",
			Resolution: "1920x1080@165.00",
			Position:   "760x0",
			Scale:      "1.00",
			Enabled:    true,
		},
	}

	workspaces = []MonitorWorkspace{
		{
			ID:         "1001",
			Monitor:    "DP-1",
			Default:    true,
			Persistent: true,
			Decorate:   true,
		},
		{
			ID:         "1002",
			Monitor:    "DP-1",
			Default:    false,
			Persistent: true,
			Decorate:   true,
		},
		{
			ID:         "2001",
			Monitor:    "DP-2",
			Default:    true,
			Persistent: true,
			Decorate:   true,
		},
	}

	err = Apply(monitors, workspaces)
	if err != nil {
		t.Fatalf("configuration application failure: %v", err)
	}

	var (
		actualPath string
		actualData []byte
		expectData []byte
	)

	actualPath, err = getConfigPath()
	if err != nil {
		t.Fatalf("failed to get configuration path: %v", err)
	}
	actualData, err = os.ReadFile(path.Join(
		actualPath,
		"actual.hypr-displays.conf",
	))
	if err != nil {
		t.Fatalf("failed reading output file: %v", err)
	}

	expectData, err = os.ReadFile("testdata/expect.valid-displays.conf")
	if err != nil {
		t.Fatalf("failed reading sample file: %v", err)
	}

	if string(expectData) != string(actualData) {
		t.Fatalf("expected configurations to be the same")
	}
}
