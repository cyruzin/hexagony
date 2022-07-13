package postgres

import (
	"context"
	"database/sql"
	"hexagony/internal/albums/domain"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestFindAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	mockAlbums := []domain.Album{
		{
			UUID:      uuid.New(),
			Name:      "St. Anger", // Metallica
			Length:    75,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			UUID:      uuid.New(),
			Name:      "De Profundis", // Vader
			Length:    34,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{
		"uuid",
		"name",
		"length",
		"created_at",
		"updated_at",
	}).AddRow(
		mockAlbums[0].UUID,
		mockAlbums[0].Name,
		mockAlbums[0].Length,
		mockAlbums[0].CreatedAt,
		mockAlbums[0].UpdatedAt,
	).AddRow(
		mockAlbums[1].UUID,
		mockAlbums[1].Name,
		mockAlbums[1].Length,
		mockAlbums[1].CreatedAt,
		mockAlbums[1].UpdatedAt,
	)

	query := "SELECT \\* FROM albums"

	mock.ExpectQuery(query).WillReturnRows(rows)

	albumRepo := NewPostgresRepository(dbx)
	albumList, err := albumRepo.FindAll(context.TODO())

	assert.NoError(t, err)
	assert.Len(t, albumList, 2)
	assert.Equal(t, albumList[0].Name, "St. Anger")
}

func TestFindAllFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	rows := sqlmock.NewRows([]string{
		"uuid",
		"name",
		"length",
		"created_at",
		"updated_at",
	}).
		AddRow("", "", "", "", "")

	query := "SELECT \\* FROM albums"
	mock.ExpectQuery(query).WillReturnRows(rows)

	albumRepo := NewPostgresRepository(dbx)
	_, err = albumRepo.FindAll(context.TODO())

	assert.NotNil(t, err)
}

func TestFindAllFailNoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	rows := sqlmock.NewRows([]string{
		"uuid",
		"name",
		"length",
		"created_at",
		"updated_at",
	}).
		AddRow("", "", "", "", "")

	query := "SELECT \\* FROM albums"
	mock.ExpectQuery(query).WillReturnRows(rows).WillReturnError(sql.ErrNoRows)

	albumRepo := NewPostgresRepository(dbx)
	albums, _ := albumRepo.FindAll(context.TODO())

	assert.Nil(t, albums)
}

func TestFindByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	newUUID := uuid.New()

	rows := sqlmock.NewRows([]string{
		"uuid",
		"name",
		"length",
		"created_at",
		"updated_at",
	}).
		AddRow(newUUID, "St. Anger", 75, time.Now(), time.Now())

	query := "SELECT \\* FROM albums WHERE uuid=\\$1"
	mock.ExpectQuery(query).WillReturnRows(rows)

	albumRepo := NewPostgresRepository(dbx)
	currentAlbum, err := albumRepo.FindByID(context.TODO(), newUUID)

	assert.NoError(t, err)
	assert.NotNil(t, currentAlbum)
	assert.Equal(t, "St. Anger", currentAlbum.Name)
}

func TestGetByIDFail(t *testing.T) {
	newUUID := uuid.New()
	ctx := context.TODO()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	rows := sqlmock.NewRows([]string{
		"uuid",
		"name",
		"length",
		"created_at",
		"updated_at",
	}).
		AddRow("", "", "", "", "")

	query := "SELECT \\* FROM albums WHERE uuid=\\$1"
	mock.ExpectQuery(query).WillReturnRows(rows)

	albumRepo := NewPostgresRepository(dbx)
	_, err = albumRepo.FindByID(ctx, newUUID)

	assert.NotNil(t, err)
}

func TestGetByIDFailNoRows(t *testing.T) {
	newUUID := uuid.New()
	ctx := context.TODO()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	rows := sqlmock.NewRows([]string{
		"uuid",
		"name",
		"length",
		"created_at",
		"updated_at",
	}).
		AddRow("", "", "", "", "")

	query := "SELECT \\* FROM albums WHERE uuid=\\$1"
	mock.ExpectQuery(query).WillReturnRows(rows).WillReturnError(sql.ErrNoRows)

	albumRepo := NewPostgresRepository(dbx)
	_, err = albumRepo.FindByID(ctx, newUUID)

	assert.NotNil(t, err)
}

func TestAdd(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	album := &domain.Album{
		UUID:      newUUID,
		Name:      "St. Anger",
		Length:    75,
		CreatedAt: now,
		UpdatedAt: now,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	query := `INSERT INTO 
	albums (uuid, name, length, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5)`

	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(newUUID, album.Name, album.Length, album.CreatedAt, album.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1)) // Using UUID

	albumRepo := NewPostgresRepository(dbx)
	err = albumRepo.Add(context.TODO(), album)

	assert.NoError(t, err)
}

func TestAddFail(t *testing.T) {
	album := &domain.Album{}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	query := `
	  INSERT INTO albums (
		uuid,
		name,
		length,
		created_at,
		updated_at
		)
		VALUES ($1, $2, $3, $4, $5)
		`

	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs("", "", "", "", "").
		WillReturnResult(sqlmock.NewResult(1, 1)) // Using UUID

	albumRepo := NewPostgresRepository(dbx)
	err = albumRepo.Add(context.TODO(), album)

	assert.NotNil(t, err)
}

