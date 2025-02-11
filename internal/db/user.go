package db

import (
	"database/sql"

	"carshift/internal/misc"
)

type User struct {
	Id        int
	Username  string
	Firstname string
	Lastname  string
	Passhash  string
	Phone     string
	Email     string
}

func IsUserExists(username string) (bool, error) {
	var userID int
	err := DB.QueryRow("SELECT user_id FROM user WHERE user_name=$1", username).Scan(&userID)
	// DB error
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	// user exists
	if err != sql.ErrNoRows {
		return true, nil
	}

	// user dont exists
	return false, nil
}

func (u *User) AddUser() error {
	err := DB.QueryRow(`
  INSERT INTO user
    (user_name,user_firstname,user_lastname,user_passhash,user_phone,user_email)
    VALUES ($1,$2,$3,$4,$5,$6)
    RETURNING user_id`,
		u.Username, u.Firstname, u.Lastname, u.Passhash, u.Phone, u.Email,
	).Scan(&u.Id)

	return err
}

func (u *User) FetchUser() error {
	if u.Username != "" {
		err := DB.QueryRow("SELECT * FROM user WHERE user_name=$1", u.Username).
			Scan(&u.Id, &u.Username, &u.Firstname, &u.Lastname, &u.Passhash, &u.Phone, &u.Email)

		return err
	} else if u.Id != 0 {
		err := DB.QueryRow("SELECT * FROM user WHERE user_id=$1", u.Id).
			Scan(&u.Id, &u.Username, &u.Firstname, &u.Lastname, &u.Passhash, &u.Phone, &u.Email)

		return err
	}

	return misc.ErrNoIdentifier
}
