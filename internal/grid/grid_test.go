package grid_test

import (
	"berlin-heatmap/internal/grid"
	"testing"

	"berlin-heatmap/internal/gql"
)

func TestHaversine(t *testing.T) {
	// Distance between same point should be 0
	d := grid.Haversine(52.5, 13.4, 52.5, 13.4)
	if d != 0 {
		t.Errorf("expected 0, got %f", d)
	}
}

func TestNearestKStops(t *testing.T) {
	stops := []gql.Stop{
		{ID: "A", Lat: 52.5, Lon: 13.4},
		{ID: "B", Lat: 52.6, Lon: 13.5},
	}
	idxs, walks := grid.NearestKStops(52.55, 13.45, stops, 1)
	if len(idxs) != 1 || len(walks) != 1 {
		t.Fatalf("expected 1 nearest stop, got %d", len(idxs))
	}
}
