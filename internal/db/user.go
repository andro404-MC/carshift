package db

import "carshift/internal/misc"

func IsUserExists(username string) (bool, error) {
	tx := db.Limit(1).Where("Username = ?", username).Find(&User{})
	return tx.RowsAffected > 0, tx.Error
}

func (u *User) AddUser() error {
	return db.Create(&u).Error
}

func (u *User) FetchUser() error {
	if u.Username != "" {
		return db.Where("Username = ?", u.Username).First(&u).Error
	} else if u.ID != 0 {
		return db.First(&u, u.ID).Error
	}

	return misc.ErrNoIdentifier
}
