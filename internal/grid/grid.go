package grid

import (
	"berlin-heatmap/internal/config"
	"berlin-heatmap/internal/gql"
	"encoding/json"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"

	"berlin-heatmap/internal/walk"
)

type Cell struct {
	Lat, Lon float64
	Stops    []int
	Walks    []uint16
}

func BuildGrid(stops []gql.Stop, cfg config.Config) []Cell {
	log.Printf("building grid links")

	bbox := cfg.BBox
	cellSize := float64(cfg.Grid)

	latDegPerM := 1.0 / 111000.0
	lonDegPerM := 1.0 / (111000.0 * math.Cos((bbox.MinLat+bbox.MaxLat)/2*math.Pi/180))
	cols := int((bbox.MaxLon - bbox.MinLon) / (lonDegPerM * cellSize))
	rows := int((bbox.MaxLat - bbox.MinLat) / (latDegPerM * cellSize))
	log.Printf("grid %dx%d", cols, rows)

	var gridCells []Cell
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			clat := bbox.MinLat + float64(r)*latDegPerM*cellSize
			clon := bbox.MinLon + float64(c)*lonDegPerM*cellSize

			stopIdxs, walks := NearestKStops(clat, clon, stops, cfg.K)
			gridCells = append(gridCells, Cell{clat, clon, stopIdxs, walks})
		}
	}
	return gridCells
}

func WriteGrid(cells []Cell, cfg config.Config) {
	path := filepath.Join(cfg.OutDir, "grid_links.json")
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(cells); err != nil {
		log.Fatal(err)
	}
	log.Printf("wrote %s", path)
}

// naive implementation — replace with KD-tree for speed later
func NearestKStops(lat, lon float64, stops []gql.Stop, k int) ([]int, []uint16) {
	type candidate struct {
		idx int
		d   float64
	}
	all := make([]candidate, 0, len(stops))
	for i, s := range stops {
		d := Haversine(lat, lon, s.Lat, s.Lon)
		all = append(all, candidate{i, d})
	}
	sort.Slice(all, func(i, j int) bool { return all[i].d < all[j].d })

	stopIdxs := make([]int, 0, k)
	walks := make([]uint16, 0, k)
	for i := 0; i < k && i < len(all); i++ {
		stopIdxs = append(stopIdxs, all[i].idx)
		walks = append(walks, walk.WalkMinutesMeters(all[i].d))
	}
	return stopIdxs, walks
}

func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000
	φ1, λ1 := lat1*math.Pi/180, lon1*math.Pi/180
	φ2, λ2 := lat2*math.Pi/180, lon2*math.Pi/180
	dφ := φ2 - φ1
	dλ := λ2 - λ1
	a := math.Sin(dφ/2)*math.Sin(dφ/2) +
		math.Cos(φ1)*math.Cos(φ2)*math.Sin(dλ/2)*math.Sin(dλ/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
