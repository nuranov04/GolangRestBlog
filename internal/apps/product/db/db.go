package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"go.mod/internal/apperror"
	"go.mod/internal/apps/product"
	"go.mod/pkg/client/postgresql"
	"go.mod/pkg/logging"
	"go.mod/pkg/utils"
)

type ProductRepository struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewProductRepository(client postgresql.Client, logger *logging.Logger) product.Storage {
	return &ProductRepository{
		client: client,
		logger: logger,
	}
}

func (r *ProductRepository) Create(ctx context.Context, ProductObj product.CreateProductDTO) (u *product.Product, err error) {
	q := `
	INSERT INTO public.product (title, description, owner_id) VALUES ($1, $2, $3) RETURNING id, title, description, owner_id
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))
	var ProductDTO product.Product
	if err := r.client.QueryRow(ctx, q, ProductObj.Title, ProductObj.Description, ProductObj.OwnerId).Scan(&ProductDTO.ID, &ProductDTO.Title, &ProductDTO.Description, &ProductDTO.OwnerId); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			if pgErr.Code == "23505" {
				return nil, apperror.PostTitleAlreadyExist
			}
			return nil, newErr
		}
		return nil, err
	}
	return &ProductDTO, nil
}

func (r *ProductRepository) Update(ctx context.Context, ProductObj *product.Product, ProductUpdate product.UpdateProductDTO) (u *product.Product, err error) {
	q := `
		UPDATE public.product 
		SET title = $1, description = $2 
		WHERE id = (
			SELECT id
			FROM public.product
			WHERE id = $3
			AND owner_id = $4
			LIMIT 1
			FOR UPDATE 
		)
	RETURNING title, description;`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))

	if err := r.client.QueryRow(ctx, q, ProductUpdate.Title, ProductUpdate.Description, ProductObj.ID, ProductObj.OwnerId).Scan(&ProductObj.Title, &ProductObj.Description); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			if pgErr.Code == "23505" {
				return nil, apperror.PostTitleAlreadyExist
			}
			return nil, newErr
		}
		return nil, err
	}
	return ProductObj, nil
}

func (r *ProductRepository) FindOne(ctx context.Context, id int) (u *product.Product, err error) {
	q := `SELECT id, title, description, owner_id FROM public.product WHERE id = $1`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))
	var ProductObj product.Product
	if err := r.client.QueryRow(ctx, q, id).Scan(&ProductObj.ID, &ProductObj.Title, &ProductObj.Description, &ProductObj.OwnerId); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			return nil, newErr
		}
		r.logger.Debug(err)
		return nil, err
	}
	return &ProductObj, nil
}

func (r *ProductRepository) FindAll(ctx context.Context) (u []product.Product, err error) {
	q := `SELECT id, title, description, owner_id FROM public.product`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))
	query, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	Products := make([]product.Product, 0)
	r.logger.Debug(query)
	for query.Next() {
		var ProductInfo product.Product
		err := query.Scan(&ProductInfo.ID, &ProductInfo.Title, &ProductInfo.Description, &ProductInfo.OwnerId)
		if err != nil {
			return nil, err
		}

		Products = append(Products, ProductInfo)
	}
	return Products, nil
}

func (r *ProductRepository) FindUserAllProducts(ctx context.Context, userId int) ([]product.Product, error) {
	q := `
			SELECT id, title, description, owner_id FROM public.product WHERE owner_id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))
	query, err := r.client.Query(ctx, q, userId)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	Products := make([]product.Product, 0)
	for query.Next() {
		var ProductInfo product.Product
		err := query.Scan(&ProductInfo.ID, &ProductInfo.Title, &ProductInfo.Description, &ProductInfo.OwnerId)
		if err != nil {
			return nil, err
		}

		Products = append(Products, ProductInfo)
	}
	return Products, nil
}

func (r *ProductRepository) Delete(ctx context.Context, id int) error {
	q := `
	DELETE FROM public.product WHERE id=$1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))
	if err := r.client.QueryRow(ctx, q, id).Scan(&id); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			return newErr
		}
	}
	return nil
}
