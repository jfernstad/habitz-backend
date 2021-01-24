package internal

import (
	"time"
)

type HabitTemplate struct {
	Name    string `json:"name" db:"name"`
	Weekday string `json:"weekday" db:"weekday"`
	Habit   string `json:"habit" db:"habit"`
}

type HabitEntry struct {
	ID         int        `json:"id" db:"id"`
	Name       string     `json:"name,omitempty" db:"name"`
	Weekday    string     `json:"weekday" db:"weekday"`
	Habit      string     `json:"habit" db:"habit"`
	Complete   bool       `json:"complete" db:"complete"`
	Date       string     `json:"date,omitempty" db:"date"`
	CompleteAt *time.Time `json:"complete_at,omitempty" db:"complete_at"`
}
