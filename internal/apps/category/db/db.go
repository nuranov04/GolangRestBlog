package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"go.mod/internal/apperror"
	"go.mod/internal/apps/category"
	"go.mod/pkg/client/postgresql"
	"go.mod/pkg/logging"
	"go.mod/pkg/utils"
)

type categoryRepository struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewCategoryRepository(client postgresql.Client, logger *logging.Logger) category.Storage {
	return &categoryRepository{
		client: client,
		logger: logger,
	}
}

func (r *categoryRepository) FindOne(ctx context.Context, id int) (c *category.Category, err error) {
	q := `
	SELECT id, title, child_id FROM public.category WHERE id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))
	var categoryDTO category.Category
	if err := r.client.QueryRow(ctx, q, id).Scan(&categoryDTO.Id, &categoryDTO.Title, &categoryDTO.ChildId); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			return nil, newErr
		}
		return nil, err
	}
	return &categoryDTO, nil
}
func (r *categoryRepository) FindOneByTitle(ctx context.Context, title string) (c *category.Category, err error) {
	q := `
	SELECT id, title, child_id FROM public.category WHERE title = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))
	var categoryDTO category.Category
	if err := r.client.QueryRow(ctx, q, title).Scan(&categoryDTO.Id, &categoryDTO.Title, &categoryDTO.ChildId); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			return nil, newErr
		}
		return nil, err
	}
	return &categoryDTO, nil
}

func (r *categoryRepository) FindAll(ctx context.Context) (c []category.Category, err error) {
	q := `
	SELECT id, title, child_id FROM public.category
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))
	query, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	categories := make([]category.Category, 0)

	for query.Next() {
		var categoryInfo category.Category
		err := query.Scan(&categoryInfo.Id, &categoryInfo.Title, &categoryInfo.ChildId)
		if err != nil {
			return nil, err
		}
		categories = append(categories, categoryInfo)

	}
	if err = query.Err(); err != nil {
		return nil, err
	}
	return categories, err
}

func (r *categoryRepository) Create(ctx context.Context, categoryDTO category.CreateUpdateCategory) (c *category.Category, err error) {
	q := `
	INSERT INTO public.category (title, child_id) VALUES ($1, $2) RETURNING id, title, child_id 
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))
	var categoryInfo category.Category
	if err := r.client.QueryRow(ctx, q, categoryDTO.Title, categoryDTO.ChildId).
		Scan(&categoryInfo.Id, &categoryDTO.Title, &categoryDTO.ChildId); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			if pgErr.Code == "23505" {
				return nil, apperror.CategoryTileAlreadyExist
			}
			return nil, newErr
		}
		return nil, err
	}
	return &categoryInfo, nil
}

func (r *categoryRepository) Update(ctx context.Context, categoryUpdate category.CreateUpdateCategory, categoryDTO category.Category) (c *category.Category, err error) {
	q := `
	UPDATE public.category 
	SET title=$1, child_id = $2
	WHERE id = (
	    SELECT id 
	    FROM public.category
	    WHERE id = $3
	    LIMIT 1
	    FOR UPDATE 
	)
	RETURNING id, title, child_id;`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))
	if err := r.client.QueryRow(ctx, q, categoryUpdate.Title, categoryUpdate.ChildId, categoryDTO.Id).Scan(&categoryDTO.Id, &categoryDTO.Title, &categoryDTO.ChildId); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			if pgErr.Code == "23505" {
				return nil, apperror.CategoryTileAlreadyExist
			}
			return nil, newErr
		}
		return nil, err
	}
	return nil, err
}

func (r *categoryRepository) Delete(ctx context.Context, id int) error {
	q := `	
	DELETE FROM public.category WHERE id = $1;`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))
	if err := r.client.QueryRow(ctx, q, id).Scan(&id); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			return newErr
		}
		return err
	}
	return nil
}
