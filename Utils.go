package Binance

import (
	"math"
	"strconv"
	"strings"
)

func ParseInt(intStr string) (int64, error) {
	return strconv.ParseInt(intStr, 10, 64)
}

func ParseFloat(floatStr string) (float64, error) {
	return strconv.ParseFloat(floatStr, 64)
}

func DetectDotNumIndexes(numStr string) (dotIndex int, numIndex int) {
	dotIndex = -1
	numIndex = -1
	for i, char := range numStr {
		switch char {
		case '.':
			dotIndex = i
		case '0':
		default:
			numIndex = i
		}
	}

	return dotIndex, numIndex
}

func Format_TickSize_str(priceStr string, tickSize string) string {
	_, numIndex := DetectDotNumIndexes(tickSize)
	if numIndex == -1 {
		return "0"
	}

	tickSize_dotIndex, tickSize_numIndex := DetectDotNumIndexes(tickSize)
	if tickSize_numIndex == -1 {
		return priceStr
	}

	var precision int

	if tickSize_dotIndex == -1 {
		precision = tickSize_numIndex - len(tickSize)
	} else {
		precision = tickSize_numIndex - tickSize_dotIndex
	}

	if precision < 0 {
		precision++
	}

	return Round_priceStr(priceStr, precision)
}

func Round_priceStr(priceStr string, precision int) string {
	if precision == 0 {
		return strings.Split(priceStr, ".")[0]
	}

	if precision < 0 {
		abs_precision := -precision
		priceStr = strings.Split(priceStr, ".")[0]
		length := len(priceStr)
		endIndex := length - abs_precision

		if abs_precision >= length {
			return "0"
		}

		return priceStr[:endIndex] + strings.Repeat("0", abs_precision)
	} else {
		dotIndex, _ := DetectDotNumIndexes(priceStr)
		if dotIndex == -1 {
			return priceStr + "." + strings.Repeat("0", precision)
		}

		arr := strings.Split(priceStr, ".")
		intStr, decimalStr := arr[0], arr[1]
		decimalLength := len(decimalStr)

		if decimalLength >= precision {
			decimalStr = decimalStr[:precision]
		} else {
			decimalStr += strings.Repeat("0", precision-decimalLength)
		}

		return intStr + "." + decimalStr
	}
}

func ToFixed_Floor(price float64, precision int) float64 {
	return math.Floor(price*math.Pow10(precision)) / math.Pow10(precision)
}

func ToFixed_Round(price float64, precision int) float64 {
	return math.Round(price*math.Pow10(precision)) / math.Pow10(precision)
}

func ToFixed_Ceil(price float64, precision int) float64 {
	return math.Ceil(price*math.Pow10(precision)) / math.Pow10(precision)
}
