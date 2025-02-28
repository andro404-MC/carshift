package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	gm "github.com/go-chi/chi/v5/middleware"

	"github.com/untemi/carshift/internal/db"
	h "github.com/untemi/carshift/internal/handler"
	m "github.com/untemi/carshift/internal/middleware"
	"github.com/untemi/carshift/internal/view"
)

func main() {
	var err error
	var adr string

	// Flags
	flag.StringVar(&adr, "a", ":8000", "address")
	flag.Parse()

	// Setup
	if err = db.Init(); err != nil {
		log.Printf("DB: Error fetching user %v", err)
		return
	}

	if err = h.Init(); err != nil {
		log.Printf("SM: Error fetching user %v", err)
		return
	}

	r := chi.NewRouter()
	r.Use(gm.Logger, gm.Recoverer, h.SM.LoadAndSave)

	// Static and general stuff
	r.Group(func(r chi.Router) {
		r.Get("/favicon.ico", view.ServeFavicon)
	})

	// HTMX thingys
	r.Group(func(r chi.Router) {
		r.Get("/htmx/alert", h.HtmxAlert)
	})

	// OUR routes
	r.Group(func(r chi.Router) {
		r.Use(m.FetchLogin)

		r.Get("/", h.GEThome)
		r.Get("/profile/{username}", h.GETprofile)

		r.Get("/carfinder", h.GETcarFinder)
		r.Post("/carfinder", h.POSTcarFinder)

		r.Get("/userfinder", h.GETuserFinder)
		r.Post("/userfinder", h.POSTuserFinder)
	})

	// User routes
	r.Group(func(r chi.Router) {
		r.Use(m.FetchLogin, m.UserOnly)

		r.Get("/logout", h.EndSession)
		r.Get("/me", h.GETprofileSelf)
		r.Get("/settings", h.GETsettings)
		r.Get("/settings/{tab}", h.GETsettingsTabs)
	})

	// Guest routes
	r.Group(func(r chi.Router) {
		r.Use(m.FetchLogin, m.GuestOnly)

		r.Get("/login", h.GETlogin)
		r.Get("/register", h.GETregister)
		r.Post("/login", h.POSTlogin)
		r.Post("/register", h.POSTregister)
	})

	// Files serving
	view.FileServer(r, "/static", "static")

	server := http.Server{
		Addr:    adr,
		Handler: r,
	}

	log.Println("SERVER: running on", server.Addr)
	err = server.ListenAndServe()
	if err != nil {
		log.Printf("SERVER: Error fetching user %v", err)
	}
}
