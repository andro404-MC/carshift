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
	if err = db.Setup(); err != nil {
		log.Printf("DB: Error fetching user %v", err)
		return
	}

	if err = h.Setup(); err != nil {
		log.Printf("SM: Error fetching user %v", err)
		return
	}

	r := chi.NewRouter()
	r.Use(gm.Logger, h.SM.LoadAndSave)

	// Static and general stuff
	r.Group(func(r chi.Router) {
		r.Use(gm.Compress(5, "text/css", "text/javascript"))

		r.Get("/favicon.ico", view.ServeFavicon)
	})

	// HTMX thingys
	r.Group(func(r chi.Router) {
		r.Get("/htmx/alert", h.HtmxAlert)
	})

	// OUR routes
	r.Group(func(r chi.Router) {
		r.Use(m.FetchLogin)

		r.Get("/", h.HandleHome)
		r.Get("/profile/{username}", h.HandleProfile)
	})

	// User routes
	r.Group(func(r chi.Router) {
		r.Use(m.FetchLogin, m.UserOnly)

		r.Get("/logout", h.EndSession)
		r.Get("/me", h.HandleProfileSelf)
		r.Get("/settings", h.HandleSettings)
		r.Get("/settings/{tab}", h.HandleSettingsTabs)
	})

	// Guest routes
	r.Group(func(r chi.Router) {
		r.Use(m.FetchLogin, m.GuestOnly)

		r.Get("/login", h.HandleLogin)
		r.Get("/register", h.HandleRegister)
		r.Post("/login", h.HandlePostLogin)
		r.Post("/register", h.HandlePostRegister)
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
