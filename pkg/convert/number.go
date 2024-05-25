package convert

import (
	"math"
	"strconv"
)

// TODO : change to Generic
func StrToInt(text string, defaultReturn int) int {
	number := defaultReturn
	if text != "" {
		var err error
		number, err = strconv.Atoi(text)
		if err != nil {
			number = defaultReturn
		}
	}
	return number
}

func StrToInt64(text string, defaultReturn int64) int64 {
	number := defaultReturn
	if text != "" {
		var err error
		number, err = strconv.ParseInt(text, 10, 64)
		if err != nil {
			number = defaultReturn
		}
	}
	return number
}

func StrToUint64(text string, defaultReturn uint64) uint64 {
	number := defaultReturn
	if text != "" {
		var err error
		number, err = strconv.ParseUint(text, 10, 64)
		if err != nil {
			number = defaultReturn
		}
	}
	return number
}

// This function is used to round the fraction value after the comma to a certain number of digits
func RoundFloat64(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

// helper function for compare float
const float64EqualityThreshold = 0.1

// AlmostEqual return true if two float diferent < float64EqualityThreshold
func AlmostEqual(a, b float64) bool {
	dif := math.Abs(a - b)
	if dif < 0 {
		dif = -dif
	}
	return dif <= float64EqualityThreshold
}
