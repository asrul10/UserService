package helper

type HelperInterface interface {
	HashPassword(password string) (string, error)
	ComparePassword(password string, hashedPassword string) error
	GenerateAccessToken(token *string, id int) error
	GenerateRefreshToken(token *string, id int) error
	VerifyToken(tokenString string) (string, error)
	GetToken(authorization string) string
}
