package product

import (
	"context"
	"go.mod/pkg/logging"
)

type Service interface {
	Create(ctx context.Context, post CreateProductDTO) (*Product, error)
	Delete(ctx context.Context, postId int) error
	Update(ctx context.Context, post *Product, postUpdate UpdateProductDTO) (u *Product, err error)
	FindAll(ctx context.Context) ([]Product, error)
	FindOneById(ctx context.Context, id int) (u *Product, err error)
	FindUserPosts(ctx context.Context, userId int) ([]Product, error)
}

type postService struct {
	storage Storage
	logger  *logging.Logger
}

func NewPostService(storage Storage, logger *logging.Logger) Service {
	return &postService{
		storage: storage,
		logger:  logger,
	}
}

func (s *postService) Create(ctx context.Context, post CreateProductDTO) (*Product, error) {
	postObj, err := s.storage.Create(ctx, post)
	if err != nil {
		return nil, err
	}
	return postObj, nil
}

func (s *postService) Delete(ctx context.Context, postId int) error {
	err := s.storage.Delete(ctx, postId)
	if err != nil {
		return err
	}
	return nil
}

func (s *postService) Update(ctx context.Context, post *Product, postUpdate UpdateProductDTO) (u *Product, err error) {
	updated, err := s.storage.Update(ctx, post, postUpdate)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *postService) FindAll(ctx context.Context) ([]Product, error) {
	all, err := s.storage.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return all, nil
}

func (s *postService) FindOneById(ctx context.Context, id int) (u *Product, err error) {
	post, err := s.storage.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *postService) FindUserPosts(ctx context.Context, userId int) ([]Product, error) {
	posts, err := s.storage.FindUserAllProducts(ctx, userId)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
