package endpoints

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jfernstad/habitz/web/internal"
)

type www struct {
	DefaultEndpoint
	service internal.HabitzServicer
}

func NewWWWEndpoint(hs internal.HabitzServicer) EndpointRouter {
	return &www{
		service: hs,
	}
}

func (ww *www) Routes() chi.Router {
	router := NewRouter()

	router.Route("/", func(r chi.Router) {
		r.Get("/", ErrorHandler(ww.todaysHabitz))
		r.Get("/new", ErrorHandler(ww.newHabit))
	})

	return router
}

func (ww *www) todaysHabitz(w http.ResponseWriter, r *http.Request) error {

	// Load the template
	htmlTemplate, err := template.ParseFiles("./cmd/backend/templates/today.tmpl")
	if err != nil {
		return newInternalServerErr("could not create template").Wrap(err)
	}

	// Load habits for all users
	allUsers, err := ww.service.Users()
	if err != nil {
		return newInternalServerErr("could not load users").Wrap(err)
	}

	// What day is it?
	today := internal.Today()
	allHabitz := []*habitState{}

	// Try to retrive todays habitz for all users
	for _, user := range allUsers {
		habitz, err := ww.service.HabitEntries(user, today)
		if err != nil {
			return newInternalServerErr("could not load habitz for today").Wrap(err)
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
		States  []*habitState
	}

	states := stateRender{
		Width:   100.0 / float32(len(allHabitz)),
		Today:   internal.Today(),
		Weekday: internal.Weekday(),
		States:  allHabitz,
	}
	// Load the data into the template

	writeHTML(w, http.StatusOK, htmlTemplate, states)
	return nil
}

func (ww *www) newHabit(w http.ResponseWriter, r *http.Request) error {
	// Load the template
	htmlTemplate, err := template.ParseFiles("./cmd/backend/templates/new.tmpl")
	if err != nil {
		return newInternalServerErr("could not create template").Wrap(err)
	}
	writeHTML(w, http.StatusOK, htmlTemplate, nil)
	return nil
}
