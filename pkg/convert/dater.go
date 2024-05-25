package convert

import (
	"strings"
	"time"
)

const (
	dateTimeFormat             = "2006-01-02 15:04:05"
	dateFormat                 = "2006-01-02"
	dateFormatWithoutDelimiter = "20060102"
	locationJakarta            = "Asia/Jakarta"
)

func StringToTime(timeString string, formatDateType string) (time.Time, error) {
	loc, _ := time.LoadLocation(locationJakarta)

	formatDate := dateTimeFormat
	if formatDateType == "2" {
		formatDate = dateFormat
	} else if formatDateType == "3" {
		formatDate = dateFormatWithoutDelimiter
	}

	t, err := time.ParseInLocation(formatDate, timeString, loc)
	if err != nil {
		return time.Now(), err
	}
	return t, nil
}

func TimeToString(inputTime time.Time) string {
	return inputTime.Format(dateTimeFormat)
}

// StringToDateWIB input "2006-01-02" to time
func StringToDateWIB(input string) (time.Time, error) {
	loc, err := time.LoadLocation(locationJakarta)
	if err != nil {
		return time.Now(), err
	}
	date, err := time.Parse(dateFormat, input)
	if err != nil {
		return time.Now(), err
	}
	return date.In(loc), nil
}

// DateToStringWIB input time.Time to "2006-01-02"
func DateToStringWIB(input time.Time) (string, error) {
	loc, err := time.LoadLocation(locationJakarta)
	if err != nil {
		return "", err
	}
	dateWIB := input.In(loc)
	dateString := dateWIB.Format(dateFormat)

	return dateString, nil
}

func SplitStringToTime(timeString string) (time.Time, error) {
	loc, _ := time.LoadLocation(locationJakarta)
	trxDate := strings.Split(timeString, "TRX")
	t, err := time.ParseInLocation(dateTimeFormat, trxDate[1], loc)
	if err != nil {
		return time.Now(), err
	}
	return t, nil
}

// TODO: time string conversion to time format, when an error occurs the input will be returned 0001-01-01 00:00:00
func ParseToDateTime(date string) time.Time {
	res, err := time.ParseInLocation(dateTimeFormat, date, time.Local)
	if err != nil {
		return time.Date(1, time.Month(1), 1, 0, 0, 0, 0, time.Local) // 0001-01-01 00:00:00 +0700
	}

	return res
}

func GetBeforeDateString(timeString string) (string, error) {
	date, err := time.ParseInLocation(dateFormat, timeString, time.Local)
	if err != nil {
		return "", err
	}
	return date.AddDate(0, 0, -1).Format(dateFormat), nil
}

func ParseStringToDuration(str string) (duration time.Duration, err error) {
	duration, err = time.ParseDuration(str)
	if err != nil {
		return
	}

	return
}
