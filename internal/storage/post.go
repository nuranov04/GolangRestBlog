package storage

import (
	"context"
	"go.mod/internal/model"
)

type PostStorage interface {
	Create(ctx context.Context, post model.CreatePostDTO) (u *model.Post, err error)
	FindOne(ctx context.Context, id int) (u *model.Post, err error)
	FindAll(ctx context.Context) (u []model.Post, err error)
	FindUserAllPosts(ctx context.Context, userId int) ([]model.Post, error)
	Update(ctx context.Context, postObj *model.Post, postUpdate model.UpdatePostDTO) (u *model.Post, err error)
	Delete(ctx context.Context, id int) error
}
