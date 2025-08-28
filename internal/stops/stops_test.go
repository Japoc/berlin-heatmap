package stops_test

import (
	"berlin-heatmap/internal/stops"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteStopsCSV(t *testing.T) {
	tmp := t.TempDir()
	testStops := []stops.Stop{
		{ID: "A", Lat: 52.5, Lon: 13.4},
		{ID: "B", Lat: 52.6, Lon: 13.5},
	}
	stops.WriteStopsCSV(testStops, tmp)

	path := filepath.Join(tmp, "stops.csv")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected file at %s, got error: %v", path, err)
	}
}
