package parser

import "strconv"

func ParseFloat(value string) float64 {
	result, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0
	}
	return result
}
