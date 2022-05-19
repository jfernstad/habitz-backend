package repository

import (
	"time"
)

type WeekdayHabitTemplate struct {
	Name    string `json:"name" db:"name"`
	Weekday string `json:"weekday" db:"weekday"`
	Habit   string `json:"habit" db:"habit"`
}

type WeekHabitTemplates struct {
	Name     string   `json:"name" db:"name"`
	Weekdays []string `json:"weekdays" db:"weekdays"`
	Habit    string   `json:"habit" db:"habit"`
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

type User struct {
	ID              string `json:"id" db:"id"`
	Email           string `json:"email" db:"email"`
	Firstname       string `json:"name" db:"firstname"`
	Lastname        string `json:"lastname" db:"lastname"`
	ProfileImageURL string `json:"profile_image" db:"profile_image"`
}

type ExternalUser struct {
	User
	Provider   string `json:"auth_provider" db:"auth_provider"`
	ExternalID string `json:"external_id" db:"external_id"`
}
