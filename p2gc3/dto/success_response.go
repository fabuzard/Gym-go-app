package dto

import "p2gc3/model"

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Weight   int    `json:"weight"`
	Height   int    `json:"height"`
}

type WorkoutResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      uint   `json:"user_id"`
}
type WorkoutByIDResponse struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	UserID      uint             `json:"user_id"`
	Exercises   []model.Exercise `json:"exercises"`
}

type WorkoutPutDeleteResponse struct {
	Message string        `json:"message"`
	Data    model.Workout `json:"data"`
}

type WorkoutDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ExerciseResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ExerciseLogResponse struct {
	ExerciseID uint `json:"exercise_id"`
	SetCount   uint `json:"set_count"`
	RepCount   uint `json:"rep_count"`
	Weight     uint `json:"weight"`
}

type UserInfoWithBMIResponse struct {
	ID             uint    `json:"id"`
	Email          string  `json:"email"`
	FullName       string  `json:"full_name"`
	Weight         int     `json:"weight"`
	Height         int     `json:"height"`
	BMI            float64 `json:"bmi"`
	WeightCategory string  `json:"weight_category"`
}
