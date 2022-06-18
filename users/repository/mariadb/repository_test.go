package mariadb

import (
	"context"
	"database/sql"
	"hexagony/users/domain"
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

	mockUsers := []domain.User{
		{
			UUID:      uuid.New(),
			Name:      "Cyro Dubeux",
			Email:     "xorycx@gmail.com",
			Password:  "12345678",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			UUID:      uuid.New(),
			Name:      "John Doe",
			Email:     "john@doe.com",
			Password:  "12345678",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{
		"uuid",
		"name",
		"email",
		"password",
		"created_at",
		"updated_at",
	}).AddRow(
		mockUsers[0].UUID,
		mockUsers[0].Name,
		mockUsers[0].Email,
		mockUsers[0].Password,
		mockUsers[0].CreatedAt,
		mockUsers[0].UpdatedAt,
	).AddRow(
		mockUsers[1].UUID,
		mockUsers[1].Name,
		mockUsers[1].Email,
		mockUsers[1].Password,
		mockUsers[1].CreatedAt,
		mockUsers[1].UpdatedAt,
	)

	query := "SELECT \\* FROM users"

	mock.ExpectQuery(query).WillReturnRows(rows)

	userRepo := NewMariaDBRepository(dbx)
	userList, err := userRepo.FindAll(context.TODO())

	assert.NoError(t, err)
	assert.Len(t, userList, 2)
	assert.Equal(t, userList[0].Name, "Cyro Dubeux")
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
		"email",
		"password",
		"created_at",
		"updated_at",
	}).
		AddRow("", "", "", "", "", "")

	query := "SELECT \\* FROM users"
	mock.ExpectQuery(query).WillReturnRows(rows)

	userRepo := NewMariaDBRepository(dbx)
	_, err = userRepo.FindAll(context.TODO())

	assert.NotNil(t, err)
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
		"email",
		"password",
		"created_at",
		"updated_at",
	}).
		AddRow(newUUID, "Cyro Dubeux", "xorycx@gmail.com", "12345678", time.Now(), time.Now())

	query := "SELECT \\* FROM users WHERE uuid=\\?"
	mock.ExpectQuery(query).WillReturnRows(rows)

	userRepo := NewMariaDBRepository(dbx)
	currentUser, err := userRepo.FindByID(context.TODO(), newUUID)

	assert.NoError(t, err)
	assert.NotNil(t, currentUser)
	assert.Equal(t, "Cyro Dubeux", currentUser.Name)
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
		"email",
		"password",
		"created_at",
		"updated_at",
	}).
		AddRow("", "", "", "", "", "")

	query := "SELECT \\* FROM users WHERE uuid=\\?"
	mock.ExpectQuery(query).WillReturnRows(rows)

	userRepo := NewMariaDBRepository(dbx)
	_, err = userRepo.FindByID(ctx, newUUID)

	assert.NotNil(t, err)
}

