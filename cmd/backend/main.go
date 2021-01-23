package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/jfernstad/habitz/web/cmd/backend/endpoints"
	"github.com/jfernstad/habitz/web/internal/sqlite"
)

func main() {

	dbFile := os.Getenv("SQLITE_DB")
	if dbFile == "" {
		dbFile = "habitz.sqlite"
	}

	db, err := sqlx.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-XSRF-TOKEN"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	// habitzService := &mock.HabitzService{}
	habitzService := sqlite.NewHabitzService(db)
	habitzEndpoint := endpoints.NewHabitzEndpoint(habitzService)

	r := endpoints.NewRouter()

	r.Route("/api/habitz", func(v chi.Router) {
		v.Use(cors.Handler)
		v.Mount("/", habitzEndpoint.Routes())
	})

	log.Println("HTTP routes:")
	printRoutes(r)

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}

func printRoutes(routes chi.Routes) {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("%s \t%s\n", method, route)
		return nil
	}

	if err := chi.Walk(routes, walkFunc); err != nil {
		fmt.Printf("printRoutes error: %s\n", err.Error())
	}
}
