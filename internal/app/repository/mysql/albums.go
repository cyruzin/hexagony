package mysql

import (
	"context"
	"database/sql"
	"hexagony/internal/app/albums"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	// Mysql connection
	_ "github.com/go-sql-driver/mysql"
)

type mysqlRepository struct {
	client *sqlx.DB
}

// NewMysqlRepository creates a instance of MySQL that access the albums repository.
func NewMysqlRepository(
	ctx context.Context,
	dataSourceName string,
) (albums.Repository, error) {
	client, err := sqlx.ConnectContext(ctx, "mysql", dataSourceName)
	if err != nil {
		return nil, err
	}

	if err := client.PingContext(ctx); err != nil {
		return nil, err
	}

	return &mysqlRepository{client}, nil
}

// FindAll finds the latest albums.
func (r *mysqlRepository) FindAll(
	ctx context.Context,
) ([]*albums.Album, error) {
	var album []*albums.Album

	err := r.client.SelectContext(
		ctx,
		&album,
		sqlFindAll,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return album, nil
}

// FindByID finds an album by ID.
func (r *mysqlRepository) FindByID(
	ctx context.Context,
	uuid uuid.UUID,
) (*albums.Album, error) {
	var album albums.Album

	err := r.client.GetContext(
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

// Add adds a new album.
func (r *mysqlRepository) Add(
	ctx context.Context,
	album *albums.Album,
) error {
	if _, err := r.client.ExecContext(
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

// Update updates an album by ID.
func (r *mysqlRepository) Update(
	ctx context.Context,
	uuid uuid.UUID,
	album *albums.Album,
) error {
	if _, err := r.client.ExecContext(
		ctx,
		sqlUpdate,
		album.Name,
		album.Length,
		album.UpdatedAt,
		uuid,
	); err != nil {
		return err
	}

	return nil
}

// Delete deletes an album by ID.
func (r *mysqlRepository) Delete(
	ctx context.Context,
	uuid uuid.UUID,
) error {
	if _, err := r.client.ExecContext(
		ctx,
		sqlDelete,
		uuid,
	); err != nil {
		return err
	}

	return nil
}
