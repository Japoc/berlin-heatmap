package walk

import (
	"math"
)

const WalkSpeedMps = 1.4 // ~5 km/h

type Link struct {
	StopIdx int
	WalkMin uint16 // ceil(distance / speed / 60)
}

func WalkMinutesMeters(d float64) uint16 {
	sec := d / WalkSpeedMps
	min := math.Ceil(sec / 60.0)
	if min > 65534 {
		return 65535
	}
	return uint16(min)
}
