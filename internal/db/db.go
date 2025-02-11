package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Setup() {
	var err error
	DB, err = sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}

	// Sessions
	_, err = DB.Exec(`
  CREATE TABLE IF NOT EXISTS sessions (
	  token TEXT PRIMARY KEY,
	  data BLOB NOT NULL,
	  expiry REAL NOT NULL
  );`)
	if err != nil {
		log.Fatalf("Error creating table sessions: %q\n", err)
	}

	// User
	_, err = DB.Exec(`
  CREATE TABLE IF NOT EXISTS user (
	  user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_name TEXT,
    user_firstname TEXT,
    user_lastname TEXT,
    user_passhash TEXT,
    user_phone TEXT,
    user_email TEXT
  )`)
	if err != nil {
		log.Fatalf("Error creating table user: %q\n", err)
	}

	// Car
	_, err = DB.Exec(`
  CREATE TABLE IF NOT EXISTS car (
	  car_id INTEGER PRIMARY KEY AUTOINCREMENT,
    car_title TEXT,
    car_price REAL,
    post_date TEXT,
    owner INTEGER,
    FOREIGN KEY(owner) REFERENCES user(user_id)
  );`)
	if err != nil {
		log.Fatalf("Error creating table sessions: %q\n", err)
	}

	// Indexes
	_, err = DB.Exec("CREATE INDEX IF NOT EXISTS sessions_expiry_idx ON sessions(expiry)")
	if err != nil {
		log.Fatalf("Error creating index sessions_expiry_idx: %q\n", err)
	}

	_, err = DB.Exec("CREATE INDEX IF NOT EXISTS user_name_idx ON user(user_name)")
	if err != nil {
		log.Fatalf("Error creating index user_name_idx: %q\n", err)
	}

	_, err = DB.Exec("CREATE INDEX IF NOT EXISTS car_owner_idx ON car(owner)")
	if err != nil {
		log.Fatalf("Error creating index user_name_idx: %q\n", err)
	}

	log.Println("DB: up and running")
}

func Close() {
	DB.Close()
}
