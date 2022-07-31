package postgres

import (
	"context"
	"database/sql"
	"errors"
	"hexagony/app/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type usersRepository struct {
	conn *sqlx.DB
}

func NewUsersRepository(conn *sqlx.DB) domain.UsersRepository {
	return &usersRepository{conn}
}

func (r *usersRepository) FindAll(
	ctx context.Context,
) ([]*domain.UsersList, error) {
	var users []*domain.UsersList

	err := r.conn.SelectContext(
		ctx,
		&users,
		sqlUsersFindAll,
	)
	noRows := errors.Is(err, sql.ErrNoRows)
	if noRows {
		return users, nil
	}

	if err != nil {
		return nil, domain.ErrUsersFindAll
	}

	return users, nil
}

func (r *usersRepository) FindByID(
	ctx context.Context,
	uuid uuid.UUID,
) (*domain.UsersList, error) {
	var user domain.UsersList

	err := r.conn.GetContext(
		ctx,
		&user,
		sqlUsersFindByID,
		uuid,
	)
	noRows := errors.Is(err, sql.ErrNoRows)
	if noRows {
		return nil, domain.ErrResourceNotFound
	}

	if err != nil {
		return nil, domain.ErrUsersFindByID
	}

	return &user, nil
}

func (r *usersRepository) Add(
	ctx context.Context,
	user *domain.Users,
) error {
	exists, err := r.checkDuplicate(ctx, user.Email)
	if err != nil {
		return err
	}

	if exists {
		return domain.ErrUsersDuplicateEmail
	}

	if _, err := r.conn.ExecContext(
		ctx,
		sqlUsersAdd,
		user.UUID,
		user.Name,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	); err != nil {
		return domain.ErrUsersAdd
	}

	return nil
}

func (r *usersRepository) Update(
	ctx context.Context,
	uuid uuid.UUID,
	user *domain.Users,
) error {
	result, err := r.conn.ExecContext(
		ctx,
		sqlUsersUpdate,
		user.Name,
		user.Email,
		user.UpdatedAt,
		uuid,
	)
	if err != nil {
		return domain.ErrUsersUpdate
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.ErrUsersUpdate
	}

	if rowsAffected == 0 {
		return domain.ErrResourceNotFound
	}

	return nil
}

func (r *usersRepository) Delete(
	ctx context.Context,
	uuid uuid.UUID,
) error {
	result, err := r.conn.ExecContext(
		ctx,
		sqlUsersDelete,
		uuid,
	)
	if err != nil {
		return domain.ErrUsersDelete
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.ErrUsersDelete
	}

	if rowsAffected == 0 {
		return domain.ErrResourceNotFound
	}

	return nil
}

func (r *usersRepository) checkDuplicate(ctx context.Context, email string) (bool, error) {
	var userEmail string
	exists := false

	err := r.conn.GetContext(ctx, &userEmail, sqlUsersCheckDuplicate, email)
	noRows := errors.Is(err, sql.ErrNoRows)
	if noRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	if userEmail != "" {
		exists = true
	}

	return exists, nil
}
