package server

import (
	"image"
	"image/color"
)

// heatmapToPNG maps uint16 minutes → RGBA
func heatmapToPNG(heat []uint16, grid []GridCell, max int) *image.RGBA {
	//width := int(math.Sqrt(float64(len(grid))))
	//height := width
	width := 113
	height := 93
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := y*width + x
			//img.SetRGBA(x, y, colorForMinutes(heat[idx]))
			img.SetRGBA(x, y, colorRamp(int(heat[idx]), max))
		}
	}
	flipVerticalRGBA(img)
	return img
}

// heatmapToPNG maps uint16 minutes → RGBA
func heatmapScale() *image.RGBA {
	width := 113
	scaleSize := 10
	img := image.NewRGBA(image.Rect(0, 0, width, scaleSize))
	for y := 0; y < scaleSize; y++ {
		for x := 0; x < width; x++ {
			img.SetRGBA(x, y, colorForMinutes(uint16((float64(x)/float64(width))*60)))
			img.SetRGBA(x, y, colorRamp(int((float64(x)/float64(width))*60), 60))
		}
	}
	return img
}

func colorForMinutes(min uint16) color.RGBA {
	if min == 65535 {
		return color.RGBA{R: 0, G: 0, B: 0, A: 0}
	}
	if min > 60 {
		min = 60
	}
	t := float64(min) / 60.0
	r := uint8(255 * t)
	g := uint8(255 * (1 - 0.5*t))
	b := uint8(255 * (1 - t))
	var a uint8 = 180
	return color.RGBA{R: r, G: g, B: b, A: a}
}

func lerp(a, b uint8, t float64) uint8 {
	return uint8(float64(a) + t*(float64(b)-float64(a)))
}

func colorRamp(val, max int) color.RGBA {
	if val == 65535 {
		return color.RGBA{0, 0, 0, 0} // transparent
	}
	if val > max {
		val = max
	}
	t := float64(val) / float64(max)

	// Example: blue → green → yellow → red
	switch {
	case t < 0.33:
		return color.RGBA{
			R: 0,
			G: uint8(lerp(0, 255, t/0.33)),
			B: 255,
			A: 255,
		}
	case t < 0.66:
		return color.RGBA{
			R: uint8(lerp(0, 255, (t-0.33)/0.33)),
			G: 255,
			B: uint8(lerp(255, 0, (t-0.33)/0.33)),
			A: 255,
		}
	default:
		return color.RGBA{
			R: 255,
			G: uint8(lerp(255, 0, (t-0.66)/0.34)),
			B: 0,
			A: 255,
		}
	}
}

// flipVerticalRGBA flips an *image.RGBA vertically in place.
func flipVerticalRGBA(img *image.RGBA) {
	b := img.Bounds()
	stride := img.Stride
	tmp := make([]byte, stride)

	for y := 0; y < b.Dy()/2; y++ {
		top := img.Pix[y*stride : (y+1)*stride]
		bottom := img.Pix[(b.Dy()-1-y)*stride : (b.Dy()-y)*stride]

		// copy top row into tmp
		copy(tmp, top)
		// move bottom row into top
		copy(top, bottom)
		// move tmp (original top) into bottom
		copy(bottom, tmp)
	}
}
