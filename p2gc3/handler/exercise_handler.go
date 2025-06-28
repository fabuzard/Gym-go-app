package handler

import (
	"errors"
	"net/http"
	"p2gc3/config"
	"p2gc3/dto"
	helper "p2gc3/helpers"
	"p2gc3/model"
	"strconv"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// CreateExercise godoc
// @Summary      Create a new exercise
// @Description  Adds a new exercise to the database
// @Tags         exercises
// @Accept       json
// @Produce      json
// @Param        exercise  body  dto.ExerciseCreateRequest  true  "Exercise payload"
// @Success      201  {object}  dto.SuccessResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      403  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /exercises [post]
// @Security     BearerAuth
func CreateExercise(c echo.Context) error {
	var req dto.ExerciseCreateRequest

	userID, err := helper.ExtractUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Failed to extract user information",
			Details: err.Error(),
		})
	}

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid input",
			Details: err.Error(),
		})
	}

	if req.Description == "" || req.Name == "" {
		return echo.NewHTTPError(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Exercise name and description are required",
		})
	}

	var workout model.Workout

	if err := config.DB.First(&workout, req.WorkoutID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, dto.ErrorResponse{
				Message: "Workout not found",
			})
		}
		return echo.NewHTTPError(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to fetch workout",
			Details: err.Error(),
		})
	}

	if workout.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, dto.ErrorResponse{
			Message: "You are not authorized to create exercise to this workout",
		})
	}

	exercise := model.Exercise{
		WorkoutID:   req.WorkoutID,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := config.DB.Create(&exercise).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to create exercise",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse{
		Message: "Exercise successfully created",
		Data: dto.ExerciseResponse{
			Name:        exercise.Name,
			Description: exercise.Description,
		},
	})

}

// DeleteExercise godoc
// @Summary      Delete an exercise
// @Description  Deletes an exercise by ID
// @Tags         exercises
// @Produce      json
// @Param        id  path  int  true  "Exercise ID"
// @Success      200  {object}  dto.SuccessResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      403  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /exercises/{id} [delete]
// @Security     BearerAuth
func DeleteExercise(c echo.Context) error {
	userID, err := helper.ExtractUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized",
			Details: err.Error(),
		})
	}

	exerciseIDParam := c.Param("id")
	exerciseID, err := strconv.Atoi(exerciseIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid exercise ID",
			Details: err.Error(),
		})
	}

	var exercise model.Exercise
	err = config.DB.Preload("Workout").First(&exercise, exerciseID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, dto.ErrorResponse{
			Message: "Exercise not found",
		})
	} else if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to retrieve exercise",
			Details: err.Error(),
		})
	}

	if exercise.Workout.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, dto.ErrorResponse{
			Message: "You are not authorized to delete this exercise",
		})
	}

	// Delete associated logs
	if err := config.DB.Where("exercise_id = ?", exercise.ID).Delete(&model.ExerciseLog{}).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to delete associated exercise logs",
			Details: err.Error(),
		})
	}

	// Delete exercise
	if err := config.DB.Delete(&exercise).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to delete exercise",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Exercise deleted successfully",
		Data: dto.ExerciseResponse{
			Name:        exercise.Name,
			Description: exercise.Description,
		},
	})
}
