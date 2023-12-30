package dbrepo

import (
	"context"
	"errors"
	"nastenka_udalosti/internal/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// TODO: Udělej všechny databázové dotazy

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

func (m *postgresDBRepo) InsertEvent(event models.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		insert into event(event_header, event_body, event_author_id, event_created_at, event_updated_at)
		values ($1,$2,$3,$4,$5)
	`
	_, err := m.DB.ExecContext(ctx, query,
		event.Header,
		event.Body,
		event.AuthorID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}
	return nil
}

// Authenticate authenticates a user
func (m *postgresDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select user_id, user_password from users where user_email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}
