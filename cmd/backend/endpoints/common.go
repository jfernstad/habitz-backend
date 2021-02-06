package endpoints

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type EndpointRouter interface {
	Routes() chi.Router
}

// DefaultNotFoundHandler always responds with a generic 404 error.
func DefaultNotFoundHandler(w http.ResponseWriter, req *http.Request) error {
	return newNotFoundErr(fmt.Sprintf("handler not found: %s", req.RequestURI))
}

// DefaultMethodNotAllowedHandler always responds with a generic 405 error.
func DefaultMethodNotAllowedHandler(w http.ResponseWriter, req *http.Request) error {
	return ErrMethodNotAllowed
}

// EndpointMethods configure a specific endpoint
// Made specially to work with the `httprouter` package
type EndpointMethods interface {
	MissingHandler() http.HandlerFunc
	MethodNotAllowedHandler() http.HandlerFunc
	Routes() chi.Router
}

// DefaultEndpoint all endpoints should `inherit`, i.e compose, this one
type DefaultEndpoint struct {
}

// MissingHandler default HTTP 404 StatusNotFound
func (e DefaultEndpoint) MissingHandler(w http.ResponseWriter, req *http.Request) {
	DefaultNotFoundHandler(w, req)
}

// MethodNotAllowedHandler default HTTP 405 StatusMethodNotAllowed
func (e DefaultEndpoint) MethodNotAllowedHandler(w http.ResponseWriter, req *http.Request) {
	DefaultMethodNotAllowedHandler(w, req)
}

// NewRouter is a convenience method to create most routes. Used here and in tests
func NewRouter() chi.Router {
	r := chi.NewRouter()

	r.NotFound(ErrorHandler(DefaultNotFoundHandler))
	r.MethodNotAllowed(ErrorHandler(DefaultMethodNotAllowedHandler))

	return r
}

// WebserviceHandler is the type of web service function we define per endpoint
type WebserviceHandler func(rw http.ResponseWriter, req *http.Request) error

func (h WebserviceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var msg *errMsg

	if err := h(w, r); err != nil {
		switch err.(type) {
		case *errMsg:
			{
				msg = err.(*errMsg)
			}
		default: // Internal service error
			{
				msg = newInternalServerErr("internal error").Wrap(err)
			}
		}
		fmt.Println(msg)

		rsp := struct {
			errMsg
			RequestID string `json:"requestId"`
		}{
			errMsg:    *msg,
			RequestID: middleware.GetReqID(r.Context()),
		}

		writeJSON(w, rsp.HTTPCode, rsp)
	}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)

	if v == nil {
		v = struct{}{} // This empty struct will be encoded as "{}"
	}

	// Don't encode any response on this status
	if status == http.StatusNoContent {
		return
	}

	if err := json.NewEncoder(w).Encode(v); err != nil {
		// It's too late to tell the client--we already sent the status code.
		// The best we can do is log the error.
		log.Println("json encode error: ", err)
	}
}

func writeHTML(w http.ResponseWriter, status int, htmlTemplate *template.Template, content interface{}) {
	w.Header().Set("content-type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	// Don't encode any response on this status
	if status == http.StatusNoContent {
		return
	}

	if err := htmlTemplate.Execute(w, content); err != nil {
		// It's too late to tell the client--we already sent the status code.
		// The best we can do is log the error.
		log.Println("html encode error: ", err)
	}
}

func readJSON(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return WrapJSONDecodeError(err)
	}
	return nil
}
