package endpoints

import (
	"net/http"
)

// ErrorHandler should decorate all HTTP WebserviceHandlers
// Convenience to convert to httpFunc
func ErrorHandler(handler WebserviceHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}
