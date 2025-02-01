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
	SM.Lifetime = 30 * time.Minute
	SM.Store = sqlite3store.New(db.DB)
	log.Println("SCS: up and running")
}

func IsLogged(ctx context.Context) bool {
	return SM.Exists(ctx, "userId")
}

func EndSession(w http.ResponseWriter, r *http.Request) {
	SM.Destroy(r.Context())
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
