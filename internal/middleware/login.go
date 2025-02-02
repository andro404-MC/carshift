package middleware

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	"carshift/internal/db"
	"carshift/internal/handler"
)

func CheckLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "logged", handler.IsLogged(r.Context()))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GuestOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l, ok := r.Context().Value("logged").(bool)
		if !ok {
			log.Println("SERVER: error fetching prop logged")
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

		if l {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func UserOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l, ok := r.Context().Value("logged").(bool)
		if !ok {
			log.Println("SERVER: error fetching prop logged")
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

		if !l {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		u := db.User{Id: handler.SM.GetInt(r.Context(), "userId")}

		err := db.GetUser(&u)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Redirect(w, r, "/logout", http.StatusSeeOther)
				return
			}

			log.Printf("SERVER: Error fetching user %v", err)
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), "userdata", u)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
