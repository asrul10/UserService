package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/asrul10/UserService/generated"
	"github.com/asrul10/UserService/repository"
	"github.com/labstack/echo/v4"
)

// (POST /users)
func (s *Server) RegisterUser(ctx echo.Context) error {
	user := new(generated.RegisterUserJSONRequestBody)

	if err := ctx.Bind(user); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Invalid request body",
		})
	}

	// Validate request body
	if err := ctx.Validate(user); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Check if phone number already registered
	if _, err := s.Repository.GetUserByPhoneNumber(context.Background(), repository.GetUserByPhoneNumberInput{
		PhoneNumber: user.PhoneNumber,
	}); err == nil {
		return ctx.JSON(http.StatusConflict, generated.ErrorResponse{
			Message: "Phone number already registered",
		})
	}

	// Create user
	hashPassword, err := s.Helper.HashPassword(user.Password)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "Failed to hash password",
		})
	}
	resp, err := s.Repository.CreateUser(ctx.Request().Context(), repository.CreateUserInput{
		PhoneNumber: user.PhoneNumber,
		FullName:    user.FullName,
		Password:    hashPassword,
	})

	if err != nil {
		log.Println(err)
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, resp)
}

// (POST /users/login)
func (s *Server) LoginUser(ctx echo.Context) error {
	user := new(generated.LoginUserJSONRequestBody)

	if err := ctx.Bind(user); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Invalid request body",
		})
	}

	// Validate request body
	if err := ctx.Validate(user); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Get user by phone number
	resp, err := s.Repository.GetUserByPhoneNumber(ctx.Request().Context(), repository.GetUserByPhoneNumberInput{
		PhoneNumber: user.PhoneNumber,
	})
	if err != nil {
		return ctx.JSON(http.StatusNotFound, generated.ErrorResponse{
			Message: "User not found",
		})
	}

	// Check if password is correct
	if err := s.Helper.ComparePassword(user.Password, resp.Password); err != nil {
		return ctx.JSON(http.StatusUnauthorized, generated.ErrorResponse{
			Message: "Invalid password",
		})
	}

	// Generate token
	token := ""
	if err := s.Helper.GenerateAccessToken(&token, resp.UserId); err != nil {
		log.Println(err)
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "Failed to generate token",
		})
	}
	refreshToken := ""
	if err := s.Helper.GenerateRefreshToken(&refreshToken, resp.UserId); err != nil {
		log.Println(err)
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{
			Message: "Failed to generate refresh token",
		})
	}

	// Update success login
	_, err = s.Repository.SuccessLoginCount(ctx.Request().Context(), repository.SuccessLoginCountInput{
		UserId: resp.UserId,
	})
	if err != nil {
		log.Println(err)
	}

	return ctx.JSON(http.StatusOK, generated.LoginUserResponse{
		UserId:       resp.UserId,
		AccessToken:  token,
		RefreshToken: refreshToken,
	})
}

// (GET /users/{id})
func (s *Server) GetUser(ctx echo.Context) error {
	token := s.Helper.GetToken(ctx.Request().Header.Get("Authorization"))
	userIdStr, err := s.Helper.VerifyToken(token)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, generated.ErrorResponse{
			Message: "Unauthorized",
		})
	}

	// Get user by id
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Invalid user id",
		})
	}
	resp, err := s.Repository.GetUserById(ctx.Request().Context(), repository.GetUserByIdInput{
		UserId: userId,
	})
	if err != nil {
		return ctx.JSON(http.StatusNotFound, generated.ErrorResponse{
			Message: "User not found",
		})
	}

	return ctx.JSON(http.StatusOK, generated.GetUserResponse{
		UserId:      resp.UserId,
		PhoneNumber: resp.PhoneNumber,
		FullName:    resp.FullName,
	})
}

// (PUT /users/{id})
func (s *Server) UpdateUser(ctx echo.Context) error {
	token := s.Helper.GetToken(ctx.Request().Header.Get("Authorization"))
	userIdStr, err := s.Helper.VerifyToken(token)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, generated.ErrorResponse{
			Message: "Unauthorized",
		})
	}

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Invalid user id",
		})
	}

	user := new(generated.UpdateUserJSONRequestBody)
	if err := ctx.Bind(user); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: "Invalid request body",
		})
	}

	// Validate request body
	if err := ctx.Validate(user); err != nil {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{
			Message: err.Error(),
		})
	}

	// Prevent updating phone number
	isChanged, err := s.Repository.IsPhoneNumberChanged(ctx.Request().Context(), repository.IsPhoneNumberChangedInput{
		UserId:      userId,
		PhoneNumber: user.PhoneNumber,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	if isChanged.IsChanged {
		return ctx.JSON(http.StatusConflict, generated.ErrorResponse{
			Message: "Phone number cannot be changed",
		})
	}

	// Update user
	resp, err := s.Repository.UpdateUserById(ctx.Request().Context(), repository.UpdateUserByIdInput{
		UserId:      userId,
		PhoneNumber: user.PhoneNumber,
		FullName:    user.FullName,
	})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, resp)
}
