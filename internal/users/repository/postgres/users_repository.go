package postgres

import (
	"context"
	"database/sql"
	"errors"
	"hexagony/internal/users/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type postgresRepository struct {
	conn *sqlx.DB
}

func NewPostgresRepository(conn *sqlx.DB) domain.UserRepository {
	return &postgresRepository{conn}
}

func (r *postgresRepository) FindAll(
	ctx context.Context,
) ([]*domain.UserList, error) {
	var users []*domain.UserList

	err := r.conn.SelectContext(
		ctx,
		&users,
		sqlFindAll,
	)
	noRows := errors.Is(err, sql.ErrNoRows)
	if noRows {
		return users, nil
	}

	if err != nil {
		return nil, domain.ErrFindAll
	}

	return users, nil
}

func (r *postgresRepository) FindByID(
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
	noRows := errors.Is(err, sql.ErrNoRows)
	if noRows {
		return nil, domain.ErrResourceNotFound
	}

	if err != nil {
		return nil, domain.ErrFindByID
	}

	return &user, nil
}

func (r *postgresRepository) Add(
	ctx context.Context,
	user *domain.User,
) error {
	exists, err := r.checkDuplicate(ctx, user.Email)
	if err != nil {
		return domain.ErrAdd
	}

	if exists {
		return domain.ErrDuplicateEmail
	}

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
		return domain.ErrAdd
	}

	return nil
}

func (r *postgresRepository) Update(
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
		return domain.ErrUpdate
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.ErrUpdate
	}

	if rowsAffected == 0 {
		return domain.ErrResourceNotFound
	}

	return nil
}

func (r *postgresRepository) Delete(
	ctx context.Context,
	uuid uuid.UUID,
) error {
	result, err := r.conn.ExecContext(
		ctx,
		sqlDelete,
		uuid,
	)
	if err != nil {
		return domain.ErrDelete
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.ErrDelete
	}

	if rowsAffected == 0 {
		return domain.ErrResourceNotFound
	}

	return nil
}

func (r *postgresRepository) checkDuplicate(ctx context.Context, email string) (bool, error) {
	var userEmail string
	exists := false

	err := r.conn.GetContext(ctx, &userEmail, sqlCheckDuplicate, email)
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
