package main

import (
	"log"
	"net/http"

	"carshift/internal/db"
	h "carshift/internal/handler"
	m "carshift/internal/middleware"
	"carshift/internal/view"
)

func main() {
	db.Setup()
	defer db.Close()
	h.Setup()

	router := http.NewServeMux()
	stackMain := m.Stack(
		m.Log,
		h.SM.LoadAndSave,
	)
	stackLogged := m.Stack(
		m.CheckLogin,
		m.UserOnly,
	)
	stackGuest := m.Stack(
		m.CheckLogin,
		m.GuestOnly,
	)

	router.HandleFunc("GET /favicon.ico", view.ServeFavicon)
	router.HandleFunc("GET /static/", view.ServeStaticFiles)

	// User routes
	router.Handle("GET /", stackLogged(http.HandlerFunc(h.HandleHome)))
	router.Handle("GET /logout", stackLogged(http.HandlerFunc(h.EndSession)))

	// Guest routes
	router.Handle("GET /login", stackGuest(http.HandlerFunc(h.HandleLogin)))
	router.Handle("GET /register", stackGuest(http.HandlerFunc(h.HandleRegister)))
	router.Handle("POST /login", stackGuest(http.HandlerFunc(h.HandlePostLogin)))
	router.Handle("POST /register", stackGuest(http.HandlerFunc(h.HandlePostRegister)))

	server := http.Server{
		Addr:    ":8000",
		Handler: stackMain(router),
	}

	log.Println("SERVER: running on port", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
