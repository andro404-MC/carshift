package handler

import (
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"carshift/internal/db"
	"carshift/internal/template"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	template.Login().Render(r.Context(), w)
}

func HandlePostLogin(w http.ResponseWriter, r *http.Request) {
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
	password := r.FormValue("password")

	// Check if all data is provided
	if username == "" || password == "" {
		template.AlertError("missing data").Render(r.Context(), w)
		return
	}

	// Cheching if username exists
	e, err := db.IsUserExists(username)
	if err != nil {
		log.Printf("DB: Error checking user existence: %v", err)
		template.AlertError("internal error").Render(r.Context(), w)
		return
	}
	if !e {
		template.AlertError("username or password incorrect").Render(r.Context(), w)
		return
	}

	// Fetching user
	u := db.User{Username: username}
	err = u.FetchUser()
	if err != nil {
		log.Printf("SERVER: Error fetching user %v", err)
		template.AlertError("internal error").Render(r.Context(), w)
		return
	}

	// Cheching password
	if bcrypt.CompareHashAndPassword([]byte(u.Passhash), []byte(password)) != nil {
		template.AlertError("username or password incorrect").Render(r.Context(), w)
		return
	}

	// Session renew
	err = SM.RenewToken(r.Context())
	if err != nil {
		log.Printf("SCS: Error session renew %v", err)
		template.AlertError("internal error").Render(r.Context(), w)
		return
	}

	// Registering session
	SM.Put(r.Context(), "userId", u.Id)
	w.Header().Set("HX-Redirect", "/")
}
