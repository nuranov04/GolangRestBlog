package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"go.mod/internal/post"
	"go.mod/pkg/client/postgresql"
	"go.mod/pkg/logging"
	"go.mod/pkg/utils"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewRepository(client postgresql.Client, logger *logging.Logger) post.Storage {
	return &repository{
		client: client,
		logger: logger,
	}
}

func (r *repository) Create(ctx context.Context, postObj post.CreatePostDTO) (p *post.CreatePostDTO, err error) {
	q := `
	INSERT INTO (title, description, owner_id) VALUES ($1, $2, $3) RETURNING id, title, description, owner_id
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))
	if err := r.client.QueryRow(ctx, q, postObj.Title, postObj.Description, postObj.OwnerId).Scan(&postObj.Id, postObj.Title, postObj.Description, postObj.OwnerId); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return nil, newErr
		}
		return nil, err
	}
	return &postObj, nil
}

func (r *repository) FindOne(ctx context.Context, id int) (u *post.Post, err error) {
	q := `SELECT id, title, description, owner_id FROM public.post WHERE id = $1`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))
	var userObj post.Post
	if err := r.client.QueryRow(ctx, q, id).Scan(&userObj.ID, userObj.Title, userObj.Description, userObj.OwnerId); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return &userObj, newErr
		}
		return &userObj, err
	}
	return &userObj, err
}

func (r *repository) FindAll(ctx context.Context) (u []post.Post, err error) {
	q := `SELECT id, title, description, owner_id FROM public.post`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))
	query, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	posts := make([]post.Post, 0)

	for query.Next() {
		var postInfo post.Post
		err := query.Scan(&postInfo.ID, &postInfo.Title, &postInfo.Description, &postInfo.OwnerId)
		if err != nil {
			return nil, err
		}

		posts = append(posts, postInfo)
	}
	return posts, nil
}

func (r *repository) FindUserAllPosts(ctx context.Context, userId int) ([]post.Post, error) {
	q := `
			SELECT id, title, description, owner_id FROM public.post WHERE owner_id = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))
	query, err := r.client.Query(ctx, q, userId)
	if err != nil {
		return nil, err
	}
	defer query.Close()

	posts := make([]post.Post, 0)
	for query.Next() {
		var postInfo post.Post
		err := query.Scan(&postInfo.ID, &postInfo.Title, &postInfo.Description, &postInfo.OwnerId)
		if err != nil {
			return nil, err
		}

		posts = append(posts, postInfo)
	}
	return posts, nil
}

func (r *repository) Update(ctx context.Context, postObj *post.Post, postUpdate post.UpdatePostDTO) (u *post.Post, err error) {
	q := `
		UPDATE public.post 
		SET title = $1, description = $2 
		WHERE id = (
			SELECT id
			FROM public.post
			WHERE id = $3
			AND owner_id = $4
			LIMIT 1
			FOR UPDATE 
		)
	RETURNING title, description;`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))

	if err := r.client.QueryRow(ctx, q, postUpdate.Title, postUpdate.Description, postObj.ID, postObj.OwnerId).Scan(&postObj.Title, &postObj.Description); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return &post.Post{}, newErr
		}
		return &post.Post{}, err
	}
	return postObj, nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	panic("das")
}
