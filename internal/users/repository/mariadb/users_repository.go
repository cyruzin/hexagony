package mariadb

import (
	"context"
	"database/sql"
	"hexagony/internal/users/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type mariadbRepository struct {
	conn *sqlx.DB
}

func NewMariaDBRepository(conn *sqlx.DB) domain.UserRepository {
	return &mariadbRepository{conn}
}

func (r *mariadbRepository) FindAll(
	ctx context.Context,
) ([]*domain.UserList, error) {
	var users []*domain.UserList

	err := r.conn.SelectContext(
		ctx,
		&users,
		sqlFindAll,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return users, nil
}

func (r *mariadbRepository) FindByID(
	ctx context.Context,
	uuid uuid.UUID,
) (*domain.UserList, error) {
	var user domain.UserList

	err := r.conn.GetContext(
		ctx,
		&user,
		sqlFindByID,
		uuid,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &user, nil
}

func (r *mariadbRepository) Add(
	ctx context.Context,
	user *domain.User,
) error {
	if _, err := r.conn.ExecContext(
		ctx,
		sqlAdd,
		user.UUID,
		user.Name,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (r *mariadbRepository) Update(
	ctx context.Context,
	uuid uuid.UUID,
	user *domain.User,
) error {
	result, err := r.conn.ExecContext(
		ctx,
		sqlUpdate,
		user.Name,
		user.Email,
		user.UpdatedAt,
		uuid,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrResourceNotFound
	}

	return nil
}

func (r *mariadbRepository) Delete(
	ctx context.Context,
	uuid uuid.UUID,
) error {
	result, err := r.conn.ExecContext(
		ctx,
		sqlDelete,
		uuid,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrResourceNotFound
	}

	return nil
}
