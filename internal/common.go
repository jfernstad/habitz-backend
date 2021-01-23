package internal

import "time"

func Today() time.Time {
	return time.Now().UTC().Truncate(24 * time.Hour) // UTC midnight
}

func Weekday() string {
	return Today().Weekday().String()
}
