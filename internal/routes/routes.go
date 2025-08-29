package routes

import (
	"berlin-heatmap/internal/config"
	"berlin-heatmap/internal/gql"
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

type Route = gql.Route

func FetchRoutes(ctx context.Context, cfg config.Config, routeNames []string, mode string) []Route {
	client := gql.New(cfg.APIURL)
	routes, err := client.Routes(ctx, routeNames, mode)
	if err != nil {
		log.Fatalf("fetch routes: %v", err)
	}
	log.Printf("fetched %d routes", len(routes))
	return routes
}

func WriteRoutes(routes []Route, outDir string, fileName string) {
	path := filepath.Join(outDir, fileName)
	log.Println(path)
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(routes); err != nil {
		log.Fatal(err)
	}
	log.Printf("wrote %s", path)
}
