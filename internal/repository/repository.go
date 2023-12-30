package repository

import "nastenka_udalosti/internal/models"

// DatabaseRepo uchovává postgres funkce
type DatabaseRepo interface {
	AllUsers() bool
	InsertEvent(event models.Event) error
	Authenticate(email, testPassword string) (int, string, error)
	ShowEvents() ([]models.Event, error)
	ShowUserEvents(id int) ([]models.Event, error)
}
