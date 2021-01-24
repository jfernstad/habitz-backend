package sqlite

import (
	"log"
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
	db *sqlx.DB
}

func NewHabitzService(db *sqlx.DB) internal.HabitzServicer {
	hs := &habitzService{
		db: db,
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

func (m *habitzService) Users() ([]string, error) {
	sql, _, _ := sq.Select("*").From("users").ToSql()
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

	if _, err := m.db.Exec(sql, args...); err != nil {
		return err
	}

	return nil
}

func (m *habitzService) Templates(user, weekday string) ([]*internal.HabitTemplate, error) {
	sql, args, _ := sq.Select("name", "weekday", "habit").
		From("habit_templates").
		Where(sq.Eq{"name": user, "weekday": weekday}).
		ToSql()

	rows, err := m.db.Queryx(sql, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	userTemplates := []*internal.HabitTemplate{}
	for rows.Next() {
		var tmpl internal.HabitTemplate

		if err = rows.StructScan(&tmpl); err != nil {
			return nil, err
		}
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

	if _, err := m.db.Exec(sql, args...); err != nil {
		return err
	}

	return nil
}

func (m *habitzService) RemoveTemplate(user, weekday, habit string) error { return nil }
func (m *habitzService) HabitEntries(user string, date string) ([]*internal.HabitEntry, error) {
	sql, args, _ := sq.Select("*").
		From("habit_entries").
		Where(sq.Eq{"name": user, "date": date}).
		ToSql()

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

	return &entry, nil
}

func (m *habitzService) UpdateHabitEntry(id int, complete bool) (*internal.HabitEntry, error) {

	sql, args, _ := sq.Update("habit_entries").
		Set("complete", complete).
		Set("complete_at", time.Now().UTC().Format(sqlTimeFormat)).
		Where(sq.Eq{"id": id}).ToSql()

	log.Println(sql, args)

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

	return &entry, nil
}
