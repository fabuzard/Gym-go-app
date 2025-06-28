package handler

import (
	"p2gc3/config"
	"p2gc3/dto"
	helper "p2gc3/helpers"
	"p2gc3/model"
	"p2gc3/utils"

	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Register godoc
// @Summary      Register a new user
// @Description  Creates a new user account
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body  dto.UserRegisterRequest  true  "Register payload"
// @Success      201  {object}  dto.SuccessResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /users/register [post]
func Register(c echo.Context) error {
	var req dto.UserRegisterRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid input",
			Details: err.Error(),
		})
	}

	if req.Email == "" || req.FullName == "" || req.Password == "" || req.Weight <= 0 || req.Height <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid input: required fields must be filled",
			Details: "One or more required fields are missing or invalid",
		})
	}

	var existingUser model.User
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return echo.NewHTTPError(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Email already registered",
		})
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to hash password",
			Details: err.Error(),
		})
	}

	User := model.User{
		Email:    req.Email,
		FullName: req.FullName,
		Password: string(bytes),
		Weight:   req.Weight,
		Height:   req.Height,
	}

	if err := config.DB.Create(&User).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Error creating data",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.SuccessResponse{
		Message: "Register sukses",
		Data:    User,
	})
}

// Login godoc
// @Summary      Login user
// @Description  Authenticates a user and returns JWT
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body  dto.UserLoginRequest  true  "Login credentials"
// @Success      200  {object} dto.SuccessResponse
// @Failure      400  {object} dto.ErrorResponse
// @Router       /users/login [post]
func Login(c echo.Context) error {
	var req dto.UserLoginRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Invalid input",
			Details: err.Error(),
		})
	}

	if req.Email == "" || req.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, dto.ErrorResponse{
			Message: "Email and password are required",
		})
	}

	var user model.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Invalid credentials",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Invalid credentials",
		})
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, dto.ErrorResponse{
			Message: "Failed to generate token",
			Details: err.Error(),
		})
	}
	user.Password = ""

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Login successful",
		Data: map[string]interface{}{
			"user":      user,
			"jwt_token": token,
		},
	})
}

// UserInfo godoc
// @Summary      Get authenticated user
// @Description  Returns the authenticated user's data, BMI, and weight category
// @Tags         users
// @Produce      json
// @Success      200  {object}  dto.SuccessResponse
// @Failure      401  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Router       /users [get]
// @Security     BearerAuth
func UserInfo(c echo.Context) error {
	userID, err := helper.ExtractUserID(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, dto.ErrorResponse{
			Message: "Failed to extract user information",
			Details: err.Error(),
		})
	}

	var u model.User
	if err := config.DB.First(&u, userID).Error; err != nil {
		return echo.NewHTTPError(http.StatusNotFound, dto.ErrorResponse{
			Message: "User not found",
			Details: err.Error(),
		})
	}

	// 3rd party BMI usage
	bmi, category, err := helper.GetBMIAndCategory(u.Weight, u.Height)
	if err != nil {
		category = "Unknown"
	}

	return c.JSON(http.StatusOK, dto.SuccessResponse{
		Message: "Authenticated user retrieved",
		Data: dto.UserInfoWithBMIResponse{
			ID:             u.ID,
			Email:          u.Email,
			FullName:       u.FullName,
			Weight:         u.Weight,
			Height:         u.Height,
			BMI:            bmi,
			WeightCategory: category,
		},
	})
}
