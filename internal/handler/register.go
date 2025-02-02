package handler

import (
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"carshift/internal/db"
	"carshift/internal/misc"
	"carshift/internal/template"
)

var loginDisable = 0

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	template.Register().Render(r.Context(), w)
}

func HandlePostRegister(w http.ResponseWriter, r *http.Request) {
	if loginDisable == 1 {
		template.AlertError("Login Disabled").Render(r.Context(), w)
		return
	}

	// DPS parses
	err := r.ParseForm()
	if err != nil {
		template.AlertError("bad data").Render(r.Context(), w)
		return
	}

	username := strings.TrimSpace(r.FormValue("username"))
	firstname := strings.TrimSpace(r.FormValue("firstname"))
	lastname := strings.TrimSpace(r.FormValue("lastname"))
	password := r.FormValue("password")
	passwordcon := r.FormValue("passwordcon")

	// Check if all data is provided
	if username == "" || firstname == "" || lastname == "" || password == "" || passwordcon == "" {
		template.AlertError("missing data").Render(r.Context(), w)
		return
	}

	// Verifying username
	if !misc.ValidateUsername(username) {
		template.AlertError("invalid username").Render(r.Context(), w)
		return
	}

	// Convert names to proper format
	firstname = misc.FormaterName(firstname)
	lastname = misc.FormaterName(lastname)

	// Verifying first & last name
	if !misc.ValidateName(firstname) || !misc.ValidateName(lastname) {
		template.AlertError("invalid first or lastname").Render(r.Context(), w)
		return
	}

	// Verifying Password
	if passwordcon != password {
		template.AlertError("passwords do not match").Render(r.Context(), w)
		return
	}
	if !misc.ValidatePassword(password) {
		template.AlertError("invalid password").Render(r.Context(), w)
		return
	}

	// Cheching if username is used
	e, err := db.IsUserExists(username)
	if err != nil {
		log.Printf("DB: Error checking user existence: %v", err)
		template.AlertError("internal error").Render(r.Context(), w)
		return
	}
	if e {
		template.AlertError("username is already taken").Render(r.Context(), w)
		return
	}

	// Hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("SERVER: Error hashing password: %v", err)
		template.AlertError("internal error").Render(r.Context(), w)
		return
	}

	// Registering the user
	u := db.User{
		Username:  username,
		Firstname: firstname,
		Lastname:  lastname,
		Passhash:  string(hashedPassword),
	}
	err = db.AddUser(&u)
	if err != nil {
		log.Printf("DB: Error creating user: %v", err)
		template.AlertError("internal error").Render(r.Context(), w)
		return
	}

	SM.Put(r.Context(), "userId", u.Id)
	w.Header().Set("HX-Redirect", "/")
}
