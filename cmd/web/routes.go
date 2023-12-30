package main

import (
	"nastenka_udalosti/internal/config"
	"nastenka_udalosti/internal/handlers"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)

	mux.Get("/user/login", handlers.Repo.Login)
	mux.Post("/user/login", handlers.Repo.PostLogin)

	// mux.Get("/user/signup", handlers.Repo.Signup)
	// mux.Post("/user/signup", handlers.Repo.PostSignup)

	mux.Get("/user/logout", handlers.Repo.Logout)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Route("/admin", func(mux chi.Router) {
		// TODO: Zruš komentář
		mux.Use(Auth)
		mux.Get("/make-event", handlers.Repo.MakeEvent)
		mux.Post("/make-event", handlers.Repo.PostMakeEvent)

		mux.Get("/my-events", handlers.Repo.MyEvents)
	})

	return mux
}
