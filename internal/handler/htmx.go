package handler

import (
	"net/http"

	"github.com/untemi/carshift/internal/template"
)

func HtmxAlert(w http.ResponseWriter, r *http.Request) {
	message := r.URL.Query().Get("message")
	template.AlertError(message).Render(r.Context(), w)
}
