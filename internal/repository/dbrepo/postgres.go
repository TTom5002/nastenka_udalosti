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

// ShowEvents zobrazí určitý počet příspěvků
func (m *postgresDBRepo) ShowEvents() ([]models.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var events []models.Event

	// TODO: Kde je limit, tak budeš moct přidávat více příspěvků na stránku a offset jakou stránku
	query := `
		select e.event_id, e.event_header, e.event_body, e.event_created_at, e.event_updated_at,
		u.user_id, u.user_lastname
		from events e
		left join users u on (e.event_author_id = u.user_id)
		order by e.event_created_at desc
		LIMIT 25 offset 0
	`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return events, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Event
		err := rows.Scan(
			&i.ID,
			&i.Header,
			&i.Body,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.User.ID,
			&i.User.LastName,
		)

		if err != nil {
			return events, err
		}
		events = append(events, i)
	}

	if err = rows.Err(); err != nil {
		return events, err
	}

	return events, nil
}

func (m *postgresDBRepo) InsertEvent(event models.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		insert into events(event_header, event_body, event_author_id, event_created_at, event_updated_at)
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

// Authenticate ověří uživatele, že je přihlášen
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

// Query až budu chtít vybírat podle role
/*
select e.event_id, e.event_header, e.event_body, e.event_created_at, e.event_updated_at,
		u.user_id, u.user_lastname
		from events e
		left join users u on (e.event_author_id = u.user_id)
		where u.user_access_level = 3
		order by e.event_created_at asc
		LIMIT 5 offset 0
*/

func (m *postgresDBRepo) ShowUserEvents(id int) ([]models.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var events []models.Event

	// TODO: Kde je limit, tak budeš moct přidávat více příspěvků na stránku a offset jakou stránku
	query := `
		select e.event_id, e.event_header, e.event_body, e.event_created_at, e.event_updated_at,
		u.user_id, u.user_lastname
		from events e
		left join users u on (e.event_author_id = u.user_id)
		where u.user_id  = $1
		order by e.event_created_at desc
		LIMIT 25 offset 0
	`

	rows, err := m.DB.QueryContext(ctx, query, id)
	if err != nil {
		return events, err
	}
	defer rows.Close()

	for rows.Next() {
		var i models.Event
		err := rows.Scan(
			&i.ID,
			&i.Header,
			&i.Body,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.User.ID,
			&i.User.LastName,
		)

		if err != nil {
			return events, err
		}
		events = append(events, i)
	}

	if err = rows.Err(); err != nil {
		return events, err
	}

	return events, nil
}

// TODO: Kde je limit, tak budeš moct přidávat více příspěvků na stránku a offset jakou stránku
// query := `
// select e.event_id, e.event_header, e.event_body, e.event_created_at, e.event_updated_at,
// u.user_id, u.user_lastname
// from events e
// left join users u on (e.event_author_id = u.user_id)
// where u.user_id  = $1
// order by e.event_created_at asc
// LIMIT 5 offset 0
// `
