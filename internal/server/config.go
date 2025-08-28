package server

import "flag"

type Config struct {
	Port       string
	MatrixPath string
	GridFile   string
	StopsFile  string
	K          int
}

func ParseFlags() *Config {
	port := flag.String("port", "8088", "server port")
	matrixPath := flag.String("matrix", "../data/artifacts/matrix_weekday_08.bin", "stop->stop matrix")
	gridFile := flag.String("gridFile", "../data/artifacts/grid_links.json", "grid links JSON")
	stopsFile := flag.String("stops", "../data/artifacts/stops.csv", "stops CSV")

	flag.Parse()

	return &Config{
		Port:       *port,
		MatrixPath: *matrixPath,
		GridFile:   *gridFile,
		StopsFile:  *stopsFile,
		K:          35,
	}
}
