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

// CreateWorkout godoc
// @Summary      Create a new workout
// @Description  Creates a workout with name and description
// @Tags         workouts
// @Accept       json
// @Produce      json
// @Param        workout  body  dto.CreateOrUpdateWorkoutRequest  true  "Workout data"
// @Success      201  {object}  dto.SuccessResponse{data=dto.WorkoutResponse}
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Router       /api/workouts [post]
// @Security     BearerAuth
func CreateWorkout(c echo.Context) error {
	var req dto.WorkoutCreateRequest

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
			Message: "workout name and description are required",
		})
	}

	Workout := model.Workout{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := config.DB.Create(&Workout).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Error creating data",
			Details: err.Error(),
		})
	}

	response := dto.WorkoutResponse{
		ID:          Workout.ID,
		Name:        Workout.Name,
		Description: Workout.Description,
		UserID:      Workout.UserID,
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse{
		Message: "workout successfully created",
		Data:    response,
	})

}

// GetWorkout godoc
// @Summary      Get all workouts
// @Description  Retrieves all workouts for the authenticated user
// @Tags         workouts
// @Produce      json
// @Success      200  {object}  dto.SuccessResponse{data=[]dto.WorkoutResponse}
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/workouts [get]
// @Security     BearerAuth
func GetWorkout(c echo.Context) error {
	userID, err := helper.ExtractUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Failed to extract user information",
			Details: err.Error(),
		})
	}

	var w []model.Workout

	if err := config.DB.
		Where("user_id = ?", userID).
		Find(&w).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to retrieve workouts",
			Details: err.Error(),
		})
	}

	var response []dto.WorkoutResponse

	for _, wo := range w {
		response = append(response, dto.WorkoutResponse{
			ID:          wo.ID,
			Name:        wo.Name,
			Description: wo.Description,
			UserID:      wo.UserID,
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Success Retrieving Workouts",
		Data:    response,
	})
}

// GetWorkoutByID godoc
// @Summary      Get workout by ID
// @Description  Retrieves a specific workout with exercises
// @Tags         workouts
// @Produce      json
// @Param        id  path  int  true  "Workout ID"
// @Success      200  {object}  dto.WorkoutByIDResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      403  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Router       /api/workouts/{id} [get]
// @Security     BearerAuth
func GetWorkoutByID(c echo.Context) error {
	userID, err := helper.ExtractUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Failed to extract user information",
			Details: err.Error(),
		})
	}

	workoutIDParam := c.Param("id")
	workoutID, err := strconv.Atoi(workoutIDParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid workout ID",
			Details: err.Error(),
		})
	}

	var w model.Workout
	err = config.DB.Preload("Exercises").First(&w, workoutID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, dto.ErrorResponse{
			Message: "Workout not found",
			Details: err.Error(),
		})
	}

	if w.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, dto.ErrorResponse{
			Message: "You are not authorized to view this workout",
		})
	}

	// Success
	return c.JSON(http.StatusOK, dto.WorkoutByIDResponse{
		ID:          w.ID,
		Name:        w.Name,
		Description: w.Description,
		UserID:      w.UserID,
		Exercises:   w.Exercises,
	})
}

// UpdateWorkout godoc
// @Summary      Update a workout
// @Description  Updates workout name and description
// @Tags         workouts
// @Accept       json
// @Produce      json
// @Param        id       path  int                                 true  "Workout ID"
// @Param        workout  body  dto.CreateOrUpdateWorkoutRequest    true  "Updated workout data"
// @Success      200  {object}  dto.SuccessResponse{data=dto.WorkoutDTO}
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      403  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Router       /api/workouts/{id} [put]
// @Security     BearerAuth
func UpdateWorkout(c echo.Context) error {
	userID, err := helper.ExtractUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Failed to extract user information",
			Details: err.Error(),
		})
	}

	workoutID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid workout ID",
			Details: err.Error(),
		})
	}

	var workout model.Workout
	err = config.DB.First(&workout, workoutID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, dto.ErrorResponse{
			Message: "Workout not found",
		})
	}
	if workout.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, dto.ErrorResponse{
			Message: "You are not authorized to update this workout",
		})
	}

	var req dto.CreateOrUpdateWorkoutRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid request body",
			Details: err.Error(),
		})
	}

	workout.Name = req.Name
	workout.Description = req.Description

	if err := config.DB.Save(&workout).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to update workout",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Workout updated successfully",
		Data: dto.WorkoutDTO{
			ID:          workout.ID,
			Name:        req.Name,
			Description: req.Description,
		},
	})
}

// DeleteWorkout godoc
// @Summary      Delete a workout
// @Description  Deletes a workout and all associated exercises and logs
// @Tags         workouts
// @Produce      json
// @Param        id  path  int  true  "Workout ID"
// @Success      200  {object}  dto.SuccessResponse{data=dto.WorkoutDTO}
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      403  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/workouts/{id} [delete]
// @Security     BearerAuth
func DeleteWorkout(c echo.Context) error {
	userID, err := helper.ExtractUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Unauthorized",
			Details: err.Error(),
		})
	}

	workoutID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid workout ID",
			Details: err.Error(),
		})
	}

	var workout model.Workout
	err = config.DB.Preload("Exercises").First(&workout, workoutID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return echo.NewHTTPError(http.StatusNotFound, dto.ErrorResponse{
			Message: "Workout not found",
		})
	}
	if workout.UserID != userID {
		return echo.NewHTTPError(http.StatusForbidden, dto.ErrorResponse{
			Message: "You are not authorized to delete this workout",
		})
	}

	tx := config.DB.Begin()

	for _, ex := range workout.Exercises {
		if err := tx.Where("exercise_id = ?", ex.ID).Delete(&model.ExerciseLog{}).Error; err != nil {
			tx.Rollback()
			return echo.NewHTTPError(http.StatusInternalServerError, dto.ErrorResponse{
				Message: "Failed to delete exercise logs",
				Details: err.Error(),
			})
		}
	}

	if err := tx.Where("workout_id = ?", workout.ID).Delete(&model.Exercise{}).Error; err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to delete exercises",
			Details: err.Error(),
		})
	}

	if err := tx.Delete(&workout).Error; err != nil {
		tx.Rollback()
		return echo.NewHTTPError(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to delete workout",
			Details: err.Error(),
		})
	}

	tx.Commit()

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Workout deleted successfully",
		Data: dto.WorkoutDTO{
			ID:          workout.ID,
			Name:        workout.Name,
			Description: workout.Description,
		},
	})
}
