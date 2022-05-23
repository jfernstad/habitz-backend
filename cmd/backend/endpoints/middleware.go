package endpoints

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-chi/chi/middleware"
	"github.com/jfernstad/habitz/web/internal/auth"
)

type ContextKey string

const (
	ContextFirstnameKey ContextKey = "firstname"
	ContextUserIDKey    ContextKey = "user-id"
)

// ErrorHandler should decorate all HTTP WebserviceHandlers
// Convenience to convert to httpFunc
func ErrorHandler(handler WebserviceHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)
	}
}

func JWTValidation(jwtService auth.JWTServicer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Read authorization header
			authToken := r.Header.Get("Authorization")        // bearer eyJ...
			splitToken := strings.Split(authToken, "Bearer ") // The only type we support

			// Bad authorization
			if len(splitToken) != 2 {
				err := newNotAuthenticatedErr("Bearer token missing or malformed")
				rsp := errHttpResponse{
					errMsg:    *err,
					RequestID: middleware.GetReqID(r.Context()),
				}
				writeJSON(w, err.HTTPCode, rsp)
				return
			}

			bearerToken := strings.TrimSpace(splitToken[1])

			// Validate token
			ok, claims, err := jwtService.VerifyToken(bearerToken)
			// If bad, return 401
			if !ok {
				unauthErr := newNotAuthenticatedErr("could not parse Bearer token").Wrap(err)
				rsp := errHttpResponse{
					errMsg:    *unauthErr,
					RequestID: middleware.GetReqID(r.Context()),
				}
				writeJSON(w, unauthErr.HTTPCode, rsp)
				return
			}

			// If good, parse token, make variables available in context
			ctx := context.WithValue(r.Context(), ContextFirstnameKey, claims.Firstname)
			ctx = context.WithValue(ctx, ContextUserIDKey, claims.Subject)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
