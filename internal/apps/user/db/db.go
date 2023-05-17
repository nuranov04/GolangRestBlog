package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgconn"
	"go.mod/internal/apperror"
	"go.mod/internal/apps/user"
	"go.mod/pkg/client/postgresql"
	"go.mod/pkg/logging"
	"go.mod/pkg/utils"
)

type userRepository struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewUserRepository(client postgresql.Client, logger *logging.Logger) user.Storage {
	return &userRepository{
		client: client,
		logger: logger,
	}
}

func (r *userRepository) Create(ctx context.Context, userDTO user.User) (u *user.User, err error) {
	q := `INSERT INTO public.user (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id, username, email`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))
	if err := r.client.QueryRow(ctx, q, userDTO.Username, userDTO.Email, userDTO.Password).Scan(&userDTO.ID, &userDTO.Username, &userDTO.Email); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			if pgErr.Code == "23505" {
				return nil, apperror.UserAlreadyExist
			}
			return nil, newErr
		}
		return nil, err
	}

	return &userDTO, nil
}

func (r *userRepository) Update(ctx context.Context, userObj user.User, userUpdate user.UpdateUserDTO) (u *user.User, err error) {
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

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))
	if err := r.client.QueryRow(ctx, q, userUpdate.Username, userUpdate.Email, userUpdate.PasswordHash, userObj.ID).Scan(&userObj.ID, &userObj.Username, &userObj.Email, &userObj.Password); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			if pgErr.Code == "23505" {
				return nil, apperror.UserAlreadyExist
			}
			return nil, newErr
		}
	}
	return &userObj, nil
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	q := `
	DELETE FROM public.user WHERE id=$1
	`

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

func (r *userRepository) FindAll(ctx context.Context) (u []user.User, err error) {
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

func (r *userRepository) FindOneById(ctx context.Context, id int) (u *user.User, err error) {
	q := `
		SELECT id, username, email FROM public.user WHERE id = $1
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))

	var userInfo user.User
	if err := r.client.QueryRow(ctx, q, id).Scan(&userInfo.ID, &userInfo.Username, &userInfo.Email); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			return nil, newErr
		}
		return nil, err
	}
	return &userInfo, nil

}

func (r *userRepository) FindOneByUsername(ctx context.Context, username string) (u *user.User, err error) {
	q := `
		SELECT id, username, email, password_hash FROM public.user WHERE username = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))

	var userInfo user.User

	err = r.client.QueryRow(ctx, q, username).Scan(&userInfo.ID, &userInfo.Username, &userInfo.Email, &userInfo.Password)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil

}

func (r *userRepository) FindOneByEmail(ctx context.Context, email string) (u *user.User, err error) {
	q := `
		SELECT id, username, email FROM public.user WHERE email = $1
	`
	r.logger.Trace(fmt.Sprintf("SQL Query: %s", utils.FormatQuery(q)))

	var userInfo user.User

	err = r.client.QueryRow(ctx, q, email).Scan(&userInfo.ID, &userInfo.Username, &userInfo.Email)
	if err != nil {
		return nil, err
	}

	return &userInfo, nil

}
