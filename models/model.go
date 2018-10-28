package models

import (
	"database/sql"
	"log"

	"github.com/hideshi/echo-sample/structs"
)

func CreateConnection() *sql.DB {
	db, err := sql.Open("sqlite3", "./sample.db")
	if err != nil {
		panic(err)
	}
	return db
}

func InitDB() {
	db := CreateConnection()
	defer db.Close()
	db.Exec(`
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

func FindUser(userID int64) (structs.User, error) {
	db := CreateConnection()
	defer db.Close()
	user := structs.User{}
	err := db.QueryRow(
		`SELECT id, email, activated, activation_key FROM users WHERE id = ?`,
		userID,
	).Scan(&user.ID, &user.Email, &user.Activated, &user.ActivationKey)
	return user, err
}

func CreateUser(email string, password string, activationKey string, expirationOfActivationKey string) (int64, error) {
	db := CreateConnection()
	defer db.Close()

	stmt, err := db.Prepare(`
		INSERT INTO users (
			email,
			password,
			activated,
			activation_key,
			expiration_of_activation_key
			) VALUES (?, ?, 0, ?, ?)
		`)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		email,
		password,
		activationKey,
		expirationOfActivationKey,
	)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	return res.LastInsertId()
}

func ActivateUser(activationKey string, unixtime string) (sql.Result, error) {
	db := CreateConnection()
	defer db.Close()

	stmt, err := db.Prepare(`
	UPDATE users
		SET activated = 1
		WHERE activation_key = ?
		AND expiration_of_activation_key >= ?
	`)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		activationKey,
		unixtime,
	)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return res, nil
}
