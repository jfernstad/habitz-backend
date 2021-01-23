package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jfernstad/habitz/web/internal"
)

type habitz struct {
	DefaultEndpoint
	service internal.HabitzServicer
}

func NewHabitzEndpoint(hs internal.HabitzServicer) EndpointRouter {
	return &habitz{
		service: hs,
	}
}

type incomingHabitTemplate struct {
	Name     string
	Habit    string
	Weekdays []string
}

type habitState struct {
	Name   string
	Habitz []*internal.HabitEntry
}

func (h *habitz) Routes() chi.Router {
	router := NewRouter()

	router.Route("/", func(r chi.Router) {
		r.Get("/users", ErrorHandler(h.loadUsers))
		r.Post("/", ErrorHandler(h.createHabitTemplate))
		r.Delete("/", ErrorHandler(h.deleteHabit))
		r.Get("/today", ErrorHandler(h.loadTodaysHabitz))
		r.Post("/today", ErrorHandler(h.updateTodaysHabitz))
	})

	return router
}

func (h *habitz) loadUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := h.service.Users()
	if err != nil {
		return newInternalServerErr("could not load users").Wrap(err)
	}
	writeJSON(w, http.StatusOK, &users)
	return nil
}

func (h *habitz) createHabitTemplate(w http.ResponseWriter, r *http.Request) error {

	ht := incomingHabitTemplate{}
	if err := json.NewDecoder(r.Body).Decode(&ht); err != nil {
		return newBadRequestErr("invalid input").Wrap(err)
	}

	allUsers, err := h.service.Users()
	if err != nil {
		return newInternalServerErr("could not create template").Wrap(err)
	}

	userExist := false
	// Check if we need to create this user
	// TODO: Check in db instead
	for _, user := range allUsers {
		if ht.Name == user {
			userExist = true
		}
	}

	if userExist == false {
		// No such user, create
		if err := h.service.CreateUser(ht.Name); err != nil {
			return newInternalServerErr("could not create user").Wrap(err)
		}
	}

	// Create Habit template
	for _, weekday := range ht.Weekdays {
		if err := h.service.CreateTemplate(ht.Name, weekday, ht.Habit); err != nil {
			return newInternalServerErr("could not create template").Wrap(err)
		}
	}

	writeJSON(w, http.StatusCreated, nil)
	return nil
}

func (h *habitz) deleteHabit(w http.ResponseWriter, r *http.Request) error {

	ht := incomingHabitTemplate{}
	if err := json.NewDecoder(r.Body).Decode(&ht); err != nil {
		return newBadRequestErr("invalid input").Wrap(err)
	}

	for _, weekday := range ht.Weekdays {
		if err := h.service.RemoveTemplate(ht.Name, weekday, ht.Habit); err != nil {
			return newInternalServerErr("could not remove template").Wrap(err)
		}
	}

	writeJSON(w, http.StatusOK, nil)
	return nil
}

func (h *habitz) loadTodaysHabitz(w http.ResponseWriter, r *http.Request) error {

	allUsers, err := h.service.Users()
	if err != nil {
		return newInternalServerErr("could not load users").Wrap(err)
	}

	// What day is it?
	today := internal.Today()
	weekday := internal.Weekday()

	response := []habitState{}

	// Try to retrive todays habitz for all users
	for _, user := range allUsers {
		habitz, err := h.service.HabitEntries(user, today)
		if err != nil {
			return newInternalServerErr("could not load habitz for today").Wrap(err)
		}

		// Todays entries might not have been created yet, lets create them
		if len(habitz) == 0 {
			habitz = []*internal.HabitEntry{}
			templates, err := h.service.Templates(user, weekday)
			if err != nil {
				return newInternalServerErr("could not load templates for today").Wrap(err)
			}

			for _, t := range templates {
				entry, err := h.service.CreateHabitEntry(user, weekday, t.Habit, false)
				if err != nil {
					return newInternalServerErr("could not create habit entry for today").Wrap(err)
				}
				habitz = append(habitz, entry)
			}
		}
		userHabitz := habitState{
			Name:   user,
			Habitz: habitz,
		}
		response = append(response, userHabitz)
	}
	writeJSON(w, http.StatusOK, &response)
	return nil
}

func (h *habitz) updateTodaysHabitz(w http.ResponseWriter, r *http.Request) error {
	hh := []habitState{}
	if err := json.NewDecoder(r.Body).Decode(&hh); err != nil {
		return newBadRequestErr("invalid input").Wrap(err)
	}

	for _, userEntries := range hh {
		for _, entry := range userEntries.Habitz {
			_, err := h.service.UpdateHabitEntry(entry.ID, entry.Complete)
			if err != nil {
				return newInternalServerErr("could not update habit entry").Wrap(err)
			}
		}
	}

	writeJSON(w, http.StatusOK, nil)
	return nil
}
