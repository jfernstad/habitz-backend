package sqlite

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	sq "github.com/Masterminds/squirrel"
	"github.com/jfernstad/habitz/web/internal"
)

const createUserTable = `
CREATE TABLE IF NOT EXISTS users(
	name TEXT PRIMARY KEY
);
`

const createHabitTemplateTable = `
CREATE TABLE IF NOT EXISTS habit_templates (
	name TEXT,
	weekday TEXT,
	habit TEXT,
	PRIMARY KEY (name, weekday, habit)
) WITHOUT ROWID;
`

const createEntryTable = `
CREATE TABLE IF NOT EXISTS habit_entries(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	weekday TEXT,
	date TEXT,
	habit TEXT,
	complete INTEGER,
	complete_at TIMESTAMP
);
`

const sqlTimeFormat = "2006-01-02 15:04:05"

type habitzService struct {
	db    *sqlx.DB
	debug bool
}

func NewHabitzService(db *sqlx.DB, debug bool) internal.HabitzServicer {
	hs := &habitzService{
		db:    db,
		debug: debug,
	}

	if err := hs.initSQLDatabase(); err != nil {
		log.Fatal(err)
	}

	return hs
}

func (m *habitzService) initSQLDatabase() error {
	_, err := m.db.Exec(createUserTable)
	if err != nil {
		return err
	}

	_, err = m.db.Exec(createHabitTemplateTable)
	if err != nil {
		return err
	}

	_, err = m.db.Exec(createEntryTable)
	if err != nil {
		return err
	}

	return nil
}

func (m *habitzService) log(msg string) {
	if m.debug {
		log.Println("sql: " + msg)
	}
}

