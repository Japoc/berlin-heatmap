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

	routeConfig := []struct {
		nameFilter  []string
		mode        string
		outFileName string
	}{
		{
			nameFilter:  []string{"S"},
			mode:        "RAIL",
			outFileName: "sbahn_routes",
		},
		{
			nameFilter:  []string{"U"},
			mode:        "RAIL",
			outFileName: "metro_routes",
		},
		{
			nameFilter:  []string{"S"},
			mode:        "RAIL",
			outFileName: "sbahn_routes",
		},
		{
			nameFilter:  []string{""},
			mode:        "BUS",
			outFileName: "bus_routes",
		},
		{
			nameFilter:  []string{""},
			mode:        "TRAM",
			outFileName: "tram_routes",
		},
	}

	for _, conf := range routeConfig {
		fetched := routes.FetchRoutes(ctx, cfg, conf.nameFilter, conf.mode)
		routes.WriteRoutes(fetched, cfg.OutDir, conf.outFileName+".json")
	}

	log.Printf("done. routes in %s", cfg.OutDir)
}
