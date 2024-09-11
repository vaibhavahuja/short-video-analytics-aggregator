package utils

import (
	"fmt"
	"strconv"
	"time"
)

func ConvertStringToTime(timeStamp string, layout string) (time.Time, error) {
	t, err := time.Parse(layout, timeStamp)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
func ConvertUnixTimestampToTime(unixTimestamp string) (time.Time, error) {
	// Parse the Unix timestamp string to int64
	seconds, err := strconv.ParseInt(unixTimestamp, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	// Convert the seconds to time.Time
	t := time.Unix(seconds, 0)
	return t, nil
}

func ConvertTimeToStringUnix(timeStamp time.Time) string {
	return fmt.Sprintf("%d", timeStamp.Unix())
}

func TrimTimeToMinute(timeStamp string) string {
	return fmt.Sprintf("%02s", timeStamp)
}
