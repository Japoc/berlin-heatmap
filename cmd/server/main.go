package main

import (
	server2 "berlin-heatmap/internal/server"
	"log"
)

func main() {
	cfg := server2.ParseFlags()

	store, err := server2.NewHeatmapStore(cfg)
	if err != nil {
		log.Fatalf("failed to initialize store: %v", err)
	}

	server2.StartHTTPServer(cfg, store)
}
