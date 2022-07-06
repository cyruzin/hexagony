package postgres

import (
	"context"
	"database/sql"
	"hexagony/internal/albums/domain"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type postgresRepository struct {
	conn *sqlx.DB
}

func NewPostgresRepository(conn *sqlx.DB) domain.AlbumRepository {
	return &postgresRepository{conn}
}

func (r *postgresRepository) FindAll(
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

func (r *postgresRepository) FindByID(
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

func (r *postgresRepository) Add(
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

func (r *postgresRepository) Update(
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
