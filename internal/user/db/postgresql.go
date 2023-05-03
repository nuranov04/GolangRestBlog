package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"go.mod/internal/user"
	"go.mod/pkg/client/postgresql"
	"go.mod/pkg/logging"
	"strings"
)

type repository struct {
	client postgresql.Client
	logger *logging.Logger
}

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func NewRepository(client postgresql.Client, logger *logging.Logger) user.Storage {
	return &repository{
		client: client,
		logger: logger,
	}
}

func (r *repository) Create(ctx context.Context, user user.User) error {
	q := `INSERT INTO public.user (username, email, password_hash) VALUES ($1, $2, $3) `
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	if err := r.client.QueryRow(ctx, q, user.Username, user.Email, user.PasswordHash).Scan(&user.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return newErr
		}
		return err
	}

	return nil
}

func (r *repository) FindOneById(ctx context.Context, id int) (user.User, error) {
	q := `
		SELECT id, username, email FROM public.user WHERE id = $1
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var userInfo user.User
	err := r.client.QueryRow(ctx, q, id).Scan(&userInfo.ID, &userInfo.Username, &userInfo.Email)
	r.logger.Info(err)
	if err != nil {
		return user.User{}, err
	}

	return userInfo, nil

}

func (r *repository) FindAll(ctx context.Context) (u []user.User, err error) {
	q := `SELECT id, username, email FROM public.user`
	query, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	users := make([]user.User, 0)

	for query.Next() {
		var userInfo user.User
		err := query.Scan(&userInfo.ID, &userInfo.Username, &userInfo.Email)
		if err != nil {
			return nil, err
		}

		users = append(users, userInfo)
	}

	if err = query.Err(); err != nil {
		return nil, err
	}

	return users, nil

}

func (r *repository) Update(ctx context.Context, user user.User) error {
	//TODO implement me
	panic("implement me")
}

func (r *repository) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
