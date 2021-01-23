package internal

import "time"

type Services struct {
	HabitzService HabitzServicer
}

type HabitzServicer interface {
	Users() ([]string, error)
	CreateUser(name string) error

	Templates(user, weekday string) ([]*HabitTemplate, error)
	CreateTemplate(user, weekday, habit string) error
	RemoveTemplate(user, weekday, habit string) error

	HabitEntries(user string, date time.Time) ([]*HabitEntry, error)
	CreateHabitEntry(user, weekday, habit string) (*HabitEntry, error)
	UpdateHabitEntry(id int, complete bool) (*HabitEntry, error)
}
