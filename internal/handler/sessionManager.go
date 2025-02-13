package handler

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	_ "github.com/mattn/go-sqlite3"
)

var SM = scs.New()

func Setup() error {
	db, err := sql.Open("sqlite3", "session.db")
	if err != nil {
		return err
	}

	_, err = db.Exec(`
  CREATE TABLE IF NOT EXISTS sessions (
	  token TEXT PRIMARY KEY,
	  data BLOB NOT NULL,
	  expiry REAL NOT NULL
  );`)
	if err != nil {
		return err
	}

	if db.Exec("CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions(expiry)"); err != nil {
		return err
	}

	SM.Lifetime = time.Hour * 365 * 24
	SM.Store = sqlite3store.New(db)

	log.Println("SCS: up and running")
	return nil
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
