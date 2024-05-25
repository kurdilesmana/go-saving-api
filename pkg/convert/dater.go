package convert

import (
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
