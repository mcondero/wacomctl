package main

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseMonitors(t *testing.T) {
	output := `Monitors: 2
 0: +*HDMI-1 2560/597x1440/336+0+0 0
 1: +DP-1 1920/598x1080/336+0+0 0
`

	monitors, err := parseMonitors(output)
	if err != nil {
		t.Fatalf("parseMonitors returned error: %v", err)
	}

	want := []string{"HDMI-1", "DP-1"}
	if !reflect.DeepEqual(monitors, want) {
		t.Fatalf("parseMonitors() = %v, want %v", monitors, want)
	}
}

func TestBuildSelectionOptions(t *testing.T) {
	monitors := []string{"HDMI-1", "DP-1"}
	options := buildSelectionOptions(monitors)
	want := []string{"HDMI-1", "DP-1", "both", "turn off", "turn on"}
	if !reflect.DeepEqual(options, want) {
		t.Fatalf("buildSelectionOptions() = %v, want %v", options, want)
	}
}

func TestResolveWacomTargetName(t *testing.T) {
	resolved := resolveWacomTargetName("HDMI-0", 0)
	if resolved != "HEAD-0" {
		t.Fatalf("resolveWacomTargetName() = %q, want %q", resolved, "HEAD-0")
	}
}

func TestPersistedTargetRoundTrip(t *testing.T) {
	stateFile := filepath.Join(t.TempDir(), "state")
	t.Setenv("WACOMCTL_STATE_FILE", stateFile)

	if err := savePersistedTarget("HEAD-1"); err != nil {
		t.Fatalf("savePersistedTarget returned error: %v", err)
	}

	got, err := loadPersistedTarget()
	if err != nil {
		t.Fatalf("loadPersistedTarget returned error: %v", err)
	}

	if got != "HEAD-1" {
		t.Fatalf("loadPersistedTarget() = %q, want %q", got, "HEAD-1")
	}
}
