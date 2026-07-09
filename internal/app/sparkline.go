package app

import "math"

var brailleBase = '\u2800'

func RenderBrailleSparkline(values []float64, width, height int) string {
	if width < 1 || height < 1 || len(values) < 1 {
		return ""
	}

	cols := width * 2
	rows := height * 4

	maxVal := 0.0
	for _, v := range values {
		v = math.Abs(v)
		if v > maxVal {
			maxVal = v
		}
	}
	if maxVal == 0 {
		maxVal = 1
	}

	grid := make([][]bool, cols)
	for i := range grid {
		grid[i] = make([]bool, rows)
	}

	for i := 0; i < cols && i < len(values); i++ {
		idx := i
		if idx >= len(values) {
			idx = len(values) - 1
		}
		ratio := math.Abs(values[idx]) / maxVal
		if ratio > 1 {
			ratio = 1
		}
		filledRows := int(math.Round(ratio * float64(rows)))
		if filledRows > rows {
			filledRows = rows
		}
		for r := 0; r < filledRows; r++ {
			grid[cols-1-i][rows-1-r] = true
		}
	}

	var out string
	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			code := 0
			for dy := 0; dy < 4; dy++ {
				for dx := 0; dx < 2; dx++ {
					gx := col*2 + dx
					gy := row*4 + dy
					if gx < cols && gy < rows && grid[gx][gy] {
						bitPos := uint(dy*2 + dx)
						code |= 1 << bitPos
					}
				}
			}
			out += string(rune(int(brailleBase) + code))
		}
		if row < height-1 {
			out += "\n"
		}
	}

	return out
}
