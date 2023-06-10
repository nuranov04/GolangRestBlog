package category

import (
	"context"
	"go.mod/pkg/logging"
)

type Service interface {
	Create(ctx context.Context, createUser CreateUpdateCategory) (u *Category, err error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, updateDTO CreateUpdateCategory, categoryDTO Category) (u *Category, err error)
	FindAll(ctx context.Context) (u []Category, err error)
	FindOneById(ctx context.Context, id int) (u *Category, err error)
	FindOneByTitle(ctx context.Context, title string) (u *Category, err error)
}

type categoryService struct {
	storage Storage
	logger  *logging.Logger
}

func NewService(storage Storage, logger *logging.Logger) Service {
	return &categoryService{
		storage: storage,
		logger:  logger,
	}
}

func (c *categoryService) Create(ctx context.Context, createUser CreateUpdateCategory) (u *Category, err error) {
	create, err := c.storage.Create(ctx, createUser)
	if err != nil {
		return nil, err
	}
	return create, nil
}

func (c *categoryService) Delete(ctx context.Context, id int) error {
	return c.storage.Delete(ctx, id)
}

func (c *categoryService) Update(ctx context.Context, updateDTO CreateUpdateCategory, categoryDTO Category) (u *Category, err error) {
	update, err := c.storage.Update(ctx, updateDTO, categoryDTO)
	if err != nil {
		return nil, err
	}
	return update, nil
}

func (c *categoryService) FindAll(ctx context.Context) (u []Category, err error) {
	all, err := c.storage.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return all, nil
}

func (c *categoryService) FindOneById(ctx context.Context, id int) (u *Category, err error) {
	one, err := c.storage.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}
	return one, nil
}

func (c *categoryService) FindOneByTitle(ctx context.Context, title string) (u *Category, err error) {
	byTitle, err := c.storage.FindOneByTitle(ctx, title)
	if err != nil {
		return nil, err
	}
	return byTitle, nil
}
