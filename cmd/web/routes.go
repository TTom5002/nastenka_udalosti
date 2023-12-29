package main

import (
	"nastenka_udalosti/internal/config"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	// TODO: Odstraň příklad
	// mux.Get("/search-availability", handlers.Repo.Availability)
	// mux.Post("/search-availability", handlers.Repo.PostAvailability)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Route("/admin", func(mux chi.Router) {
		// TODO: Zruš komentář
		// mux.Use(Auth)

		// TODO: Odstraň příklad
		// mux.Get("/reservations-calendar", handlers.Repo.AdminReservationsCalendar)
		// mux.Post("/reservations-calendar", handlers.Repo.AdminPostReservationsCalendar)

	})

	return mux
}
