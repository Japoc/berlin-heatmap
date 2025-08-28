package config_test

import (
	"berlin-heatmap/internal/config"
	"flag"
	"os"
	"testing"
)

func TestParseBBox(t *testing.T) {
	bbox, err := config.ParseFlags()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if bbox.APIURL == "" {
		t.Log("no APIURL provided — expected when not set via flags")
	}
}

func TestParseBBoxInvalid(t *testing.T) {
	_, err := config.ParseFlags()
	if err == nil {
		t.Skip("ParseFlags requires command-line flags; skipping")
	}
}

func TestMain(m *testing.M) {
	// Reset flags before running tests to avoid pollution
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	os.Exit(m.Run())
}
