package postgres

import (
	"context"
	"database/sql"
	"errors"
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
	noRows := errors.Is(err, sql.ErrNoRows)
	if noRows {
		return albums, nil
	}

	if err != nil {
		return nil, domain.ErrFindAll
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

	noRows := errors.Is(err, sql.ErrNoRows)
	if noRows {
		return nil, domain.ErrResourceNotFound
	}

	if err != nil {
		return nil, domain.ErrFindByID
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
		return domain.ErrAdd
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
