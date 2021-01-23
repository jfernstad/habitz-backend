package mock

import (
	"time"

	"github.com/jfernstad/habitz/web/internal"
)

type HabitzService struct{}

func (m *HabitzService) Users() ([]string, error)     { return nil, nil }
func (m *HabitzService) CreateUser(name string) error { return nil }
func (m *HabitzService) Templates(user, weekday string) ([]internal.HabitTemplate, error) {
	return nil, nil
}
func (m *HabitzService) CreateTemplate(user, weekday, habit string) error { return nil }
func (m *HabitzService) RemoveTemplate(user, weekday, habit string) error { return nil }
func (m *HabitzService) HabitEntries(user string, date time.Time) ([]*internal.HabitEntry, error) {
	return nil, nil
}
func (m *HabitzService) CreateHabitEntry(user, weekday, habit string, complete bool) (*internal.HabitEntry, error) {
	return nil, nil
}
func (m *HabitzService) UpdateHabitEntry(id int, complete bool) (*internal.HabitEntry, error) {
	return nil, nil
}
