package endpoints

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt"
	"github.com/jfernstad/habitz/web/internal"
)

type authEndpoint struct {
	DefaultEndpoint
	service           internal.HabitzServicer
	jwtSigningSecret  []byte
	jwtAudience       string
	cachedGoogleCerts map[string]string
}

type googleClaims struct {
	jwt.StandardClaims
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	ProfileImage  string `json:"picture"`
}

type HabitzJWTClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	UserID   string `json:"user_id"`
}

func NewAuthEndpoint(hs internal.HabitzServicer, jwtSigningSecret string, jwtAudience string) EndpointRouter {
	return &authEndpoint{
		service:           hs,
		jwtSigningSecret:  []byte(jwtSigningSecret), // TODO: verify its 32 bytes long
		jwtAudience:       jwtAudience,              // Verify the incoming JWT token was intended for us
		cachedGoogleCerts: map[string]string{},      // Optimization
	}
}

func (a *authEndpoint) Routes() chi.Router {
	router := NewRouter()

	router.Route("/", func(r chi.Router) {
		r.Post("/google", ErrorHandler(a.google))
	})

	return router
}

func (a *authEndpoint) google(w http.ResponseWriter, r *http.Request) error {
	// User fetched a JQT token from Google
	// Lets transform it to our JWT token
	defer r.Body.Close()

	type token struct {
		Token string `json:"token"`
	}

	loginToken := token{}

	err := json.NewDecoder(r.Body).Decode(&loginToken)
	if err != nil {
		return newBadRequestErr("not a valid Google JWT").Wrap(err)
	}

	gToken, err := a.parseGoogleJWTToken(loginToken.Token)
	if err != nil {
		return newBadRequestErr("could not validate Google JWT").Wrap(err)
	}

	// TODO:
	// Is this the first time the user logs in?
	// If so, create an account, store basic info

	// Create a JWT token for the API
	// 30 days it last
	expirationTime := time.Now().Add(30 * 24 * time.Hour)

	claims := &HabitzJWTClaims{
		Username: gToken.FirstName,
		UserID:   gToken.Subject,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	habitzToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the new JWT string
	tokenString, err := habitzToken.SignedString(a.jwtSigningSecret)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		return newBadRequestErr("could not sign habitz JWT").Wrap(err)
	}

	resp := token{
		Token: tokenString,
	}

	writeJSON(w, http.StatusOK, resp)

	return nil
}

func (a *authEndpoint) parseGoogleJWTToken(tokenString string) (*googleClaims, error) {
	claimsStruct := googleClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) {
			pem, err := a.getGooglePublicKey(fmt.Sprintf("%s", token.Header["kid"]))
			if err != nil {
				return nil, err
			}
			key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
			if err != nil {
				return nil, err
			}
			return key, nil
		},
	)
	if err != nil {
		return &googleClaims{}, err
	}

	claims, ok := token.Claims.(*googleClaims)
	if !ok {
		return &googleClaims{}, errors.New("Invalid Google JWT")
	}

	if claims.Issuer != "accounts.google.com" && claims.Issuer != "https://accounts.google.com" {
		return &googleClaims{}, errors.New("iss is invalid")
	}

	// TODO: Pass Google App ClientID into Auth service
	if claims.Audience != a.jwtAudience {
		return &googleClaims{}, errors.New("aud is invalid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return &googleClaims{}, errors.New("JWT is expired")
	}

	return claims, nil
}

// From: https://blog.boot.dev/golang/how-to-implement-sign-in-with-google-in-golang/
func (a *authEndpoint) getGooglePublicKey(keyID string) (string, error) {

	// Check cache first
	if key, ok := a.cachedGoogleCerts[keyID]; ok {
		return key, nil
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
	if err != nil {
		return "", err
	}
	dat, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	myResp := map[string]string{}
	err = json.Unmarshal(dat, &myResp)
	if err != nil {
		return "", err
	}
	key, ok := myResp[keyID]
	if !ok {
		return "", errors.New("key not found")
	}

	// Cache key
	a.cachedGoogleCerts[keyID] = key
	return key, nil
}