func TestAddFailDuplicate(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	album := &domain.Album{
		UUID:      newUUID,
		Name:      "St. Anger",
		Length:    75,
		CreatedAt: now,
		UpdatedAt: now,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	rows := sqlmock.NewRows([]string{
		"uuid",
		"name",
		"length",
		"created_at",
		"updated_at",
	}).
		AddRow("", "", "", "", "")

	queryDuplicate := "SELECT \\* FROM albums"
	mock.ExpectQuery(queryDuplicate).WillReturnRows(rows).WillReturnError(sql.ErrNoRows)

	query := `INSERT INTO 
	albums (uuid, name, length, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5)`

	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(newUUID, album.Name, album.Length, album.CreatedAt, album.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1)) // Using UUID

	albumRepo := NewPostgresRepository(dbx)
	_ = albumRepo.Add(context.TODO(), album)

}

func TestUpdate(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	album := &domain.Album{
		UUID:      newUUID,
		Name:      "St. Anger",
		Length:    75,
		UpdatedAt: now,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	query := `
		UPDATE albums
		SET
		name=$1,
		length=$2,
		updated_at=$3
		WHERE uuid=$4
	`

	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(album.Name, album.Length, album.UpdatedAt, album.UUID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	albumRepo := NewPostgresRepository(dbx)
	err = albumRepo.Update(context.TODO(), newUUID, album)

	assert.NoError(t, err)
}

func TestUpdateFail(t *testing.T) {
	newUUID := uuid.New()
	album := &domain.Album{}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	query := `
		UPDATE albums
		SET
		name=$1,
		length=$2,
		updated_at=$3
		WHERE uuid=$4
	`

	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs("", "", "", "", "", "").
		WillReturnResult(sqlmock.NewResult(1, 1))

	albumRepo := NewPostgresRepository(dbx)
	err = albumRepo.Update(context.TODO(), newUUID, album)

	assert.NotNil(t, err)
}

func TestUpdateRowsAffected(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	album := &domain.Album{
		UUID:      newUUID,
		Name:      "St. Anger",
		Length:    75,
		UpdatedAt: now,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	query := `
		UPDATE albums
		SET
		name=$1,
		length=$2,
		updated_at=$3
		WHERE uuid=$4
	`

	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(
		album.Name,
		album.Length,
		album.UpdatedAt,
		album.UUID,
	).WillReturnResult(sqlmock.NewResult(1, 0))

	albumRepo := NewPostgresRepository(dbx)
	err = albumRepo.Update(context.TODO(), newUUID, album)

	assert.NotNil(t, err)
}

func TestUpdateRowsAffectedFail(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	album := &domain.Album{
		UUID:      newUUID,
		Name:      "St. Anger",
		Length:    75,
		UpdatedAt: now,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	query := `
		UPDATE albums
		SET
		name=$1,
		length=$2,
		updated_at=$3
		WHERE uuid=$4
	`

	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(
		album.Name,
		album.Length,
		album.UpdatedAt,
		album.UUID,
	).WillReturnResult(sqlmock.NewErrorResult(sql.ErrNoRows))

	albumRepo := NewPostgresRepository(dbx)
	err = albumRepo.Update(context.TODO(), newUUID, album)

	assert.NotNil(t, err)
}

func TestDelete(t *testing.T) {
	newUUID := uuid.New()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	query := "DELETE FROM albums WHERE uuid=\\$1"

	mock.ExpectExec(query).
		WithArgs(newUUID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	albumRepo := NewPostgresRepository(dbx)
	err = albumRepo.Delete(context.TODO(), newUUID)

	assert.NoError(t, err)
}

func TestDeleteFailure(t *testing.T) {
	newUUID := uuid.New()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	query := "DELETE FROM albums WHERE uuid=\\$1"

	mock.ExpectExec(query).
		WithArgs(0).
		WillReturnResult(sqlmock.NewResult(1, 1))

	albumRepo := NewPostgresRepository(dbx)
	err = albumRepo.Delete(context.TODO(), newUUID)

	assert.NotNil(t, err)
}

func TestDeleteRowsAffected(t *testing.T) {
	newUUID := uuid.New()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	query := "DELETE FROM albums WHERE uuid=\\$1"

	mock.ExpectExec(query).
		WithArgs(newUUID).
		WillReturnResult(sqlmock.NewResult(1, 0))

	albumRepo := NewPostgresRepository(dbx)
	err = albumRepo.Delete(context.TODO(), newUUID)

	assert.NotNil(t, err)
}

func TestDeleteRowsAffectedFail(t *testing.T) {
	newUUID := uuid.New()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	query := "DELETE FROM albums WHERE uuid=\\$1"

	mock.ExpectExec(query).
		WithArgs(newUUID).
		WillReturnResult(sqlmock.NewErrorResult(sql.ErrNoRows))

	albumRepo := NewPostgresRepository(dbx)
	err = albumRepo.Delete(context.TODO(), newUUID)

	assert.NotNil(t, err)
}
