package mock

import (
	"github.com/jfernstad/habitz/web/internal/repository"
)

type HabitzService struct{}

func (m *HabitzService) Users() ([]string, error)     { return nil, nil }
func (m *HabitzService) CreateUser(name string) error { return nil }
func (m *HabitzService) Templates(user, weekday string) ([]*repository.WeekHabitTemplates, error) {
	return nil, nil
}
func (m *HabitzService) CreateTemplate(user, weekday, habit string) error { return nil }
func (m *HabitzService) RemoveTemplate(user, weekday, habit string) error { return nil }
func (m *HabitzService) HabitEntries(user string, date string) ([]*repository.HabitEntry, error) {
	return nil, nil
}
func (m *HabitzService) CreateHabitEntry(user, weekday, habit string) (*repository.HabitEntry, error) {
	return nil, nil
}
func (m *HabitzService) UpdateHabitEntry(id int, complete bool) (*repository.HabitEntry, error) {
	return nil, nil
}
