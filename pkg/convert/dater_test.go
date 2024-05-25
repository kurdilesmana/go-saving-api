package convert

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStringToTime(t *testing.T) {
	data := "2021-06-30 15:00:01"
	wrongData := "2021-30-06 15:00:01" // split date and month

	loc, _ := time.LoadLocation(locationJakarta)
	want := time.Date(2021, time.Month(6), 30, 15, 0, 1, 0, loc)

	t.Run("expect success", func(t *testing.T) {
		got, err := StringToTime(data, "1")

		assert.Nil(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("expect failure", func(t *testing.T) {
		got, err := StringToTime(wrongData, "1")

		assert.NotNil(t, err)
		assert.Equal(t, got.Year(), time.Now().Year())
		assert.Equal(t, got.Month(), time.Now().Month())
		assert.Equal(t, got.Weekday(), time.Now().Weekday())
		assert.Equal(t, got.Hour(), time.Now().Hour())
	})

}

func TestTimeToString(t *testing.T) {
	loc, _ := time.LoadLocation(locationJakarta)
	data := time.Date(2021, time.Month(6), 30, 15, 0, 1, 0, loc)
	want := "2021-06-30 15:00:01"

	got := TimeToString(data)

	assert.Equal(t, got, want)
}

func TestStringToDateWIB(t *testing.T) {
	data := "2021-06-30" // utc
	wrongData := "2021-06-30 15:00:01"

	loc, _ := time.LoadLocation(locationJakarta)
	want := time.Date(2021, time.Month(6), 30, 7, 0, 0, 0, loc)

	t.Run("expect success", func(t *testing.T) {
		got, err := StringToDateWIB(data)

		assert.Nil(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("expect failure", func(t *testing.T) {
		got, err := StringToDateWIB(wrongData)

		assert.NotNil(t, err)
		assert.Equal(t, got.Year(), time.Now().Year())
		assert.Equal(t, got.Month(), time.Now().Month())
		assert.Equal(t, got.Weekday(), time.Now().Weekday())
	})
}

func TestDateToStringWIB(t *testing.T) {
	loc, _ := time.LoadLocation(locationJakarta)
	data := time.Date(2021, time.Month(6), 30, 7, 0, 0, 0, loc)
	want := "2021-06-30"

	got, err := DateToStringWIB(data)

	assert.Nil(t, err)
	assert.Equal(t, want, got)
}

func TestParseToDateTime(t *testing.T) {
	tests := []struct {
		Date string
		Want time.Time
	}{
		{
			Date: "YYYYY-MM-DD HH:mm:ss",
			Want: time.Date(1, time.Month(1), 1, 0, 0, 0, 0, time.Local),
		},
		{
			Date: "2022-07-15 21:45:19",
			Want: time.Date(2022, time.Month(7), 15, 21, 45, 19, 0, time.Local),
		},
		{
			Date: "2019-12-31 23:56:59",
			Want: time.Date(2019, time.Month(12), 31, 23, 56, 59, 0, time.Local),
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.Want, ParseToDateTime(test.Date))
	}
}

func TestGetBeforeDateString(t *testing.T) {
	tests := []struct {
		Date string
		Want string
	}{
		{
			Date: "YYYYY-MM-DD",
			Want: "",
		},
		{
			Date: "2021-06-30",
			Want: "2021-06-29",
		},
		{
			Date: "2021-02-01",
			Want: "2021-01-31",
		},
	}
	for _, test := range tests {
		got, _ := GetBeforeDateString(test.Date)
		assert.Equal(t, test.Want, got)
	}
}
