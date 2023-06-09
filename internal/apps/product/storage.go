package product

import (
	"context"
)

type Storage interface {
	Create(ctx context.Context, post CreateProductDTO) (u *Product, err error)
	FindOne(ctx context.Context, id int) (u *Product, err error)
	FindAll(ctx context.Context) (u []Product, err error)
	FindUserAllProducts(ctx context.Context, userId int) ([]Product, error)
	Update(ctx context.Context, postObj *Product, postUpdate UpdateProductDTO) (u *Product, err error)
	Delete(ctx context.Context, id int) error
}
