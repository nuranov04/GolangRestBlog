package db

import (
	"context"
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
	q := `INSERT INTO public.user (username, email, password_hash) VALUES ($1, $2, $3) RETURNING username, email`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	if err := r.client.QueryRow(ctx, q, user.Username, user.Email, user.PasswordHash).Scan(&user.Username, &user.Email); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return newErr
		}
		return err
	}

	return nil
}

func (r *repository) Update(ctx context.Context, userObj user.User, userUpdate user.UpdateUserDTO) (user.UpdateUserDTO, error) {
	q := `
	UPDATE public.user 
	SET username = $1, email = $2, password_hash = $3 
	WHERE id = (
	    SELECT id
	    FROM public.user
	    WHERE id = $4
	    LIMIT 1
	    FOR UPDATE 
	)
	RETURNING id, username, email, password_hash;`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	if err := r.client.QueryRow(ctx, q, userUpdate.Username, userUpdate.Email, userUpdate.PasswordHash, userObj.ID).Scan(&userUpdate.Username, &userUpdate.Email); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return userUpdate, newErr
		}
	}
	return userUpdate, nil
}

func (r *repository) Delete(ctx context.Context, id int) error {
	q := `
	DELETE FROM public.user WHERE id=$1
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	if err := r.client.QueryRow(ctx, q, id).Scan(&id); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return newErr
		}
	}
	return nil

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

func (r *repository) FindOneById(ctx context.Context, id int) (user.User, error) {
	q := `
		SELECT id, username, email FROM public.user WHERE id = $1
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var userInfo user.User
	err := r.client.QueryRow(ctx, q, id).Scan(&userInfo.ID, &userInfo.Username, &userInfo.Email)
	if err != nil {
		return user.User{}, err
	}

	return userInfo, nil

}

func (r *repository) FindOneByUsername(ctx context.Context, username string) (user.User, error) {
	q := `
		SELECT id, username, email FROM public.user WHERE username = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var userInfo user.User

	err := r.client.QueryRow(ctx, q, username).Scan(&userInfo.ID, &userInfo.Username, &userInfo.Email)
	if err != nil {
		return user.User{}, err
	}

	return userInfo, nil

}

func (r *repository) FindOneByEmail(ctx context.Context, email string) (user.User, error) {
	q := `
		SELECT id, username, email FROM public.user WHERE email = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", formatQuery(q)))

	var userInfo user.User

	err := r.client.QueryRow(ctx, q, email).Scan(&userInfo.ID, &userInfo.Username, &userInfo.Email)
	if err != nil {
		return user.User{}, err
	}

	return userInfo, nil

}
