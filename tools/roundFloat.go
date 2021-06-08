package tools

import (
	"fmt"
	"strconv"
)

func RoundFloat(f float64) float64 {
	r, _ := strconv.ParseFloat(fmt.Sprintf("%.10f", f), 64)
	return r
}

func RoundFloatPointer(f *float64) {
	*f, _ = strconv.ParseFloat(fmt.Sprintf("%.10f", *f), 64)
}
