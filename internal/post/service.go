package post

import (
	"context"
	"go.mod/pkg/logging"
)

type Service interface {
	Create(ctx context.Context, post Post) error
	Delete(ctx context.Context, postId int) error
	Update(ctx context.Context, post Post) (Post, error)
	FindAll(ctx context.Context) ([]Post, error)
	FindOneById(ctx context.Context, id int) (Post, error)
	FindUserPosts(ctx context.Context, userId int) ([]Post, error)
}

type service struct {
	storage Storage
	logger  *logging.Logger
}

func (s *service) Create(ctx context.Context, post Post) error {
	//TODO implement me
	panic("implement me")
}

func (s *service) Delete(ctx context.Context, postId int) error {
	//TODO implement me
	panic("implement me")
}

func (s *service) Update(ctx context.Context, post Post) (Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) FindAll(ctx context.Context) ([]Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) FindOneById(ctx context.Context, id int) (Post, error) {
	//TODO implement me
	panic("implement me")
}

func (s *service) FindUserPosts(ctx context.Context, userId int) ([]Post, error) {
	//TODO implement me
	panic("implement me")
}
