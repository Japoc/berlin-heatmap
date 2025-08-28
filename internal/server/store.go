package server

import (
	"berlin-heatmap/internal/matrix"
	"encoding/csv"
	"encoding/json"
	"os"
	"strconv"
)

type Stop struct {
	ID  string  `json:"id"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type GridCell struct {
	Lat   float64  `json:"lat"`
	Lon   float64  `json:"lon"`
	Stops []int    `json:"stops"`
	Walks []uint16 `json:"walks"`
}

type HeatmapStore struct {
	Matrix *matrix.M
	Grid   []GridCell
	Stops  []Stop
	K      int
}

func NewHeatmapStore(cfg *Config) (*HeatmapStore, error) {
	// load matrix
	M, err := matrix.Open(cfg.MatrixPath)
	if err != nil {
		return nil, err
	}

	// load grid
	gridF, err := os.Open(cfg.GridFile)
	if err != nil {
		return nil, err
	}
	defer gridF.Close()
	var grid []GridCell
	if err := json.NewDecoder(gridF).Decode(&grid); err != nil {
		return nil, err
	}

	// load stops
	stopsFile, err := os.Open(cfg.StopsFile)
	if err != nil {
		return nil, err
	}
	defer stopsFile.Close()
	r := csv.NewReader(stopsFile)
	records, _ := r.ReadAll()
	stops := make([]Stop, 0, len(records)-1)
	for i, rec := range records {
		if i == 0 {
			continue
		}
		lat, _ := strconv.ParseFloat(rec[2], 64)
		lon, _ := strconv.ParseFloat(rec[3], 64)
		stops = append(stops, Stop{ID: rec[1], Lat: lat, Lon: lon})
	}

	return &HeatmapStore{
		Matrix: M,
		Grid:   grid,
		Stops:  stops,
		K:      cfg.K,
	}, nil
}
