package user

import (
	"context"
	"fmt"
	"go.mod/internal/apperror"
	"go.mod/pkg/logging"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Create(ctx context.Context, createUser CreateUserDTO) (u *User, err error)
	Delete(ctx context.Context, userId int) error
	UserUpdate(ctx context.Context, userObj User, updateUser UpdateUserDTO) (u *User, err error)
	FindAll(ctx context.Context) ([]User, error)
	FindUserByUsernameAndPassword(ctx context.Context, username, password string) (u User, err error)
	FindOneById(ctx context.Context, id int) (u *User, err error)
	FindOneByUsername(ctx context.Context, username string) (u *User, err error)
	FindOneByEmail(ctx context.Context, email string) (u *User, err error)
}

type userService struct {
	storage Storage
	logger  *logging.Logger
}

func NewUserService(storage Storage, logger *logging.Logger) Service {
	return &userService{
		storage: storage,
		logger:  logger,
	}
}

func (s userService) FindUserByUsernameAndPassword(ctx context.Context, username, password string) (u User, err error) {
	userObjLink, err := s.FindOneByUsername(ctx, username)
	userObj := userObjLink
	if err != nil {
		return u, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userObjLink.Password), []byte(password))
	if err != nil {
		return u, apperror.NotCorrectPassword
	}
	return *userObj, nil
}

func (s userService) UserUpdate(ctx context.Context, userObj User, updateUser UpdateUserDTO) (u *User, err error) {
	UpdateUser, err := s.storage.Update(ctx, userObj, updateUser)
	if err != nil {
		return nil, err
	}
	return UpdateUser, nil
}

func (s userService) Create(ctx context.Context, createUser CreateUserDTO) (u *User, err error) {
	if createUser.Password != createUser.RepeatPassword {
		return u, apperror.BadRequestError("password does not match repeat password")
	}
	user := NewUser(createUser)
	err = user.GeneratePasswordHash()
	if err != nil {
		return nil, fmt.Errorf("failed to create user due to error %v", err)
	}
	u, err = s.storage.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s userService) Delete(ctx context.Context, userId int) error {
	err := s.storage.Delete(ctx, userId)
	if err != nil {
		return err
	}
	return nil
}

func (s userService) FindAll(ctx context.Context) ([]User, error) {
	all, err := s.storage.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return all, nil
}

func (s userService) FindOneById(ctx context.Context, id int) (u *User, err error) {
	obj, err := s.storage.FindOneById(ctx, id)
	if err != nil {
		return nil, err
	}
	return obj, nil

}

func (s userService) FindOneByUsername(ctx context.Context, username string) (u *User, err error) {
	obj, err := s.storage.FindOneByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return obj, nil

}

func (s userService) FindOneByEmail(ctx context.Context, email string) (u *User, err error) {
	obj, err := s.storage.FindOneByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return obj, nil

}
