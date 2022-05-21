package endpoints

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/jfernstad/habitz/web/internal"
	"github.com/jfernstad/habitz/web/internal/auth"
	"github.com/jfernstad/habitz/web/internal/repository"
)

type habitz struct {
	DefaultEndpoint
	service     internal.HabitzServicer
	authService auth.JWTServicer
}

func NewHabitzEndpoint(hs internal.HabitzServicer, js auth.JWTServicer) EndpointRouter {
	return &habitz{
		service:     hs,
		authService: js,
	}
}

type habitState struct {
	UserID   string                   `json:"user_id"`
	TypeName string                   `json:"type_name"`
	Habitz   []*repository.HabitEntry `json:"habitz"`
}

func (h *habitz) Routes() chi.Router {
	router := NewRouter()

	router.Use(JWTValidation(h.authService))
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

	// firstname := r.Context().Value(ContextFirstnameKey).(string)
	userID := r.Context().Value(ContextUserIDKey).(string)

	ht := repository.WeekHabitTemplates{}
	if err := json.NewDecoder(r.Body).Decode(&ht); err != nil {
		return newBadRequestErr("invalid input").Wrap(err)
	}

	thisWeekday := internal.Weekday()

	// Create Habit template
	for _, weekday := range ht.Weekdays {
		if err := h.service.CreateTemplate(userID, weekday, ht.Habit); err != nil {
			return newInternalServerErr("could not create template").Wrap(err)
		}

		// If we're adding a habit for today, make sure we use it today!
		if weekday == thisWeekday {
			// Ignore this error, less important
			h.service.CreateHabitEntry(ht.UserID, weekday, ht.Habit)
		}
	}

	writeJSON(w, http.StatusCreated, nil)
	return nil
}

func (h *habitz) deleteHabit(w http.ResponseWriter, r *http.Request) error {

	ht := repository.WeekdayHabitTemplate{}
	if err := json.NewDecoder(r.Body).Decode(&ht); err != nil {
		return newBadRequestErr("invalid input").Wrap(err)
	}

	if err := h.service.RemoveTemplate(ht.UserID, ht.Weekday, ht.Habit); err != nil {
		return newInternalServerErr("could not remove template").Wrap(err)
	}

	weekday := internal.Weekday()

	// If we're removing todays Habit
	// Also delete todays entry
	if weekday == ht.Weekday {
		h.service.RemoveEntry(ht.UserID, ht.Habit, time.Now())
	}

	writeJSON(w, http.StatusOK, nil)
	return nil
}

func (h *habitz) loadTodaysHabitz(w http.ResponseWriter, r *http.Request) error {

	firstname := r.Context().Value(ContextFirstnameKey).(string)
	userID := r.Context().Value(ContextUserIDKey).(string)

	// What day is it?
	today := internal.Today()
	weekday := internal.Weekday()

	response := []habitState{}

	// We can show multiple habitz per day
	// for multiple types of habitz.
	// For now it's just a single habit type
	// with the same name as the user
	allTypes := []string{firstname}

	// Try to retrive todays habitz for all users
	for _, habitType := range allTypes {
		habitz, err := h.service.HabitEntries(userID, today)
		if err != nil {
			return newInternalServerErr("could not load habitz for today").Wrap(err)
		}

		// Todays entries might not have been created yet, lets create them
		if len(habitz) == 0 {
			log.Println("No entries for today, lets create them")

			habitz = []*repository.HabitEntry{}
			templates, err := h.service.WeekdayTemplates(userID, weekday)
			if err != nil {
				return newInternalServerErr("could not load templates for today").Wrap(err)
			}

			for _, t := range templates {
				entry, err := h.service.CreateHabitEntry(userID, t.Weekday, t.Habit)
				if err != nil {
					return newInternalServerErr("could not create habit entry for today").Wrap(err)
				}
				habitz = append(habitz, entry)
			}
		}

		if len(habitz) > 0 {
			userHabitz := habitState{
				UserID:   userID,
				TypeName: habitType,
				Habitz:   habitz,
			}
			response = append(response, userHabitz)
		}
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
