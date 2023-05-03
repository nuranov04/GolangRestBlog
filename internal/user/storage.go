package user

import "context"

type Storage interface {
	Create(ctx context.Context, user User) error
	FindOneById(ctx context.Context, id int) (User, error)
	FindAll(ctx context.Context) (u []User, err error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id string) error
}
