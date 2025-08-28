package matrix

import (
	"berlin-heatmap/internal/config"
	"berlin-heatmap/internal/gql"
	"context"
	"encoding/binary"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const Unreachable = uint16(65535)

func BuildTravelMatrix(ctx context.Context, stops []gql.Stop, cfg config.Config) [][]uint16 {
	n := len(stops)
	log.Printf("building %dx%d matrix for bucket %s", n, n, cfg.Bucket)

	client := gql.New(cfg.APIURL)

	// init
	data := make([][]uint16, n)
	for i := range data {
		data[i] = make([]uint16, n)
		for j := range data[i] {
			if i == j {
				data[i][j] = 0
			} else {
				data[i][j] = Unreachable
			}
		}
	}

	type job struct{ i, j int }
	jobs := make(chan job, 1000)
	var wg sync.WaitGroup

	total := n * (n - 1)
	done := 0
	workerCount := 200
	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range jobs {
				minutes, err := client.TravelTime(ctx, stops[task.i], stops[task.j], cfg.DateTime)
				if err == nil && minutes >= 0 && minutes < int(Unreachable) {
					data[task.i][task.j] = uint16(minutes)
				}
				// Progress tracker
				done++
				if done%100 == 0 {
					log.Printf("progress: %d/%d (%.1f%%)", done, total, 100*float64(done)/float64(total))
				}
			}
		}()
	}

	go func() {
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if i != j {
					jobs <- job{i, j}
				}
			}
		}
		close(jobs)
	}()

	wg.Wait()
	return data
}

func WriteMatrix(data [][]uint16, cfg config.Config) {
	path := filepath.Join(cfg.OutDir, "matrix_"+cfg.Bucket+".bin")
	err := WriteUint16Matrix(path, len(data), func(i int) []uint16 {
		return data[i]
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("wrote %s", path)
}

func WriteUint16Matrix(path string, n int, rowProvider func(i int) []uint16) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// header: n (uint32)
	if err := binary.Write(f, binary.LittleEndian, uint32(n)); err != nil {
		return err
	}

	buf := make([]byte, 2*n)
	for i := 0; i < n; i++ {
		row := rowProvider(i) // len == n
		for j := 0; j < n; j++ {
			binary.LittleEndian.PutUint16(buf[2*j:], row[j])
		}
		if _, err := f.Write(buf); err != nil {
			return err
		}
	}
	return nil
}
