package user

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id" bson:"_id,omitempty"`
	Username string `json:"username" bson:"username"`
	Password string `json:"-" bson:"password"`
	Email    string `json:"email" bson:"email"`
}

func (u *User) GeneratePasswordHash() error {
	passwordHash, err := GeneratePasswordHash(u.Password)
	if err != nil {
		return err
	}
	u.Password = passwordHash
	return nil
}

func GeneratePasswordHash(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("error to password hash due error %w", err)
	}
	return string(passwordHash), nil
}

func NewUser(dto CreateUserDTO) User {
	return User{
		Username: dto.Username,
		Password: dto.Password,
		Email:    dto.Email,
	}
}
