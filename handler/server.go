package handler

import (
	"github.com/asrul10/UserService/helper"
	"github.com/asrul10/UserService/repository"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Server struct {
	Repository repository.RepositoryInterface
	Helper     helper.HelperInterface
}

type NewServerOptions struct {
	Repository repository.RepositoryInterface
	Helper     helper.HelperInterface
	Echo       *echo.Echo
}

func NewServer(opts NewServerOptions) *Server {
	// Register custom validator
	opts.Echo.Validator = &CustomValidator{
		validator: validator.New(),
	}

	return &Server{
		Repository: opts.Repository,
		Helper:     opts.Helper,
	}
}
