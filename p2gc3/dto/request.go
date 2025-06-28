package dto

import "time"

// 📥 For user registration
type UserRegisterRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	FullName string `json:"full_name" form:"full_name" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
	Weight   int    `json:"weight" form:"weight" validate:"required"`
	Height   int    `json:"height" form:"height" validate:"required"`
}

// 📥 For creating a workout
type WorkoutCreateRequest struct {
	Name        string `json:"name" form:"name" validate:"required"`
	Description string `json:"description" form:"description" validate:"required"`
	UserID      uint   `json:"user_id" form:"user_id" validate:"required"` // optional if from JWT
}

// 📥 For adding an exercise to a workout
type ExerciseCreateRequest struct {
	Name        string `json:"name" form:"name" validate:"required"`
	Description string `json:"description" form:"description" validate:"required"`
	WorkoutID   uint   `json:"workout_id" form:"workout_id" validate:"required"`
}

// 📥 For logging an exercise session
type ExerciseLogCreateRequest struct {
	ExerciseID uint `json:"exercise_id" form:"exercise_id" validate:"required"`
	UserID     uint `json:"user_id" form:"user_id" validate:"required"`
	SetCount   int  `json:"set_count" form:"set_count" validate:"required"`
	RepCount   int  `json:"rep_count" form:"rep_count" validate:"required"`
	Weight     int  `json:"weight" form:"weight" validate:"required"`
}

// Login
type UserLoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}
type CreateOrUpdateWorkoutRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ExerciseLogRequest struct {
	ExerciseID uint      `json:"exercise_id"`
	SetCount   int       `json:"set_count"`
	RepCount   int       `json:"rep_count"`
	Weight     int       `json:"weight"`
	CreatedAt  time.Time `json:"created_at"`
}
