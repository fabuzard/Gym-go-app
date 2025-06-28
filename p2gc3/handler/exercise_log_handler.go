package handler

import (
	"net/http"
	"p2gc3/config"
	"p2gc3/dto"
	helper "p2gc3/helpers"
	"p2gc3/model"
	"strconv"

	"github.com/labstack/echo/v4"
)

// CreateExerciseLog godoc
// @Summary      Create an exercise log
// @Description  Logs an exercise with set, rep, and weight
// @Tags         exercise-logs
// @Accept       json
// @Produce      json
// @Param        log  body  dto.ExerciseLogRequest  true  "Exercise log data"
// @Success      201  {object}  dto.SuccessResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      403  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /exercise-logs [post]
// @Security     BearerAuth
func CreateExerciseLog(c echo.Context) error {
	userID, err := helper.ExtractUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized",
			Details: err.Error(),
		})
	}

	var req dto.ExerciseLogRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid input",
			Details: err.Error(),
		})
	}

	// Check if the exercise exists and belongs to this user
	var exercise model.Exercise
	if err := config.DB.Preload("Workout").First(&exercise, req.ExerciseID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, dto.ErrorResponse{
			Message: "Exercise not found",
			Details: err.Error(),
		})
	}

	if exercise.Workout.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, dto.ErrorResponse{
			Message: "You are not authorized to log this exercise",
		})
	}

	log := model.ExerciseLog{
		ExerciseID: req.ExerciseID,
		UserID:     userID,
		SetCount:   req.SetCount,
		RepCount:   req.RepCount,
		Weight:     req.Weight,
		CreatedAt:  req.CreatedAt,
	}

	if err := config.DB.Create(&log).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to create log",
			Details: err.Error(),
		})
	}

	response := dto.ExerciseLogResponse{
		ExerciseID: log.ExerciseID,
		SetCount:   uint(log.SetCount),
		RepCount:   uint(log.RepCount),
		Weight:     uint(log.Weight),
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse{
		Message: "Log created successfully",
		Data:    response,
	})
}

// DeleteExerciseLog godoc
// @Summary      Delete an exercise log
// @Description  Deletes an exercise log by its ID
// @Tags         exercise-logs
// @Produce      json
// @Param        id  path  int  true  "Log ID"
// @Success      200  {object}  dto.SuccessResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      403  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /exercise-logs/{id} [delete]
// @Security     BearerAuth
func DeleteExerciseLog(c echo.Context) error {
	userID, err := helper.ExtractUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized",
			Details: err.Error(),
		})
	}

	logID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid log ID",
			Details: err.Error(),
		})
	}

	var log model.ExerciseLog
	if err := config.DB.First(&log, logID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, dto.ErrorResponse{
			Message: "Log not found",
		})
	}

	if log.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, dto.ErrorResponse{
			Message: "You are not authorized to delete this log",
		})
	}

	if err := config.DB.Delete(&log).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to delete log",
			Details: err.Error(),
		})
	}

	response := dto.ExerciseLogResponse{
		ExerciseID: log.ExerciseID,
		SetCount:   uint(log.SetCount),
		RepCount:   uint(log.RepCount),
		Weight:     uint(log.Weight),
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Log deleted successfully",
		Data:    response,
	})
}
