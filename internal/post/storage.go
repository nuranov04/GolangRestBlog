package post

import "context"

type Storage interface {
	Create(ctx context.Context, post CreatePostDTO) (u *Post, err error)
	FindOne(ctx context.Context, id int) (u *Post, err error)
	FindAll(ctx context.Context) (u []Post, err error)
	FindUserAllPosts(ctx context.Context, userId int) ([]Post, error)
	Update(ctx context.Context, postObj *Post, postUpdate UpdatePostDTO) (u *Post, err error)
	Delete(ctx context.Context, id int) error
}
