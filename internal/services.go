package internal

import (
	"time"

	"github.com/jfernstad/habitz/web/internal/repository"
)

type Services struct {
	HabitzService HabitzServicer
}

type HabitzServicer interface {
	Users() ([]string, error) // Obsolete?
	UserWithExternalID(externalID string, provider string) (*repository.User, error)

	CreateExternalUser(external *repository.ExternalUser) (*repository.User, error)

	Templates(user string) ([]*repository.WeekHabitTemplates, error)
	WeekdayTemplates(user, weekday string) ([]*repository.WeekdayHabitTemplate, error)
	CreateTemplate(user, weekday, habit string) error
	RemoveTemplate(user, weekday, habit string) error
	RemoveEntry(user, habit string, date time.Time) error

	HabitEntries(user string, date string) ([]*repository.HabitEntry, error)
	CreateHabitEntry(user, weekday, habit string) (*repository.HabitEntry, error)
	UpdateHabitEntry(id int, complete bool) (*repository.HabitEntry, error)
}
