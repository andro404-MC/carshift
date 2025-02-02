package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"

	"carshift/internal/db"
)

var SM *scs.SessionManager

func Setup() {
	SM = scs.New()

	SM.Lifetime = time.Hour * 365 * 24
	SM.Store = sqlite3store.New(db.DB)

	log.Println("SCS: up and running")
}

func IsLogged(ctx context.Context) bool {
	return SM.Exists(ctx, "userId")
}

func EndSession(w http.ResponseWriter, r *http.Request) {
	err := SM.RenewToken(r.Context())
	if err != nil {
		log.Printf("SCS: Error session renew %v", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	SM.Destroy(r.Context())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
