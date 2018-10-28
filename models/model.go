package models

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/hideshi/echo-sample/structs"
	"github.com/hideshi/echo-sample/utils"
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

func ActivateUser(activationKey string) (int64, int64) {
	db := CreateConnection()
	defer db.Close()

	unixtime := utils.GetCurrentUnixTime()

	stmt, err := db.Prepare(`
	UPDATE users
		SET activated = 1
		WHERE activation_key = ?
		AND expiration_of_activation_key >= ?
	`)
	if err != nil {
		log.Fatal(err)
		return 0, http.StatusInternalServerError
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		activationKey,
		unixtime,
	)
	if err != nil {
		log.Fatal(err)
		return 0, http.StatusInternalServerError
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return 0, http.StatusNotFound
	}

	return affected, 0
}
