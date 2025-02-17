package middleware

import (
	"context"
	"log"
	"net/http"

	"gorm.io/gorm"

	"github.com/untemi/carshift/internal/db"
	"github.com/untemi/carshift/internal/handler"
)

func FetchLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := handler.IsLogged(r.Context())
		ctx := context.WithValue(r.Context(), "logged", l)

		if l {
			u := db.User{ID: handler.SM.GetInt(r.Context(), "userId")}

			err := u.Fill()
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					http.Redirect(w, r, "/logout", http.StatusSeeOther)
					return
				}

				log.Printf("SERVER: Error fetching user %v", err)
				http.Error(w, "Internal Error", http.StatusInternalServerError)
				return
			}

			ctx = context.WithValue(ctx, "userdata", u)
		}

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

		next.ServeHTTP(w, r)
	})
}
