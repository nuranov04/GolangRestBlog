package post

import (
	"context"
	"go.mod/pkg/logging"
)

type Service interface {
	Create(ctx context.Context, post CreatePostDTO) (*Post, error)
	Delete(ctx context.Context, postId int) error
	Update(ctx context.Context, post *Post, postUpdate UpdatePostDTO) (u *Post, err error)
	FindAll(ctx context.Context) ([]Post, error)
	FindOneById(ctx context.Context, id int) (u *Post, err error)
	FindUserPosts(ctx context.Context, userId int) ([]Post, error)
}

type postService struct {
	storage Storage
	logger  *logging.Logger
}

func NewPostService(storage Storage, logger *logging.Logger) *postService {
	return &postService{
		storage: storage,
		logger:  logger,
	}
}

func (s *postService) Create(ctx context.Context, post CreatePostDTO) (*Post, error) {
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

func (s *postService) Update(ctx context.Context, post *Post, postUpdate UpdatePostDTO) (u *Post, err error) {
	updated, err := s.storage.Update(ctx, post, postUpdate)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *postService) FindAll(ctx context.Context) ([]Post, error) {
	all, err := s.storage.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return all, nil
}

func (s *postService) FindOneById(ctx context.Context, id int) (u *Post, err error) {
	post, err := s.storage.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (s *postService) FindUserPosts(ctx context.Context, userId int) ([]Post, error) {
	posts, err := s.storage.FindUserAllPosts(ctx, userId)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