func (m *habitzService) Users() ([]string, error) {
	sql, _, _ := sq.Select("*").From("users").ToSql()

	m.log("Users: " + sql)

	rows, err := m.db.Query(sql)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []string{}

	for rows.Next() {
		var name string

		if err = rows.Scan(&name); err != nil {
			return nil, err
		}
		m.log(" - " + name)
		users = append(users, name)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (m *habitzService) CreateUser(name string) error {

	sql, args, _ := sq.Insert("users").
		Columns("name").Values(name).
		ToSql()

	m.log("CreateUser: " + sql + " >>  " + name)

	if _, err := m.db.Exec(sql, args...); err != nil {
		return err
	}

	return nil
}

func (m *habitzService) Templates(user string) ([]*internal.WeekHabitTemplates, error) {
	sql, args, _ := sq.Select("name", "weekday", "habit").
		From("habit_templates").
		Where(sq.Eq{"name": user}).
		Suffix("COLLATE NOCASE").
		ToSql()

	m.log("Templates: " + sql + " >> " + user)

	rows, err := m.db.Queryx(sql, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userTemplates := []*internal.WeekHabitTemplates{}
	for rows.Next() {
		var tmpl internal.WeekdayHabitTemplate
		var weekTmpl internal.WeekHabitTemplates

		if err = rows.StructScan(&tmpl); err != nil {
			return nil, err
		}
		m.log(fmt.Sprintf(" - %+v", tmpl))

		exist := false
		for _, ut := range userTemplates {
			if ut.Habit == tmpl.Habit { // Already in array, append weekday
				ut.Weekdays = append(ut.Weekdays, tmpl.Weekday)
				exist = true
				break
			}
		}

		if !exist {
			weekTmpl.Name = tmpl.Name
			weekTmpl.Habit = tmpl.Habit
			weekTmpl.Weekdays = []string{tmpl.Weekday}
			userTemplates = append(userTemplates, &weekTmpl)
			continue
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userTemplates, nil

}

func (m *habitzService) WeekdayTemplates(user, weekday string) ([]*internal.WeekdayHabitTemplate, error) {
	sql, args, _ := sq.Select("name", "weekday", "habit").
		From("habit_templates").
		Where(sq.Eq{"name": user, "weekday": weekday}).
		ToSql()

	m.log("WeekdayTemplates: " + sql + " >> " + user + ", " + weekday)

	rows, err := m.db.Queryx(sql, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userTemplates := []*internal.WeekdayHabitTemplate{}
	for rows.Next() {
		var tmpl internal.WeekdayHabitTemplate

		if err = rows.StructScan(&tmpl); err != nil {
			return nil, err
		}
		m.log(fmt.Sprintf(" - %+v", tmpl))

		userTemplates = append(userTemplates, &tmpl)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userTemplates, nil
}

func (m *habitzService) CreateTemplate(user, weekday, habit string) error {
	sql, args, _ := sq.Insert("habit_templates").
		Columns("name", "weekday", "habit").Values(user, weekday, habit).
		ToSql()

	m.log("CreateTemplate: " + sql + " >> " + user + ", " + weekday + ", " + habit)

	if _, err := m.db.Exec(sql, args...); err != nil {
		return err
	}

	return nil
}

func (m *habitzService) RemoveTemplate(user, weekday, habit string) error {
	sql, args, _ := sq.Delete("habit_templates").
		Where(sq.Eq{"name": user, "weekday": weekday, "habit": habit}).
		ToSql()

	m.log("RemoveTemplate: " + sql + " >> " + user + ", " + weekday + ", " + habit)

	if _, err := m.db.Exec(sql, args...); err != nil {
		return err
	}

	return nil
}

func (m *habitzService) RemoveEntry(user, habit string, date time.Time) error {
	shortDate := internal.ShortDate(date)
	sql, args, _ := sq.Delete("habit_entries").
		Where(sq.Eq{"name": user, "date": shortDate, "habit": habit}).
		ToSql()

	m.log("RemoveEntry: " + sql + " >> " + user + ", " + shortDate + ", " + habit)

	if _, err := m.db.Exec(sql, args...); err != nil {
		return err
	}

	return nil
}

func (m *habitzService) HabitEntries(user string, date string) ([]*internal.HabitEntry, error) {
	sql, args, _ := sq.Select("*").
		From("habit_entries").
		Where(sq.Eq{"name": user, "date": date}).
		ToSql()

	m.log("HabitEntries: " + sql + " >> " + user + ", " + date)

	rows, err := m.db.Queryx(sql, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	habitEntries := []*internal.HabitEntry{}

	for rows.Next() {
		var entry internal.HabitEntry

		if err = rows.StructScan(&entry); err != nil {
			return nil, err
		}

		m.log(fmt.Sprintf(" - %+v", entry))

		habitEntries = append(habitEntries, &entry)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return habitEntries, nil
}

func (m *habitzService) CreateHabitEntry(user, weekday, habit string) (*internal.HabitEntry, error) {

	today := internal.Today()

	sql, args, _ := sq.Insert("habit_entries").
		Columns("name", "weekday", "habit", "date", "complete").
		Values(user, weekday, habit, today, 0).
		ToSql()

	m.log("CreateHabitEntry:" + " >> " + user + ", " + weekday + ", " + habit)

	if _, err := m.db.Exec(sql, args...); err != nil {
		return nil, err
	}

	// Retrieve last insert values
	entry := internal.HabitEntry{}

	sql, _, _ = sq.Select("*").
		From("habit_entries").
		OrderBy("id desc").
		Limit(1).ToSql()

	if err := m.db.QueryRowx(sql).StructScan(&entry); err != nil {
		return nil, err
	}

	m.log(fmt.Sprintf(" - %+v", entry))

	return &entry, nil
}

func (m *habitzService) UpdateHabitEntry(id int, complete bool) (*internal.HabitEntry, error) {

	query := sq.Update("habit_entries").
		Set("complete", complete)

	// Also update timestamp
	if complete == true {
		query = query.Set("complete_at", time.Now().UTC().Format(sqlTimeFormat))
	}

	sql, args, _ := query.
		Where(sq.Eq{"id": id}).ToSql()

	m.log("UpdateHabitEntry: " + sql + " >> " + strconv.FormatInt(int64(id), 10) + ", " + strconv.FormatBool(complete))

	if _, err := m.db.Exec(sql, args...); err != nil {
		return nil, err
	}

	// Retrieve full object
	sql, args, _ = sq.Select("*").
		From("habit_entries").
		Where(sq.Eq{"id": id}).
		ToSql()

	entry := internal.HabitEntry{}

	if err := m.db.QueryRowx(sql, args...).StructScan(&entry); err != nil {
		return nil, err
	}

	m.log(fmt.Sprintf(" - %+v", entry))

	return &entry, nil
}
