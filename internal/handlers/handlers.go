package handlers

import (
	"log"
	"nastenka_udalosti/internal/config"
	"nastenka_udalosti/internal/driver"
	"nastenka_udalosti/internal/forms"
	"nastenka_udalosti/internal/models"
	"nastenka_udalosti/internal/render"
	"nastenka_udalosti/internal/repository"
	"nastenka_udalosti/internal/repository/dbrepo"
	"net/http"
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

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// MakeEvent ukáže formulář pro vytvoření příspěvku - HOTOVO
func (m *Repository) MakeEvent(w http.ResponseWriter, r *http.Request) {
	var emptyEvent models.Event
	data := make(map[string]interface{})
	data["reservation"] = emptyEvent

	render.Template(w, r, "make-event.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (m *Repository) PostEvent(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Nelze parsovat formulář!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// TODO: Musíš kurva přijít na to co má být v té závorce (probubly budeš muset vytáhnout id ze session od přihlášeného uživatele)
	// authorID, err := strconv.Atoi("")
	// if err != nil {
	// 	m.App.Session.Put(r.Context(), "error", "Nepovedlo se ověřit uživatele!")
	// 	http.Redirect(w, r, "/", http.StatusSeeOther)
	// 	return
	// }

	authorID := 1

	event := models.Event{
		Header:   r.Form.Get("header"),
		Body:     r.Form.Get("body"),
		AuthorID: authorID,
	}

	form := forms.New(r.PostForm)

	form.Required("header", "body")
	// TODO: Délka nemusí být ale možná se bude hodit
	// form.MinLength("header", 3)
	// form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["event"] = event
		render.Template(w, r, "make-event.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	err = m.DB.InsertEvent(event)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Nelze vytvořit příspěvek!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/make-event", http.StatusSeeOther)
}

func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Úspěšně přihlášen")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
