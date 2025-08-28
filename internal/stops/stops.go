package stops

import (
	"berlin-heatmap/internal/config"
	"berlin-heatmap/internal/gql"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type Stop = gql.Stop

func FetchStops(ctx context.Context, cfg config.Config) []Stop {
	client := gql.New(cfg.APIURL)
	stops, err := client.Stops(ctx, gql.BBox(cfg.BBox))
	if err != nil {
		log.Fatalf("fetch stops: %v", err)
	}
	log.Printf("fetched %d stops", len(stops))
	return stops
}

func WriteStopsCSV(stops []Stop, outDir string) {
	filePath := filepath.Join(outDir, "stops.csv")
	csvFile, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("create stops.csv: %v", err)
	}
	defer csvFile.Close()

	w := csv.NewWriter(csvFile)
	defer w.Flush()

	_ = w.Write([]string{"index", "id", "lat", "lon"})
	for i, s := range stops {
		_ = w.Write([]string{
			strconv.Itoa(i),
			s.ID,
			fmt.Sprintf("%f", s.Lat),
			fmt.Sprintf("%f", s.Lon),
		})
	}
	log.Printf("wrote %s", filePath)
}
