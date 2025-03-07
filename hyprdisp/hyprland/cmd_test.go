package hyprland

import (
	"io"
	"os"
	"testing"
)

func Test_parseMonitorsPayload(t *testing.T) {
	var (
		file    *os.File
		data    []byte
		monitor Monitor
		err     error
	)

	file, err = os.Open("testdata/input.valid-monitors.txt")
	if err != nil {
		t.Fatalf("failed reading sample monitors.txt file: %v", err)
	}

	defer func() {
		var err error = file.Close()
		if err != nil {
			t.Fatalf("failed closing the file")
		}
	}()

	data, err = io.ReadAll(file)
	if err != nil {
		t.Fatalf("failed to read the file: %v", err)
	}

	monitors, err := parseMonitorsPayload(string(data))
	if err != nil {
		t.Fatalf("why am I returning an error???: %v", err)
	}

	if len(monitors) < 3 {
		t.Fatalf("expected 3 monitors to be extracted")
	}

	monitor = monitors[0]
	if monitor.ID != "0" ||
		monitor.Name != "DP-1" ||
		monitor.Resolution != "3440x1440@120.00000" ||
		monitor.Position != "0x1080" {
		t.Fatalf("monitor 0 has invalid values")
	}

	if monitor.Description != "Test Description 1" {
		t.Fatalf("monitor 0 has invalid description")
	}
	if monitor.Make != "Test Make 1" {
		t.Fatalf("monitor 0 has invalid make")
	}
	if monitor.Model != "Test Ultrawide Monitor" {
		t.Fatalf("monitor 0 has invalid model")
	}
	if monitor.Serial != "" {
		t.Fatalf("monitor 0 has invalid description")
	}
	if monitor.Enabled != true {
		t.Fatalf("monitor 0 has invalid enabled")
	}
}
