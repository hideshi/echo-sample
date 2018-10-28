package main

import (
	"github.com/BurntSushi/toml"
	"github.com/hideshi/echo-sample/controllers"
	"github.com/hideshi/echo-sample/models"
	"github.com/hideshi/echo-sample/structs"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Create app
	e := echo.New()

	// Create tables
	if _, err := models.InitDB(); err != nil {
		e.Logger.Fatal(err)
	}

	// Load config
	if _, err := toml.DecodeFile("config/config.toml", &structs.Conf); err != nil {
		e.Logger.Fatal(err)
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/users/:id", controllers.GetUser)
	e.POST("/users", controllers.CreateUser)
	e.GET("/users/activate", controllers.ActivateUser)
	e.PATCH("/users/:id", controllers.UpdateEmail)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
