package user

import (
	"context"
	"go.mod/pkg/logging"
)

type Service struct {
	storage Storage
	logger  *logging.Logger
}

func (s Service) Create(ctx context.Context, dto CreateUserDTO) (User, error) {
	return User{}, nil
}
