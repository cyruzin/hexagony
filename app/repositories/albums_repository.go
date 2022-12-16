package postgres

import (
	"context"
	"database/sql"
	"errors"
	"hexagony/app/domain"
	"hexagony/app/repositories/queries"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type albumsRepository struct {
	conn *sqlx.DB
}

func NewAlbumsRepository(conn *sqlx.DB) domain.AlbumsRepository {
	return &albumsRepository{conn}
}

func (r *albumsRepository) FindAll(
	ctx context.Context,
) ([]*domain.Albums, error) {
	var albums []*domain.Albums

	err := r.conn.SelectContext(
		ctx,
		&albums,
		queries.SqlAlbumsFindAll,
	)
	noRows := errors.Is(err, sql.ErrNoRows)
	if noRows {
		return albums, nil
	}

	if err != nil {
		return nil, domain.ErrAlbumsFindAll
	}

	return albums, nil
}

func (r *albumsRepository) FindByID(
	ctx context.Context,
	uuid uuid.UUID,
) (*domain.Albums, error) {
	var album domain.Albums

	err := r.conn.GetContext(
		ctx,
		&album,
		queries.SqlAlbumsFindByID,
		uuid,
	)

	noRows := errors.Is(err, sql.ErrNoRows)
	if noRows {
		return nil, domain.ErrResourceNotFound
	}

	if err != nil {
		return nil, domain.ErrAlbumsFindByID
	}

	return &album, nil
}

func (r *albumsRepository) Add(
	ctx context.Context,
	album *domain.Albums,
) error {
	if _, err := r.conn.ExecContext(
		ctx,
		queries.SqlAlbumsAdd,
		album.UUID,
		album.Name,
		album.Length,
		album.CreatedAt,
		album.UpdatedAt,
	); err != nil {
		return domain.ErrAlbumsAdd
	}

	return nil
}

func (r *albumsRepository) Update(
	ctx context.Context,
	uuid uuid.UUID,
	album *domain.Albums,
) error {
	result, err := r.conn.ExecContext(
		ctx,
		queries.SqlAlbumsUpdate,
		album.Name,
		album.Length,
		album.UpdatedAt,
		uuid,
	)
	if err != nil {
		return domain.ErrAlbumsUpdate
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.ErrAlbumsUpdate
	}

	if rowsAffected == 0 {
		return domain.ErrResourceNotFound
	}

	return nil
}

func (r *albumsRepository) Delete(
	ctx context.Context,
	uuid uuid.UUID,
) error {
	result, err := r.conn.ExecContext(
		ctx,
		queries.SqlAlbumsDelete,
		uuid,
	)
	if err != nil {
		return domain.ErrAlbumsDelete
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return domain.ErrAlbumsDelete
	}

	if rowsAffected == 0 {
		return domain.ErrResourceNotFound
	}

	return nil
}
