package helper

import "github.com/labstack/echo/v4"

type Helper struct {
	JwtPrivateKeyPath string
	JwtPublicKeyPath  string
	echo              *echo.Echo
}

type NewHelperOptions struct {
	JwtPrivateKeyPath string
	JwtPublicKeyPath  string
	echo              *echo.Echo
}

func NewHelper(options NewHelperOptions) *Helper {
	return &Helper{
		JwtPrivateKeyPath: options.JwtPrivateKeyPath,
		JwtPublicKeyPath:  options.JwtPublicKeyPath,
		echo:              options.echo,
	}
}
