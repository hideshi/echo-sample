package models

import (
	"database/sql"
	"log"

	"github.com/hideshi/echo-sample/structs"
)

func CreateConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "sample.db")
	if err != nil {
		return nil, err
	}
	return db, nil
}

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

func FindUserByID(userID int64) (structs.User, error) {
	db, err := CreateConnection()
	defer db.Close()
	if err != nil {
		return structs.User{}, err
	}

	stmt, err := db.Prepare(`SELECT id, email, activated, activation_key FROM users WHERE id = ?`)
	defer stmt.Close()
	if err != nil {
		return structs.User{}, err
	}

	rows, err := stmt.Query(userID)
	user := structs.User{}
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Email, &user.Activated, &user.ActivationKey); err != nil {
			return structs.User{}, err
		}
	}
	if err := rows.Err(); err != nil {
		return structs.User{}, err
	}

	return user, err
}

func CreateUser(email string, password string, activationKey string, expirationOfActivationKey string) (int64, error) {
	db, err := CreateConnection()
	defer db.Close()
	if err != nil {
		return 0, err
	}

	stmt, err := db.Prepare(`
		INSERT INTO users (
			email,
			password,
			activated,
			activation_key,
			expiration_of_activation_key
			) VALUES (?, ?, 0, ?, ?)
		`)
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

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
	db, err := CreateConnection()
	defer db.Close()
	if err != nil {
		return nil, err
	}

	stmt, err := db.Prepare(`
	UPDATE users
		SET activated = 1
		WHERE activation_key = ?
		AND expiration_of_activation_key >= ?
	`)
	defer stmt.Close()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

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
