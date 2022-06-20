package mariadb

import (
	"context"
	"database/sql"
	"hexagony/app/albums/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type mariadbRepository struct {
	conn *sqlx.DB
}

func NewMariaDBRepository(conn *sqlx.DB) domain.AlbumRepository {
	return &mariadbRepository{conn}
}

func (r *mariadbRepository) FindAll(
	ctx context.Context,
) ([]*domain.Album, error) {
	var albums []*domain.Album

	err := r.conn.SelectContext(
		ctx,
		&albums,
		sqlFindAll,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return albums, nil
}

func (r *mariadbRepository) FindByID(
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

	return &album, nil
}

func (r *mariadbRepository) Add(
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
		album.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}

func (r *mariadbRepository) Update(
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
