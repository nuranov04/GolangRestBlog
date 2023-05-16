package user

import (
	"context"
	"fmt"
	"go.mod/internal/apperror"
	"go.mod/pkg/logging"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Create(ctx context.Context, user User) (u *User, err error)
	Delete(ctx context.Context, userId int) error
	userUpdate(ctx context.Context, userObj User, updateUser UpdateUserDTO) (u *UpdateUserDTO, err error)
	FindAll(ctx context.Context) ([]User, error)
	FindUserByUsernameAndPassword(ctx context.Context, username, email string) (u User, err error)
	FindOneById(ctx context.Context, id int) (u *User, err error)
	FindOneByUsername(ctx context.Context, username string) (u *User, err error)
	FindOneByEmail(ctx context.Context, email string) (u *User, err error)
}

type service struct {
	storage Storage
	logger  *logging.Logger
}

func NewService(storage Storage, logger *logging.Logger) *service {
	return &service{
		storage: storage,
		logger:  logger,
	}
}

func (s service) FindUserByUsernameAndPassword(ctx context.Context, loginDTO LoginDTO) (u User, err error) {
	userObjLink, err := s.FindOneByUsername(ctx, loginDTO.username)
	userObj := *userObjLink
	if err != nil {
		return u, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(userObj.Password), []byte(userObj.Password)); err != nil {
		return u, apperror.ErrorNotFound
	}
	return userObj, nil
}

func (s service) userUpdate(ctx context.Context, userObj User, updateUser UpdateUserDTO) (u *User, err error) {
	UpdateUser, err := s.storage.Update(ctx, userObj, updateUser)
	if err != nil {
		return nil, err
	}
	return UpdateUser, nil
}

func (s service) Create(ctx context.Context, createUser CreateUserDTO) (u *User, err error) {
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

func (s service) Delete(ctx context.Context, userId int) error {
	err := s.storage.Delete(ctx, userId)
	if err != nil {
		return err
	}
	return nil
}

func (s service) FindAll(ctx context.Context) ([]User, error) {
	all, err := s.storage.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return all, nil
}

func (s service) FindOneById(ctx context.Context, id int) (u *User, err error) {
	obj, err := s.storage.FindOneById(ctx, id)
	if err != nil {
		return nil, err
	}
	return obj, nil

}

func (s service) FindOneByUsername(ctx context.Context, username string) (u *User, err error) {
	obj, err := s.storage.FindOneByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return obj, nil

}

func (s service) FindOneByEmail(ctx context.Context, email string) (u *User, err error) {
	obj, err := s.storage.FindOneByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return obj, nil

}
