package handler

import (
	"log"
	"net/http"

	"carshift/internal/db"
	"carshift/internal/template"
)

func HandleProfileSelf(w http.ResponseWriter, r *http.Request) {
	u, ok := r.Context().Value("userdata").(db.User)
	if !ok {
		log.Println("SERVER: error fetching prop logged")
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	template.Profile(u, true).Render(r.Context(), w)
}
