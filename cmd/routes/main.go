package main

import (
	"berlin-heatmap/internal/config"
	"berlin-heatmap/internal/routes"
	"context"
	"log"
)

func main() {
	cfg := config.ParseFlags()
	cfg.OutDir = "data/routes"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fetchedSBahn := routes.FetchRoutes(ctx, cfg, []string{"S"})
	fetchedMetro := routes.FetchRoutes(ctx, cfg, []string{"U"})
	routes.WriteRoutes(fetchedSBahn, cfg.OutDir, "sbahn_routes.json")
	routes.WriteRoutes(fetchedMetro, cfg.OutDir, "metro_routes.json")

	log.Printf("done. routes in %s", cfg.OutDir)
}
