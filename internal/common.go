package internal

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func ShortDate(d time.Time) string {
	return d.UTC().Truncate(24 * time.Hour).Format("2006-01-02")
}

func Today() string {
	return ShortDate(time.Now())
}

func Weekday() string {
	return strings.ToLower(time.Now().UTC().Truncate(24 * time.Hour).Weekday().String())
}

// From https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func NewRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
