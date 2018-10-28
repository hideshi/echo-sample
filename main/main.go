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
	// Create tables
	models.InitDB()

	// Create app
	e := echo.New()

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

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
