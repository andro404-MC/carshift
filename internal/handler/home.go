package handler

import (
	"net/http"

	"carshift/internal/template"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	template.Home().Render(r.Context(), w)
}
