package auth

import (
	"os"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret []byte

func init() {
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if jwtSecret == nil {
		s := uuid.NewString()
		jwtSecret = []byte(s)
	}
}

func ValidatePassword(hash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func CreatePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CreateSessionCookie() (string, error) {
	token, err := uuid.NewRandom()

	if err != nil {
		return "", err
	}

	return token.String(), nil
}

func CreateToken() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	signed, err := token.SignedString(jwtSecret)

	if err != nil {
		return "", err
	}

	return signed, nil
}

func ValidateToken(token string) error {
	return validateTokenInternal(token, jwtSecret)
}

func validateTokenInternal(token string, secret []byte) error {
	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return err
	}

	return nil
}
