package handler

import (
	"net/http"

	"github.com/untemi/carshift/internal/template"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	template.Home().Render(r.Context(), w)
}
