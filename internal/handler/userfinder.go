package handler

import (
	"log"
	"net/http"

	"github.com/untemi/carshift/internal/db"
	"github.com/untemi/carshift/internal/template"
)

func GETuserFinder(w http.ResponseWriter, r *http.Request) {
	template.UserFinder().Render(r.Context(), w)
}

func POSTuserFinder(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		template.AlertError("bad data").Render(r.Context(), w)
		return
	}

	query := r.FormValue("username")
	if query == "" {
		template.AlertError("missing data").Render(r.Context(), w)
	}

	query = "%" + query + "%"
	users, err := db.FetchUsers(query, 10, 0)
	if err != nil {
		log.Printf("DB: Error fetching users: %v", err)
		template.AlertError("internal error").Render(r.Context(), w)
		return
	}

	template.UserFinderResults(users).Render(r.Context(), w)
}
