package handler

import (
	"net/http"
	"strconv"

	"carshift/internal/template"
)

var tabs = []template.Tab{
	{Name: "Profile", Content: template.SettingsProfie(), URL: "/settings/0"},
	{Name: "Account", Content: template.AlertError("Account"), URL: "/settings/1"},
}

func HandleSettings(w http.ResponseWriter, r *http.Request) {
	template.Settings().Render(r.Context(), w)
}

func HandleSettingsTabs(w http.ResponseWriter, r *http.Request) {
	// cool but need better
	isHTMX := r.Header.Get("HX-Request")
	if isHTMX != "true" {
		http.Error(w, "Not allowed", http.StatusMethodNotAllowed)
		return
	}

	sel, err := strconv.Atoi(r.PathValue("tab"))
	if err != nil || len(tabs) < sel+1 {
		http.NotFound(w, r)
		return
	}

	template.Tabbed(tabs, sel, "#settings-tabs").Render(r.Context(), w)
}
