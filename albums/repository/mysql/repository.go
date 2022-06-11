package mysql

import (
	"context"
	"database/sql"
	"hexagony/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type mysqlRepository struct {
	conn *sqlx.DB
}

func NewMysqlRepository(conn *sqlx.DB) domain.AlbumRepository {
	return &mysqlRepository{conn}
}

func (r *mysqlRepository) FindAll(
	ctx context.Context,
) ([]*domain.Album, error) {
	var album []*domain.Album

	err := r.conn.SelectContext(
		ctx,
		&album,
		sqlFindAll,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return album, nil
}

func (r *mysqlRepository) FindByID(
	ctx context.Context,
	uuid uuid.UUID,
) (*domain.Album, error) {
	var album domain.Album

	err := r.conn.GetContext(
		ctx,
		&album,
		sqlFindByID,
		uuid,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if album.Name == "" {
		return nil, domain.ErrResourceNotFound
	}

	return &album, nil
}

func (r *mysqlRepository) Add(
	ctx context.Context,
	album *domain.Album,
) error {
	if _, err := r.conn.ExecContext(
		ctx,
		sqlAdd,
		album.UUID,
		album.Name,
		album.Length,
		album.CreatedAt,
		album.CreatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (r *mysqlRepository) Update(
	ctx context.Context,
	uuid uuid.UUID,
	album *domain.Album,
) error {
	result, err := r.conn.ExecContext(
		ctx,
		sqlUpdate,
		album.Name,
		album.Length,
		album.UpdatedAt,
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

func (r *mysqlRepository) Delete(
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
