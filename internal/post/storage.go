package post

import "context"

type Storage interface {
	Create(ctx context.Context, post CreatePostDTO) (u *CreatePostDTO, err error)
	FindOne(ctx context.Context, id string) (u *Post, err error)
	FindAll(ctx context.Context) (u []Post, err error)
	FindUserAllPosts(ctx context.Context, userId int) ([]Post, error)
	Update(ctx context.Context, post *Post) error
	Delete(ctx context.Context, id string) error
}
