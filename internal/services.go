package internal

import "time"

type Services struct {
	HabitzService HabitzServicer
}

type HabitzServicer interface {
	Users() ([]string, error) // Obsolete?
	UserWithExternalID(externalID string, provider string) (*User, error)

	CreateUser(name string) error // Obsolete
	CreateExternalUser(external *ExternalUser) (*User, error)

	Templates(user string) ([]*WeekHabitTemplates, error)
	WeekdayTemplates(user, weekday string) ([]*WeekdayHabitTemplate, error)
	CreateTemplate(user, weekday, habit string) error
	RemoveTemplate(user, weekday, habit string) error
	RemoveEntry(user, habit string, date time.Time) error

	HabitEntries(user string, date string) ([]*HabitEntry, error)
	CreateHabitEntry(user, weekday, habit string) (*HabitEntry, error)
	UpdateHabitEntry(id int, complete bool) (*HabitEntry, error)
}
