package endpoints

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/jfernstad/habitz/web/internal"
	"github.com/jfernstad/habitz/web/internal/repository"
)

type www struct {
	DefaultEndpoint
	service internal.HabitzServicer
	// HTML Templates
	indexTemplate    *template.Template
	loginTemplate    *template.Template
	todayTemplate    *template.Template
	newTemplate      *template.Template
	scheduleTemplate *template.Template

	googleClientID string
}

func NewWWWEndpoint(hs internal.HabitzServicer, googleClientID string) EndpointRouter {
	// Load HTML templates
	indexTmpl, err := template.ParseFiles("./cmd/backend/templates/index.tmpl")
	if err != nil {
		panic(err)
	}

	loginTmpl, err := template.ParseFiles("./cmd/backend/templates/login.tmpl")
	if err != nil {
		panic(err)
	}

	todayTmpl, err := template.ParseFiles("./cmd/backend/templates/today.tmpl")
	if err != nil {
		panic(err)
	}

	newTmpl, err := template.ParseFiles("./cmd/backend/templates/new.tmpl")
	if err != nil {
		panic(err)
	}

	scheduleTmpl, err := template.ParseFiles("./cmd/backend/templates/schedule.tmpl")
	if err != nil {
		panic(err)
	}

	return &www{
		service:          hs,
		indexTemplate:    indexTmpl,
		loginTemplate:    loginTmpl,
		todayTemplate:    todayTmpl,
		newTemplate:      newTmpl,
		scheduleTemplate: scheduleTmpl,
		googleClientID:   googleClientID,
	}
}

func (ww *www) Routes() chi.Router {
	router := NewRouter()

	router.Route("/", func(r chi.Router) {
		r.Get("/", ErrorHandler(ww.index))
		r.Get("/login", ErrorHandler(ww.login))
		r.Get("/today", ErrorHandler(ww.todaysHabitz))
		r.Get("/new", ErrorHandler(ww.newHabit))
		r.Get("/update/{user}", ErrorHandler(ww.updateHabit))
	})

	return router
}
func (ww *www) index(w http.ResponseWriter, r *http.Request) error {
	writeHTML(w, http.StatusOK, ww.indexTemplate, nil)
	return nil
}

func (ww *www) login(w http.ResponseWriter, r *http.Request) error {
	type loginRender struct {
		GoogleID string
	}
	state := loginRender{
		GoogleID: ww.googleClientID,
	}
	writeHTML(w, http.StatusOK, ww.loginTemplate, state)
	return nil
}

func (ww *www) todaysHabitz(w http.ResponseWriter, r *http.Request) error {

	// Load habits for all users
	allUsers, err := ww.service.Users()
	if err != nil {
		return newInternalServerErr("could not load users").Wrap(err)
	}

	// What day is it?
	today := internal.Today()
	weekday := internal.Weekday()
	allHabitz := []*habitState{}

	// Try to retrive todays habitz for all users
	for _, user := range allUsers {
		habitz, err := ww.service.HabitEntries(user, today)
		if err != nil {
			return newInternalServerErr("could not load habitz for today").Wrap(err)
		}
		// Todays entries might not have been created yet, lets create them
		if len(habitz) == 0 {
			log.Println("No entries for today, lets create them")

			habitz = []*repository.HabitEntry{}
			templates, err := ww.service.WeekdayTemplates(user, weekday)
			if err != nil {
				return newInternalServerErr("could not load templates for today").Wrap(err)
			}

			for _, t := range templates {
				entry, err := ww.service.CreateHabitEntry(user, t.Weekday, t.Habit)
				if err != nil {
					return newInternalServerErr("could not create habit entry for today").Wrap(err)
				}
				habitz = append(habitz, entry)
			}
		}

		userHabitz := &habitState{
			Name:   user,
			Habitz: habitz,
		}
		allHabitz = append(allHabitz, userHabitz)
	}

	type stateRender struct {
		Width   float32
		Today   string
		Weekday string
		Updated string
		States  []*habitState
	}

	states := stateRender{
		Width:   100.0 / float32(len(allHabitz)),
		Today:   internal.Today(),
		Weekday: internal.Weekday(),
		Updated: time.Now().Format("15:04:05"), // Update time
		States:  allHabitz,
	}
	// Load the data into the template

	writeHTML(w, http.StatusOK, ww.todayTemplate, states)
	return nil
}

func (ww *www) newHabit(w http.ResponseWriter, r *http.Request) error {
	writeHTML(w, http.StatusOK, ww.newTemplate, nil)
	return nil
}

func (ww *www) updateHabit(w http.ResponseWriter, r *http.Request) error {
	user := chi.URLParam(r, "user")

	if user == "" {
		return newBadRequestErr("missing user")
	}

	userHabitz, err := ww.service.Templates(user)
	if err != nil {
		return newInternalServerErr("could not find user schedule").Wrap(err)
	}

	type wd struct {
		Name    string
		Enabled bool
	}
	type habit struct {
		Name     string
		Habit    string
		Weekdays []*wd
	}
	type stateRender struct {
		Habitz []*habit
	}

	// Convenience
	dayEnabled := func(day string, days []string) bool {
		for _, d := range days {
			if d == day {
				return true
			}
		}
		return false
	}

	weekdays := []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}
	state := []*habit{}
	for _, uh := range userHabitz {
		s := habit{
			Name:     user,
			Habit:    uh.Habit,
			Weekdays: make([]*wd, 7),
		}

		// Mark which days this habit is enabled
		for idx, day := range weekdays {
			s.Weekdays[idx] = &wd{Name: day, Enabled: dayEnabled(day, uh.Weekdays)}
		}
		state = append(state, &s)
	}

	writeHTML(w, http.StatusOK, ww.scheduleTemplate, state)
	return nil
}
