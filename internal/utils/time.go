package utils

import (
	"fmt"
	"time"
)

func ConvertStringToTime(timeStamp string, layout string) (time.Time, error) {
	t, err := time.Parse(layout, timeStamp)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func ConvertTimeToStringUnix(timeStamp time.Time) string {
	return fmt.Sprintf("%d", timeStamp.Unix())
}

func TrimTimeToMinute(timeStamp string) string {
	return fmt.Sprintf("%02s", timeStamp)
}
