package internal

import (
	"strings"
	"time"
)

func ShortDate(d time.Time) string {
	return d.UTC().Truncate(24 * time.Hour).Format("2006-01-02")
}

func Today() string {
	return ShortDate(time.Now())
}

func Weekday() string {
	return strings.ToLower(time.Now().UTC().Truncate(24 * time.Hour).Weekday().String())
}
