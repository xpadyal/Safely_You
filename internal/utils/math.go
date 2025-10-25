package utils

// Round2 rounds a float64 to 2 decimal places
func Round2(f float64) float64 {
	return float64(int(f*100+0.5)) / 100
}
