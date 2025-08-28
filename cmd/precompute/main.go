package main

import (
	"berlin-heatmap/internal/config"
	"berlin-heatmap/internal/grid"
	"berlin-heatmap/internal/matrix"
	"berlin-heatmap/internal/stops"
	"context"
	"log"
)

func main() {
	cfg := config.ParseFlags()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Fetch stops
	stopsList := stops.FetchStops(ctx, cfg)
	stops.WriteStopsCSV(stopsList, cfg.OutDir)

	// 2. Build stop→stop matrix
	matrixData := matrix.BuildTravelMatrix(ctx, stopsList, cfg)
	matrix.WriteMatrix(matrixData, cfg)

	// 3. Build grid links
	gridCells := grid.BuildGrid(stopsList, cfg)
	grid.WriteGrid(gridCells, cfg)

	log.Printf("done. artifacts in %s", cfg.OutDir)
}