func TestAdd(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	user := &domain.User{
		UUID:      newUUID,
		Name:      "Cyro Dubeux",
		Email:     "xorycx@gmail.com",
		Password:  "12345678",
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
	users (uuid, name, email, password, created_at, updated_at) 
	VALUES (?, ?, ?, ?, ?, ?)`

	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(newUUID, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1)) // Using UUID

	userRepo := NewMariaDBRepository(dbx)
	err = userRepo.Add(context.TODO(), user)

	assert.NoError(t, err)
}

func TestStoreFail(t *testing.T) {
	user := &domain.User{}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	query := `
	  INSERT INTO users (
		uuid,
		name,
		email,
		password,
		created_at,
		updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?)
		`

	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs("", "", "", "", "").
		WillReturnResult(sqlmock.NewResult(1, 1)) // Using UUID

	userRepo := NewMariaDBRepository(dbx)
	err = userRepo.Add(context.TODO(), user)

	assert.NotNil(t, err)
}

func TestUpdate(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	user := &domain.UserUpdate{
		UUID:      newUUID,
		Name:      "Cyro Dubeux",
		Email:     "xorycx@gmail.com",
		Password:  "12345678",
		UpdatedAt: now,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	query := `
		UPDATE users
		SET
		name=?,
		email=?,
		password=?,
		updated_at=?
		WHERE uuid=?
	`

	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(user.Name, user.Email, user.Password, user.UpdatedAt, user.UUID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	userRepo := NewMariaDBRepository(dbx)
	err = userRepo.Update(context.TODO(), newUUID, user)

	assert.NoError(t, err)
}

func TestUpdateFail(t *testing.T) {
	newUUID := uuid.New()
	user := &domain.UserUpdate{}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	query := `
		UPDATE users
		SET
		name=?,
		email=?,
		password=?,
		updated_at=?
		WHERE uuid=?
	`

	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs("", "", "", "", "", "").
		WillReturnResult(sqlmock.NewResult(1, 1))

	userRepo := NewMariaDBRepository(dbx)
	err = userRepo.Update(context.TODO(), newUUID, user)

	assert.NotNil(t, err)
}

func TestUpdateRowsAffected(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	user := &domain.UserUpdate{
		UUID:      newUUID,
		Name:      "Cyro Dubeux",
		Email:     "xorycx@gmail.com",
		Password:  "12345678",
		UpdatedAt: now,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	query := `
		UPDATE users
		SET
		name=?,
		email=?,
		password=?,
		updated_at=?
		WHERE uuid=?
	`

	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(
		user.Name,
		user.Email,
		user.Password,
		user.UpdatedAt,
		user.UUID,
	).WillReturnResult(sqlmock.NewResult(1, 0))

	userRepo := NewMariaDBRepository(dbx)
	err = userRepo.Update(context.TODO(), newUUID, user)

	assert.NotNil(t, err)
}

func TestUpdateRowsAffectedFail(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	user := &domain.UserUpdate{
		UUID:      newUUID,
		Name:      "Cyro Dubeux",
		Email:     "xorycx@gmail.com",
		Password:  "12345678",
		UpdatedAt: now,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	query := `
		UPDATE users
		SET
		name=?,
		email=?,
		password=?,
		updated_at=?
		WHERE uuid=?
	`

	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(
		user.Name,
		user.Email,
		user.Password,
		user.UpdatedAt,
		user.UUID,
	).WillReturnResult(sqlmock.NewErrorResult(sql.ErrNoRows))

	userRepo := NewMariaDBRepository(dbx)
	err = userRepo.Update(context.TODO(), newUUID, user)

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

	query := "DELETE FROM users WHERE uuid=\\?"

	mock.ExpectExec(query).
		WithArgs(newUUID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	userRepo := NewMariaDBRepository(dbx)
	err = userRepo.Delete(context.TODO(), newUUID)

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

	query := "DELETE FROM users WHERE uuid=\\?"

	mock.ExpectExec(query).
		WithArgs(0).
		WillReturnResult(sqlmock.NewResult(1, 1))

	userRepo := NewMariaDBRepository(dbx)
	err = userRepo.Delete(context.TODO(), newUUID)

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

	query := "DELETE FROM users WHERE uuid=\\?"

	mock.ExpectExec(query).
		WithArgs(newUUID).
		WillReturnResult(sqlmock.NewResult(1, 0))

	userRepo := NewMariaDBRepository(dbx)
	err = userRepo.Delete(context.TODO(), newUUID)

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

	query := "DELETE FROM users WHERE uuid=\\?"

	mock.ExpectExec(query).
		WithArgs(newUUID).
		WillReturnResult(sqlmock.NewErrorResult(sql.ErrNoRows))

	userRepo := NewMariaDBRepository(dbx)
	err = userRepo.Delete(context.TODO(), newUUID)

	assert.NotNil(t, err)
}
