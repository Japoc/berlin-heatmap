package server

import (
	"berlin-heatmap/internal/gql"
	"encoding/json"
	"image/png"
	"log"
	"math"
	"net/http"
	"strconv"
)

func StartHTTPServer(cfg *Config, store *HeatmapStore) {
	http.HandleFunc("/heatmap", store.HeatmapHandler)
	http.HandleFunc("/scale", store.heatmapScaleHandler)
	http.Handle("/route", withCORS(http.HandlerFunc(store.routeHandler)))
	log.Printf("listening on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, nil))
}

func (s *HeatmapStore) HeatmapHandler(w http.ResponseWriter, r *http.Request) {
	lat, lon := parseFloat(r.URL.Query().Get("lat")), parseFloat(r.URL.Query().Get("lon"))
	maxTime := parseInt(r.URL.Query().Get("max"))
	format := r.URL.Query().Get("format")
	if format == "" {
		format = "json"
	}

	heat := s.computeHeatmap(lat, lon)

	if format == "json" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(heat)
		return
	}

	if format == "png" {
		img := heatmapToPNG(heat, s.Grid, maxTime)
		w.Header().Set("Content-Type", "image/png")
		png.Encode(w, img)
		return
	}

	http.Error(w, "format must be json or png", 400)
}

func (s *HeatmapStore) heatmapScaleHandler(w http.ResponseWriter, r *http.Request) {
	img := heatmapScale()
	w.Header().Set("Content-Type", "image/png")
	png.Encode(w, img)
}

func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func parseInt(s string) int {
	v, _ := strconv.ParseInt(s, 10, 64)
	return int(v)
}

func (s *HeatmapStore) routeHandler(w http.ResponseWriter, r *http.Request) {
	latFrom, lonFrom := parseFloat(r.URL.Query().Get("latFrom")), parseFloat(r.URL.Query().Get("lonFrom"))
	latTo, lonTo := parseFloat(r.URL.Query().Get("latTo")), parseFloat(r.URL.Query().Get("lonTo"))

	originStops, originWalks := s.nearestKStops(latFrom, lonFrom)

	idx := MapLatLonToGridIndex(latTo, lonTo, gql.BBox{
		MinLon: 13.0884, MinLat: 52.3383,
		MaxLon: 13.7612, MaxLat: 52.6755,
	}, 400.0, 93, 113)

	distance := s.computeCell(originStops, originWalks, idx)

	type routeResponse struct {
		Distance uint16 `json:"distance"`
	}
	json.NewEncoder(w).Encode(routeResponse{Distance: distance})
}

func MapLatLonToGridIndex(lat, lon float64, bbox gql.BBox,
	cellSize float64, rows, cols int) int {

	latDegPerM := 1.0 / 111000.0
	lonDegPerM := 1.0 / (111000.0 * math.Cos((bbox.MinLat+bbox.MaxLat)/2*math.Pi/180))

	// compute row, col
	row := int((lat - bbox.MinLat) / (latDegPerM * cellSize))
	col := int((lon - bbox.MinLon) / (lonDegPerM * cellSize))

	// check bounds
	if row < 0 || row >= rows || col < 0 || col >= cols {
		return -1 // outside
	}

	// flatten index (row-major order)
	return row*cols + col
}

func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// handle preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(w, r)
	})
}
