package category

import "context"

type Storage interface {
	FindOne(ctx context.Context, id int) (c *Category, err error)
	FindAll(ctx context.Context) (c []Category, err error)
	Create(ctx context.Context, categoryDTO CreateUpdateCategory) (c *Category, err error)
	Update(ctx context.Context, categoryUpdate CreateUpdateCategory, category Category) (c *Category, err error)
	Delete(ctx context.Context, id int) error
}
