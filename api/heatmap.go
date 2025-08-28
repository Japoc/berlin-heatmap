package handler

import (
	"berlin-heatmap/internal/server"
	"log"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {

	cfg := &server.Config{
		Port:       "8088",
		MatrixPath: "../data/artifacts/matrix_weekday_08.bin",
		GridFile:   "../data/artifacts/grid_links.json",
		StopsFile:  "../data/artifacts/stops.csv",
		K:          20,
	}

	store, err := server.NewHeatmapStore(cfg)
	if err != nil {
		log.Fatalf("failed to initialize store: %v", err)
	}

	switch r.URL.Path {
	case "/api/heatmap":
		store.HeatmapHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}
