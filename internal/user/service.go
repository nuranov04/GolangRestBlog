package user

import (
	"context"
	"go.mod/pkg/logging"
)

type Service interface {
	Create(ctx context.Context, user User) error
	Delete(ctx context.Context, userId int) error
	userUpdate(ctx context.Context, updateUser User) (User, error)
	FindAll(ctx context.Context) ([]User, error)
	FindOneById(ctx context.Context, id int) (User, error)
	FindOneByUsername(ctx context.Context, username string) (User, error)
	FindOneByEmail(ctx context.Context, email string) (User, error)
}

type service struct {
	storage Storage
	logger  *logging.Logger
}

func (s service) userUpdate(ctx context.Context, userObj User, updateUser UpdateUserDTO) (UpdateUserDTO, error) {
	UpdateUser, err := s.storage.Update(ctx, userObj, updateUser)
	if err != nil {
		return UpdateUser, err
	}
	return UpdateUser, nil
}

func (s service) Create(ctx context.Context, user User) error {
	err := s.storage.Create(ctx, user)
	if err != nil {
		return err
	}
	return nil
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
		s.logger.Trace(err)
		return nil, err
	}
	return all, nil
}

func (s service) FindOneById(ctx context.Context, id int) (User, error) {
	obj, err := s.storage.FindOneById(ctx, id)
	if err != nil {
		return User{}, err
	}
	return obj, nil

}

func (s service) FindOneByUsername(ctx context.Context, username string) (User, error) {
	obj, err := s.storage.FindOneByUsername(ctx, username)
	if err != nil {
		return User{}, err
	}
	return obj, nil

}

func (s service) FindOneByEmail(ctx context.Context, email string) (User, error) {
	obj, err := s.storage.FindOneByUsername(ctx, email)
	if err != nil {
		return User{}, err
	}
	return obj, nil

}

func NewService(storage Storage, logger *logging.Logger) *service {
	return &service{
		storage: storage,
		logger:  logger,
	}
}
