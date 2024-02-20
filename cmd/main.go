package main

import (
	"os"

	"github.com/asrul10/UserService/generated"
	"github.com/asrul10/UserService/handler"
	"github.com/asrul10/UserService/helper"
	"github.com/asrul10/UserService/repository"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	var server generated.ServerInterface = newServer(e)

	generated.RegisterHandlers(e, server)
	e.Logger.Fatal(e.Start(":1323"))
}

func newServer(e *echo.Echo) *handler.Server {
	dbDsn := os.Getenv("DATABASE_URL")
	jwtPrivateKeyPath := os.Getenv("JWT_PRIVATE_KEY_PATH")
	jwtPublicKeyPath := os.Getenv("JWT_PUBLIC_KEY_PATH")

	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})
	var helper helper.HelperInterface = helper.NewHelper(helper.NewHelperOptions{
		JwtPrivateKeyPath: jwtPrivateKeyPath,
		JwtPublicKeyPath:  jwtPublicKeyPath,
	})

	opts := handler.NewServerOptions{
		Repository: repo,
		Helper:     helper,
		Echo:       e,
	}
	return handler.NewServer(opts)
}
