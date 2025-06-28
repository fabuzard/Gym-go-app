// @title Fitness Tracker API
// @version 1.0
// @description API for managing users, workouts, exercises, and logs.
// @host localhost:8080
// @BasePath /

package main

import (
	"fmt"
	"os"
	"p2gc3/config"
	"p2gc3/handler"
	"p2gc3/middleware"
	"p2gc3/model"
	"p2gc3/routes"

	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	// setup configs, loading env, and initializing db
	config.LoadEnv()
	db := config.DBInit()

	// Auto migrating into DB
	err := db.AutoMigrate(&model.User{}, &model.Workout{}, &model.Exercise{}, &model.ExerciseLog{})
	if err != nil {
		panic("Failed to auto migrate: " + err.Error())
	}

	// initialize echo
	e := echo.New()

	// Using middleware
	e.Use(middleware.MiddlewareLogging)

	e.HTTPErrorHandler = handler.CustomHTTPErrorHandler

	// Testing endpoint
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "🚀 Server running and DB connected!")
	})

	// Auth Route
	routes.AllRoutes(e)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Connected to db")
	e.Logger.Fatal(e.Start(":" + port))
}
