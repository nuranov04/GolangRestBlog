package user

import (
	"context"
	"go.mod/pkg/logging"
)

type Service struct {
	storage Storage
	logger  *logging.Logger
}

func (s Service) Create(ctx context.Context, user User) error {
	err := s.storage.Create(ctx, user)
	if err != nil {
		s.logger.Fatalf("failed due create user in Service.Create. err: %e", err)
		return err
	}
	return nil
}

func (s Service) FindAll(ctx context.Context) ([]User, error) {
	all, err := s.storage.FindAll(ctx)
	if err != nil {
		s.logger.Trace(err)
		return nil, err
	}
	return all, nil
}

func (s Service) FindOneById(ctx context.Context, id int) (User, error) {
	obj, err := s.storage.FindOneById(ctx, id)
	if err != nil {
		return User{}, err
	}
	return obj, nil

}

func NewService(storage Storage, logger *logging.Logger) *Service {
	return &Service{
		storage: storage,
		logger:  logger,
	}
}
