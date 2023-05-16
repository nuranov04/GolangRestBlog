package storage

import (
	"context"
	"go.mod/internal/model"
)

type UserStorage interface {
	Create(ctx context.Context, user model.User) (u *model.User, err error)
	FindAll(ctx context.Context) (u []model.User, err error)
	FindOneById(ctx context.Context, id int) (u *model.User, err error)
	FindOneByEmail(ctx context.Context, email string) (u *model.User, err error)
	FindOneByUsername(ctx context.Context, username string) (u *model.User, err error)
	Update(ctx context.Context, user model.User, userUpdate model.UpdateUserDTO) (u *model.User, err error)
	Delete(ctx context.Context, id int) error
}
