package types

import (
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	SECRET = "MY_SECRET"
)

type Req events.APIGatewayProxyRequest
type Res events.APIGatewayProxyResponse
type Next func(Req) (Res, error)

type RegisterUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}

func NewUser(registerUser RegisterUser) (User, error) {

	byteSlice := []byte(registerUser.Password)
	hashedPassWord, err := bcrypt.GenerateFromPassword(byteSlice, 10)

	if err != nil {
		return User{}, err

	}
	return User{
		Username:     registerUser.Username,
		PasswordHash: string(hashedPassWord),
	}, err

}

func ValidatePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func CreateToken(user User) string {
	now := time.Now()

	validUntil := now.Add(time.Minute * 30).Unix()

	claims := &jwt.MapClaims{
		"user":    user.Username,
		"expires": validUntil,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	tokenString, err := token.SignedString([]byte(SECRET))

	if err != nil {
		return fmt.Sprintf("there was an error %s", err)
	}

	return tokenString

}
