package controllers

import (
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
	email := c.FormValue("email")
	password := c.FormValue("password")
	activationKey := utils.GenerateActivationKey(email, password, structs.Conf.Auth.ActivationSalt)
	expirationOfActivationKey := utils.GetExpirationOfActivationKey(structs.Conf.Auth.ExpirationOfActivationKey)
	lastIntertedID, err := models.CreateUser(
		email,
		password,
		activationKey,
		expirationOfActivationKey,
	)
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	user, err := models.FindUserByID(lastIntertedID)
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	err = sendActivationMail(user)
	if err != nil {
		log.Fatal(err)
		return c.NoContent(http.StatusInternalServerError)
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

	activationLink := fmt.Sprintf(`
	%s/users/activate?activation_key=%s`, structs.Conf.Environment.Host, user.ActivationKey)

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
	unixtime := utils.GetCurrentUnixTime()
	res, err := models.ActivateUser(c.QueryParam("activation_key"), unixtime)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	return c.JSON(http.StatusOK, affected)
}

func GetUser(c echo.Context) error {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	user, err := models.FindUserByID(userID)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, user)
}
