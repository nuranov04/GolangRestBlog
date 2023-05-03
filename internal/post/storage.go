package post

import "context"

type Storage interface {
	Create(ctx context.Context, post CreatePostDTO) (string, error)
	FindOne(ctx context.Context, id string) (u *Post, err error)
	FindAll(ctx context.Context) (u []*Post, err error)
	Update(ctx context.Context, post *Post) error
	Delete(ctx context.Context, id string) error
}
