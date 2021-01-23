package internal

import (
	"database/sql"
	"time"
)

type HabitTemplate struct {
	Name    string `json:"name" db:"name"`
	Weekday string `json:"weekday" db:"weekday"`
	Habit   string `json:"habit" db:"habit"`
}

type HabitEntry struct {
	ID         int          `json:"id" db:"id"`
	Name       string       `json:"name,omitempty" db:"name"`
	Weekday    string       `json:"weekday" db:"weekday"`
	Habit      string       `json:"habit" db:"habit"`
	Complete   bool         `json:"complete" db:"complete"`
	Date       SQLiteTime   `json:"date,omitempty" db:"date"`
	CompleteAt sql.NullTime `json:"complete_at,omitempty" db:"complete_at"`
}

type SQLiteTime time.Time

func (t *SQLiteTime) Scan(v interface{}) error {
	// Should be more strictly to check this type.
	vt, err := time.Parse("2006-01-02T15:04:05", string(v.([]byte)))
	if err != nil {
		return err
	}
	*t = SQLiteTime(vt)
	return nil
}
