package handler

import (
	"net/http"

	"carshift/internal/template"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	template.Home().Render(r.Context(), w)
}
