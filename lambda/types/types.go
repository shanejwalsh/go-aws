package types

import (
	"golang.org/x/crypto/bcrypt"
)

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
