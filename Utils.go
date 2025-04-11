package Binance

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

func GetIntervalFromString(interval string) (intervalValue int64, multiplier int64, intervalRune rune, exists bool, parseErr error) {
	// Fetch the last character
	intervalRune = rune(interval[len(interval)-1])

	// Parse the rest of the string (excluding the last character) as an integer
	restOfString := interval[:len(interval)-1]
	multiplier, parseErr = ParseInt(restOfString)
	if parseErr != nil {
		return 0, 0, 0, false, LocalError(PARSING_ERR, parseErr.Error())
	}

	INTERVALS_mu.Lock()
	defer INTERVALS_mu.Unlock()

	intervalValue, exists = STATIC_INTERVAL_CHARS[intervalRune]

	return intervalValue, multiplier, intervalRune, exists, nil
}

func GetOpenCloseTimes(currentTime int64, interval string) (openTime int64, closeTime int64, err *Error) {

	intervalValue, multiplier, intervalRune, exists, parseErr := GetIntervalFromString(interval)

	if parseErr != nil {
		LocalError(PARSING_ERR, fmt.Sprintf("Error parsing integer: %s", parseErr.Error()))
		return 0, 0, LocalError(PARSING_ERR, parseErr.Error())
	}

	if multiplier <= 0 {
		return 0, 0, LocalError(INVALID_VALUE_ERR, fmt.Sprintf("Multiplier value of '%d' must be greater than 0 in '%s' is invalid", multiplier, interval))
	}

	INTERVALS_mu.Lock()
	defer INTERVALS_mu.Unlock()

	if exists {
		openTime = currentTime - (currentTime % (intervalValue * multiplier))
		closeTime = openTime + (intervalValue * multiplier) - 1
		return openTime, closeTime, nil
	}

	baseUnix_time_obj := time.Unix(0, 0)
	current_time_obj := time.Unix(0, currentTime*int64(time.Millisecond))

	switch intervalRune {
	case COMPLEX_INTERVALS.WEEK:
		weekDay_offset := int(current_time_obj.Weekday() - time.Monday)
		monday_time := current_time_obj.AddDate(0, 0, weekDay_offset).UnixMilli()

		unixFirstWeek_Dayoffset := int(baseUnix_time_obj.Weekday() - time.Monday)
		unixFirstWeek_date := baseUnix_time_obj.AddDate(0, 0, unixFirstWeek_Dayoffset)
		unixFirstWeek_offset := unixFirstWeek_date.UnixMilli()

		timestamp_to_check := monday_time + unixFirstWeek_offset

		openTime = timestamp_to_check - (timestamp_to_check % (WEEK * multiplier))
		closeTime = openTime + (WEEK * multiplier) - 1

		return openTime - unixFirstWeek_offset, closeTime - unixFirstWeek_offset, nil
	case COMPLEX_INTERVALS.MONTH:
		yearNumber := current_time_obj.Year() - 1970
		currentMonthNumber := int(current_time_obj.Month() - 1)
		monthsSinceEpoch := yearNumber*12 + currentMonthNumber

		monthsToRemoveFromCurrentTime := monthsSinceEpoch % int(multiplier)

		openTime_date := time.Date(current_time_obj.Year(), current_time_obj.Month()-time.Month(monthsToRemoveFromCurrentTime), 1, 0, 0, 0, 0, time.UTC)
		closeTime_date := time.Date(current_time_obj.Year(), current_time_obj.Month()-time.Month(monthsToRemoveFromCurrentTime)+time.Month(multiplier), 1, 0, 0, 0, 0, time.UTC)

		openTime = openTime_date.UnixMilli()
		closeTime = closeTime_date.UnixMilli() - 1

		return openTime, closeTime, nil
	case COMPLEX_INTERVALS.YEAR:
		yearNumber := current_time_obj.Year() - 1970

		yearsToRemoveFromCurrentTime := yearNumber - (yearNumber % int(multiplier))

		openTime_date := time.Date(current_time_obj.Year()-yearsToRemoveFromCurrentTime, time.January, 1, 0, 0, 0, 0, time.UTC)
		closeTime_date := time.Date(current_time_obj.Year()-yearsToRemoveFromCurrentTime+int(multiplier), time.January, 1, 0, 0, 0, 0, time.UTC)

		openTime = openTime_date.UnixMilli()
		closeTime = closeTime_date.UnixMilli() - 1

		return openTime, closeTime, nil
	}

	return 0, 0, LocalError(INVALID_VALUE_ERR, fmt.Sprintf("Invalid interval rune of '%s' is invalid in '%s' is invalid", string(intervalRune), interval))
}

func ParseInt(intStr string) (int64, error) {
	return strconv.ParseInt(intStr, 10, 64)
}

func ParseFloat(floatStr string) (float64, error) {
	precision := GetStringNumberPrecision(floatStr)

	float, err := strconv.ParseFloat(floatStr, 64)
	if err != nil {
		return float, err
	}

	return ToFixed_Round(float, precision), nil
}

func GetStringNumberPrecision(numStr string) int {
	lastNumberIndex := 0
	dotIndex := 0

	dotFound := false

	for i, char := range numStr {
		if char == '.' {
			dotFound = true
			dotIndex = i
		} else if char != '0' {
			lastNumberIndex = i
		}
	}

	if !dotFound {
		dotIndex = len(numStr)
	}

	precision := lastNumberIndex - dotIndex

	if precision < 0 {
		precision++ // because if the number is right before the '.', then the precision must be 0, not -1 (so it's offset by 1)
	}

	return precision
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
	precision := GetStringNumberPrecision(tickSize)

	return Round_priceStr(priceStr, precision)
}

func Round_priceStr(priceStr string, precision int) string {

	for i, char := range priceStr {
		if char != '0' {
			priceStr = priceStr[i:]
			break
		}
	}

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
