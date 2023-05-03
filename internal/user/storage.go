package user

import (
	"context"
)

type Storage interface {
	Create(ctx context.Context, user User) error
	FindAll(ctx context.Context) (u []User, err error)
	FindOneById(ctx context.Context, id int) (User, error)
	FindOneByEmail(ctx context.Context, email string) (User, error)
	FindOneByUsername(ctx context.Context, username string) (User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id string) error
}
