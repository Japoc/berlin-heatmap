package server

import (
	"berlin-heatmap/internal/walk"
	"math"
	"runtime"
	"sync"
)

func (s *HeatmapStore) computeHeatmap(lat, lon float64) []uint16 {
	originStops, originWalks := s.nearestKStops(lat, lon)

	heat := make([]uint16, len(s.Grid))
	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU()
	ch := make(chan int, numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range ch {
				heat[idx] = s.computeCell(originStops, originWalks, idx)
			}
		}()
	}
	for i := range s.Grid {
		ch <- i
	}
	close(ch)
	wg.Wait()
	return heat
}

func (s *HeatmapStore) nearestKStops(lat, lon float64) ([]int, []uint16) {
	type pair struct {
		idx int
		d   float64
	}
	k := s.K
	best := make([]pair, k)
	for i := range best {
		best[i].d = 1e12
	}
	for i, stop := range s.Stops {
		d := haversine(lat, lon, stop.Lat, stop.Lon)
		for j := 0; j < k; j++ {
			if d < best[j].d {
				copy(best[j+1:], best[j:k-1])
				best[j] = pair{i, d}
				break
			}
		}
	}
	ids := make([]int, k)
	walks := make([]uint16, k)
	for i := 0; i < k; i++ {
		ids[i] = best[i].idx
		walks[i] = walk.WalkMinutesMeters(best[i].d)
	}
	return ids, walks
}

func (s *HeatmapStore) computeCell(originStops []int, originWalks []uint16, idx int) uint16 {
	cell := s.Grid[idx]
	best := uint16(65535)
	for oi := range originStops {
		o := originStops[oi]
		ow := originWalks[oi]
		if ow == 65535 {
			continue
		}
		for di, d := range cell.Stops {
			dw := cell.Walks[di]
			if dw == 65535 {
				continue
			}
			matrixValue := s.Matrix.Get(o, d)
			if matrixValue == 65535 {
				continue
			}
			tt := ow + matrixValue + dw
			if tt < best {
				best = tt
			}
		}
	}
	return best
}

// haversine in meters
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000
	φ1, λ1 := lat1*math.Pi/180, lon1*math.Pi/180
	φ2, λ2 := lat2*math.Pi/180, lon2*math.Pi/180
	dφ := φ2 - φ1
	dλ := λ2 - λ1
	a := math.Sin(dφ/2)*math.Sin(dφ/2) + math.Cos(φ1)*math.Cos(φ2)*math.Sin(dλ/2)*math.Sin(dλ/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
