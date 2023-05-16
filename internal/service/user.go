package service

import (
	"context"
	"fmt"
	"go.mod/internal/apperror"
	"go.mod/internal/model"
	"go.mod/internal/storage"
	"go.mod/pkg/logging"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Create(ctx context.Context, createUser model.CreateUserDTO) (u *model.User, err error)
	Delete(ctx context.Context, userId int) error
	UserUpdate(ctx context.Context, userObj model.User, updateUser model.UpdateUserDTO) (u *model.User, err error)
	FindAll(ctx context.Context) ([]model.User, error)
	FindUserByUsernameAndPassword(ctx context.Context, dto model.LoginDTO) (u model.User, err error)
	FindOneById(ctx context.Context, id int) (u *model.User, err error)
	FindOneByUsername(ctx context.Context, username string) (u *model.User, err error)
	FindOneByEmail(ctx context.Context, email string) (u *model.User, err error)
}

type userService struct {
	storage storage.UserStorage
	logger  *logging.Logger
}

func NewUserService(storage storage.UserStorage, logger *logging.Logger) *userService {
	return &userService{
		storage: storage,
		logger:  logger,
	}
}

func (s userService) FindUserByUsernameAndPassword(ctx context.Context, loginDTO model.LoginDTO) (u model.User, err error) {
	userObjLink, err := s.FindOneByUsername(ctx, loginDTO.Username)
	userObj := *userObjLink
	if err != nil {
		return u, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(userObj.Password), []byte(loginDTO.Password)); err != nil {
		return u, apperror.ErrorNotFound
	}
	return userObj, nil
}

func (s userService) UserUpdate(ctx context.Context, userObj model.User, updateUser model.UpdateUserDTO) (u *model.User, err error) {
	UpdateUser, err := s.storage.Update(ctx, userObj, updateUser)
	if err != nil {
		return nil, err
	}
	return UpdateUser, nil
}

func (s userService) Create(ctx context.Context, createUser model.CreateUserDTO) (u *model.User, err error) {
	if createUser.Password != createUser.RepeatPassword {
		return u, apperror.BadRequestError("password does not match repeat password")
	}
	user := model.NewUser(createUser)
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

func (s userService) FindAll(ctx context.Context) ([]model.User, error) {
	all, err := s.storage.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return all, nil
}

func (s userService) FindOneById(ctx context.Context, id int) (u *model.User, err error) {
	obj, err := s.storage.FindOneById(ctx, id)
	if err != nil {
		return nil, err
	}
	return obj, nil

}

func (s userService) FindOneByUsername(ctx context.Context, username string) (u *model.User, err error) {
	obj, err := s.storage.FindOneByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return obj, nil

}

func (s userService) FindOneByEmail(ctx context.Context, email string) (u *model.User, err error) {
	obj, err := s.storage.FindOneByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return obj, nil

}
