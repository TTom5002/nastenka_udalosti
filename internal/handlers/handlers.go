package handlers

import (
	"nastenka_udalosti/internal/config"
	"nastenka_udalosti/internal/driver"
	"nastenka_udalosti/internal/repository"
	"nastenka_udalosti/internal/repository/dbrepo"
)

// TODO: Změň komentáře

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// TODO: Odkomentuj až budeš testovat
// // NewTestRepo creates a new repository
// func NewTestRepo(a *config.AppConfig) *Repository {
// 	return &Repository{
// 		App: a,
// 		DB:  dbrepo.NewTestingRepo(a),
// 	}
// }

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}
