package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type BBox struct {
	MinLon, MinLat, MaxLon, MaxLat float64
}

type Config struct {
	APIURL   string
	Bucket   string
	BBox     BBox
	Grid     int
	K        int
	OutDir   string
	DateTime time.Time
}

var (
	apiURL   = flag.String("graphql", "", "GTFS-GraphQL endpoint")
	bucket   = flag.String("bucket", "weekday_08", "time bucket label")
	bboxStr  = flag.String("bbox", "13.0884,52.3383,13.7612,52.6755", "minLon,minLat,maxLon,maxLat")
	grid     = flag.Int("grid", 400, "grid spacing in meters")
	k        = flag.Int("k", 5, "nearest stops per cell")
	outDir   = flag.String("out", "./data/artifacts", "output dir")
	dateTime = flag.String("datetime", "2025-08-22T12:00:00Z", "datetime used for calculation")
)

func ParseFlags() Config {
	flag.Parse()
	if *apiURL == "" {
		log.Fatal("must pass -graphql endpoint")
	}

	parsedDateTime, err := time.Parse(time.RFC3339, *dateTime)
	if err != nil {
		log.Fatal("invalid format of datetime")
	}

	bbox, err := parseBBox(*bboxStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.MkdirAll(*outDir, 0755); err != nil {
		log.Fatal(err)
	}

	return Config{
		APIURL:   *apiURL,
		Bucket:   *bucket,
		BBox:     bbox,
		Grid:     *grid,
		K:        *k,
		OutDir:   *outDir,
		DateTime: parsedDateTime,
	}
}

func parseBBox(s string) (BBox, error) {
	parts := strings.Split(s, ",")
	if len(parts) != 4 {
		return BBox{}, fmt.Errorf("need 4 numbers, got %d", len(parts))
	}
	vals := make([]float64, 4)
	for i, p := range parts {
		v, err := strconv.ParseFloat(p, 64)
		if err != nil {
			return BBox{}, err
		}
		vals[i] = v
	}
	return BBox{vals[0], vals[1], vals[2], vals[3]}, nil
}
