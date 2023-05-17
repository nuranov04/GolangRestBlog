package user

import (
	"context"
)

type Storage interface {
	Create(ctx context.Context, user User) (u *User, err error)
	FindAll(ctx context.Context) (u []User, err error)
	FindOneById(ctx context.Context, id int) (u *User, err error)
	FindOneByEmail(ctx context.Context, email string) (u *User, err error)
	FindOneByUsername(ctx context.Context, username string) (u *User, err error)
	Update(ctx context.Context, user User, userUpdate UpdateUserDTO) (u *User, err error)
	Delete(ctx context.Context, id int) error
}
