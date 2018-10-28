package models

import (
	"github.com/hideshi/echo-sample/structs"
	"github.com/jinzhu/gorm"
)

func CreateConnection() (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", "sample.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

/*
func InitDB() (sql.Result, error) {
	db, err := CreateConnection()
	defer db.Close()
	if err != nil {
		return nil, err
	}

	return db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL,
			password TEXT NOT NULL,
			activated NUMBER NOT NULL,
			activation_key TEXT NOT NULL,
			expiration_of_activation_key TEXT NOT NULL
		)
	`)
}
*/
func FindUserByID(userID uint64) (structs.User, error) {
	user := structs.User{}
	db, err := CreateConnection()
	defer db.Close()
	if err != nil {
		return user, err
	}
	db.First(&user, userID)
	return user, nil
}

func CreateUser(email string, password string, activationKey string, expirationOfActivationKey string) (structs.User, error) {
	db, err := CreateConnection()
	defer db.Close()
	if err != nil {
		return structs.User{}, err
	}
	user := structs.User{
		Email:                     email,
		Password:                  password,
		ActivationKey:             activationKey,
		ExpirationOfActivationKey: expirationOfActivationKey,
	}
	db.Create(&user)
	return user, nil
}

func ActivateUser(activationKey string, unixtime string) (structs.User, error) {
	user := structs.User{}
	db, err := CreateConnection()
	defer db.Close()
	if err != nil {
		return user, err
	}
	db.Model(&user).Where("activation_key = ?", activationKey).Where("expiration_of_activation_key >= ?", unixtime).Update("activated", 1)

	return user, nil
}

/*
func UpdateEmail(userID uint64, email string) (structs.User, error) {
	db, err := CreateConnection()
	defer db.Close()
	if err != nil {
		return structs.User{}, err
	}

	stmt, err := db.Prepare(`
		UPDATE users
		SET email = ?
		WHERE id = ?
	`)
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
		return structs.User{}, err
	}

	_, err2 := stmt.Exec(email, userID)
	if err2 != nil {
		log.Fatal(err)
		return structs.User{}, err2
	}

	stmt3, err3 := db.Prepare(`SELECT id, email, activated FROM users WHERE id = ?`)
	defer stmt3.Close()
	if err3 != nil {
		return structs.User{}, err3
	}
	rows, err := stmt3.Query(userID)

	user := structs.User{}
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Email, &user.Activated); err != nil {
			return structs.User{}, err
		}
	}
	if err != nil {
		log.Fatal(err)
		return structs.User{}, err
	}

	return user, nil
}
*/
