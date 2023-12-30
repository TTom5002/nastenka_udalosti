package repository

import "nastenka_udalosti/internal/models"

// TODO: Udělej koment a sem vkládej Postgres funkce
type DatabaseRepo interface {
	AllUsers() bool
	InsertEvent(event models.Event) error
	Authenticate(email, testPassword string) (int, string, error)
}
