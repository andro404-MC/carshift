package main

import (
	"flag"
	"log"
	"net/http"

	"carshift/internal/db"
	h "carshift/internal/handler"
	m "carshift/internal/middleware"
	"carshift/internal/view"
)

func main() {
	adr := flag.String("a", ":8000", "address")
	flag.Parse()

	db.Setup()
	defer db.Close()
	h.Setup()

	router := http.NewServeMux()
	stackMain := m.Stack(
		m.Log,
		h.SM.LoadAndSave,
	)
	stackLogged := m.Stack(
		m.FetchLogin,
		m.UserOnly,
	)
	stackGuest := m.Stack(
		m.FetchLogin,
		m.GuestOnly,
	)

	router.HandleFunc("GET /favicon.ico", view.ServeFavicon)
	router.HandleFunc("GET /static/", view.ServeStaticFiles)
	router.HandleFunc("GET /logout", h.EndSession)

	// Everyones routes
	router.Handle("GET /", m.FetchLogin(http.HandlerFunc(h.HandleHome)))
	router.Handle("GET /profile/{username}", m.FetchLogin(http.HandlerFunc(h.HandleProfile)))

	// User routes
	router.Handle("GET /me", stackLogged(http.HandlerFunc(h.HandleProfileSelf)))

	// Guest routes
	router.Handle("GET /login", stackGuest(http.HandlerFunc(h.HandleLogin)))
	router.Handle("GET /register", stackGuest(http.HandlerFunc(h.HandleRegister)))

	router.Handle("POST /login", stackGuest(http.HandlerFunc(h.HandlePostLogin)))
	router.Handle("POST /register", stackGuest(http.HandlerFunc(h.HandlePostRegister)))

	server := http.Server{
		Addr:    *adr,
		Handler: stackMain(router),
	}

	log.Println("SERVER: running on", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Printf("SERVER: Error fetching user %v", err)
	}
}
