package helper

import (
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	AccessTokenExpireDuration  = time.Hour * 24
	RefreshTokenExpireDuration = time.Hour * 24 * 7
)

func (h *Helper) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (h *Helper) ComparePassword(password string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
	return err
}

func (h *Helper) getPrivateKey() (*rsa.PrivateKey, error) {
	read, err := os.ReadFile(h.JwtPrivateKeyPath)
	if err != nil {
		return nil, err
	}
	priv, err := jwt.ParseRSAPrivateKeyFromPEM(read)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

func (h *Helper) getPulicKey() (*rsa.PublicKey, error) {
	read, err := os.ReadFile(h.JwtPublicKeyPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	pub, err := jwt.ParseRSAPublicKeyFromPEM(read)
	if err != nil {
		return nil, err
	}

	return pub, nil
}

func (h *Helper) GenerateAccessToken(token *string, id int) error {
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": id,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(AccessTokenExpireDuration).Unix(),
	})
	privateKey, err := h.getPrivateKey()
	if err != nil {
		return err
	}
	tokenString, err := t.SignedString(privateKey)
	if err != nil {
		return err
	}
	*token = tokenString
	return nil
}

func (h *Helper) GenerateRefreshToken(refreshToken *string, id int) error {
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(RefreshTokenExpireDuration).Unix(),
	})
	privateKey, err := h.getPrivateKey()
	if err != nil {
		return err
	}
	tokenString, err := t.SignedString(privateKey)
	if err != nil {
		return err
	}
	*refreshToken = tokenString
	return nil
}

func (h *Helper) VerifyToken(tokenString string) (string, error) {
	pubKey, err := h.getPulicKey()
	if err != nil {
		return "", err
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method == jwt.SigningMethodES256 && token.Valid {
			return nil, fmt.Errorf("Unexpected signing method: %v, token not valid", token.Header["alg"])
		}
		return pubKey, nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", err
	}

	userId, ok := (claims)["sub"]

	if !ok {
		return "", fmt.Errorf("Invalid token")
	}

	return fmt.Sprintf("%v", userId), nil
}

func (h *Helper) GetToken(authorization string) string {
	token := ""
	if len(authorization) > 7 && authorization[:7] == "Bearer " {
		token = authorization[7:]
	}
	return token
}
