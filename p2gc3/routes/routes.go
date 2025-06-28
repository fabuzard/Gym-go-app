package routes

import (
	"os"
	"p2gc3/handler"
	"p2gc3/middleware"

	_ "p2gc3/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func AllRoutes(e *echo.Echo) {
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	userGroup := e.Group("/users")
	userGroup.POST("/register", handler.Register)
	userGroup.POST("/login", handler.Login)
	userGroup.GET("", handler.UserInfo, middleware.JWTMiddleware(os.Getenv("JWT_SECRET")))

	apiGroup := e.Group("/api", middleware.JWTMiddleware(os.Getenv("JWT_SECRET")))

	// Group for /api/workouts
	workoutGroup := apiGroup.Group("/workouts")
	workoutGroup.POST("", handler.CreateWorkout)
	workoutGroup.GET("", handler.GetWorkout)
	workoutGroup.GET("/:id", handler.GetWorkoutByID)
	workoutGroup.PUT("/:id", handler.UpdateWorkout)
	workoutGroup.DELETE("/:id", handler.DeleteWorkout)

	exerciseGroup := apiGroup.Group("/exercise", middleware.JWTMiddleware(os.Getenv("JWT_SECRET")))
	exerciseGroup.POST("", handler.CreateExercise)
	exerciseGroup.DELETE("/:id", handler.DeleteExercise)

	logGroup := apiGroup.Group("/logs")
	logGroup.POST("", handler.CreateExerciseLog)
	logGroup.DELETE("/:id", handler.DeleteExerciseLog)

}
