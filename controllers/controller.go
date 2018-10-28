package controllers

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strconv"

	"github.com/hideshi/echo-sample/models"
	"github.com/hideshi/echo-sample/structs"
	"github.com/hideshi/echo-sample/utils"
	"github.com/labstack/echo"
)

func CreateUser(c echo.Context) error {
	db := models.CreateConnection()
	defer db.Close()
	h := sha256.New()

	unixtime := utils.GetCurrentUnixTime()

	h.Write([]byte(c.FormValue("email") + c.FormValue("password") + structs.Conf.Auth.ActivationSalt + unixtime))
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
		utils.GetExpirationOfActivationKey(),
	)
	if err != nil {
		log.Fatal(err)
	}

	lastIntertedID, err := res.LastInsertId()

	user, err := models.FindUser(lastIntertedID)
	if err != nil {
		log.Fatal(err)
	}

	err = sendActivationMail(user)
	if err != nil {
		log.Fatal(err)
	}

	return c.JSON(http.StatusOK, lastIntertedID)
}

func sendActivationMail(user structs.User) error {
	auth := smtp.PlainAuth(
		"",
		structs.Conf.GMail.SenderAddress,
		structs.Conf.GMail.SenderPassword,
		"smtp.gmail.com",
	)

	fmt.Println(auth)

	activationLink := `
	http://localhost:1323/users/activate?activation_key=` + user.ActivationKey

	err := smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		structs.Conf.GMail.SenderAddress,
		[]string{user.Email},
		[]byte(activationLink),
	)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func ActivateUser(c echo.Context) error {
	affected, err := models.ActivateUser(c.QueryParam("activation_key"))
	if err != 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, affected)
}

func GetUser(c echo.Context) error {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	user, err := models.FindUser(userID)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, user)
}
