package internal

import (
	"strings"
	"time"
)

func Today() string {
	return time.Now().UTC().Truncate(24 * time.Hour).Format("2006-01-02")
}

func Weekday() string {
	return strings.ToLower(time.Now().UTC().Truncate(24 * time.Hour).Weekday().String())
}
