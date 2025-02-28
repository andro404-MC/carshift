package handler

import (
	"net/http"

	"github.com/untemi/carshift/internal/db"
	"github.com/untemi/carshift/internal/template"
)

func GETcarFinder(w http.ResponseWriter, r *http.Request) {
	template.CarFinder().Render(r.Context(), w)
}

func POSTcarFinder(w http.ResponseWriter, r *http.Request) {
	template.CarFinderResults(&[]db.Car{}).Render(r.Context(), w)
}
