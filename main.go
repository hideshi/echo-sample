package main

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
)

// Config struct
type Config struct {
	Auth  AuthConfig
	GMail GMailConfig
}

// AuthConfig struct
type AuthConfig struct {
	ActivationSalt            string
	ExpirationOfActivationKey int64
}

// GMailConfig struct
type GMailConfig struct {
	SenderAddress  string
	SenderPassword string
}

var conf Config

// User struct
type User struct {
	ID            int64  `json:"id" form:"id" query:"id"`
	Email         string `json:"email" form:"email" query:"email"`
	Password      string `json:"-" form:"password"`
	Activated     int64  `json:"activated"`
	ActivationKey string `json:"-" query:"activation_key"`
}

func createConnection() *sql.DB {
	db, err := sql.Open("sqlite3", "./sample.db")
	if err != nil {
		panic(err)
	}
	return db
}

func initDB() {
	db := createConnection()
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

func createUser(c echo.Context) error {
	db := createConnection()
	defer db.Close()
	h := sha256.New()

	unixtime := getCurrentUnixTime()

	h.Write([]byte(c.FormValue("email") + c.FormValue("password") + conf.Auth.ActivationSalt + unixtime))
	activationKey := fmt.Sprintf("%x", h.Sum(nil))
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
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		c.FormValue("email"),
		c.FormValue("password"),
		activationKey,
		getExpirationOfActivationKey(),
	)
	if err != nil {
		log.Fatal(err)
	}

	lastIntertedID, err := res.LastInsertId()

	user, err := findUser(lastIntertedID)
	if err != nil {
		log.Fatal(err)
	}

	err = sendActivationMail(user)
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, lastIntertedID)
}

func sendActivationMail(user User) error {
	auth := smtp.PlainAuth(
		"",
		conf.GMail.SenderAddress,
		conf.GMail.SenderPassword,
		"smtp.gmail.com",
	)

	activationLink := `
	http://localhost:1323/users/activate?activation_key=` + user.ActivationKey

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		conf.GMail.SenderAddress,
		[]string{user.Email},
		[]byte(activationLink),
	)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func activateUser(c echo.Context) error {
	db := createConnection()
	defer db.Close()

	unixtime := getCurrentUnixTime()

	stmt, err := db.Prepare(`
		UPDATE users
			SET activated = 1
			WHERE activation_key = ?
			AND expiration_of_activation_key >= ?
		`)
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		c.QueryParam("activation_key"),
		unixtime,
	)
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, affected)
}

func getUser(c echo.Context) error {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	user, err := findUser(userID)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, user)
}

func findUser(userID int64) (User, error) {
	db := createConnection()
	defer db.Close()
	user := User{}
	err := db.QueryRow(
		`SELECT id, email, activated, activation_key FROM users WHERE id = ?`,
		userID,
	).Scan(&user.ID, &user.Email, &user.Activated, &user.ActivationKey)
	return user, err
}

func getCurrentUnixTime() string {
	return strconv.Itoa(int(time.Now().Unix()))
}

func getExpirationOfActivationKey() string {
	return strconv.Itoa(int(time.Now().Unix() + conf.Auth.ExpirationOfActivationKey))
}

func main() {
	// Create tables
	initDB()

	// Create app
	e := echo.New()

	// Load config
	if _, err := toml.DecodeFile("config.toml", &conf); err != nil {
		e.Logger.Fatal(err)
	}
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/users/:id", getUser)
	e.POST("/users", createUser)
	e.GET("/users/activate", activateUser)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
