package postgres

import (
	"context"
	"database/sql"
	"hexagony/app/domain"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestUsersFindAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	mockUsers := []domain.UsersList{
		{
			UUID:      uuid.New(),
			Name:      "Cyro Dubeux",
			Email:     "xorycx@gmail.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			UUID:      uuid.New(),
			Name:      "John Doe",
			Email:     "john@doe.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{
		"uuid",
		"name",
		"email",
		"created_at",
		"updated_at",
	}).AddRow(
		mockUsers[0].UUID,
		mockUsers[0].Name,
		mockUsers[0].Email,
		mockUsers[0].CreatedAt,
		mockUsers[0].UpdatedAt,
	).AddRow(
		mockUsers[1].UUID,
		mockUsers[1].Name,
		mockUsers[1].Email,
		mockUsers[1].CreatedAt,
		mockUsers[1].UpdatedAt,
	)

	query := "SELECT uuid,name,email,created_at,updated_at FROM users ORDER BY updated_at DESC LIMIT 10"

	mock.ExpectQuery(query).WillReturnRows(rows)

	userRepo := NewUsersRepository(dbx)
	userList, err := userRepo.FindAll(context.TODO())

	assert.NoError(t, err)
	assert.Len(t, userList, 2)
	assert.Equal(t, userList[0].Name, "Cyro Dubeux")
}

func TestUsersFindAllFail(t *testing.T) {
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
		"created_at",
		"updated_at",
	}).
		AddRow("", "", "", "", "")

	query := "SELECT uuid,name,email,created_at,updated_at FROM users ORDER BY updated_at DESC LIMIT 10"
	mock.ExpectQuery(query).WillReturnRows(rows)

	userRepo := NewUsersRepository(dbx)
	_, err = userRepo.FindAll(context.TODO())

	assert.NotNil(t, err)
}

func TestUsersFindAllFailResource(t *testing.T) {
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
		"created_at",
		"updated_at",
	}).
		AddRow("", "", "", "", "")

	query := "SELECT uuid,name,email,created_at,updated_at FROM users ORDER BY updated_at DESC LIMIT 10"
	mock.ExpectQuery(query).WillReturnRows(rows).WillReturnError(sql.ErrNoRows)

	userRepo := NewUsersRepository(dbx)
	users, _ := userRepo.FindAll(context.TODO())

	assert.Nil(t, users)
}

func TestUsersFindByID(t *testing.T) {
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
		"created_at",
		"updated_at",
	}).
		AddRow(newUUID, "Cyro Dubeux", "xorycx@gmail.com", time.Now(), time.Now())

	query := "SELECT uuid,name,email,created_at,updated_at FROM users WHERE uuid=$1"
	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

	userRepo := NewUsersRepository(dbx)
	currentUser, err := userRepo.FindByID(context.TODO(), newUUID)

	assert.NoError(t, err)
	assert.NotNil(t, currentUser)
	assert.Equal(t, "Cyro Dubeux", currentUser.Name)
}

func TestUsersGetByIDFail(t *testing.T) {
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
		"created_at",
		"updated_at",
	}).
		AddRow("", "", "", "", "")

	query := "SELECT uuid,name,email,created_at,updated_at FROM users WHERE uuid=$1"
	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows)

	userRepo := NewUsersRepository(dbx)
	_, err = userRepo.FindByID(ctx, newUUID)

	assert.NotNil(t, err)
}

func TestUsersGetByIDFailResource(t *testing.T) {
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
		"created_at",
		"updated_at",
	}).
		AddRow("", "", "", "", "")

	query := "SELECT uuid,name,email,created_at,updated_at FROM users WHERE uuid=$1"
	mock.ExpectQuery(regexp.QuoteMeta(query)).WillReturnRows(rows).WillReturnError(sql.ErrNoRows)

	userRepo := NewUsersRepository(dbx)
	_, err = userRepo.FindByID(ctx, newUUID)

	assert.NotNil(t, err)
}

func TestUsersAdd(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	user := &domain.Users{
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

	rows := sqlmock.NewRows([]string{"email"}).AddRow("")

	queryCheckDuplicate := "SELECT email FROM users WHERE email=$1"
	mock.ExpectQuery(regexp.QuoteMeta(queryCheckDuplicate)).WillReturnRows(rows)

	query := `INSERT INTO 
	users (uuid, name, email, password, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6)`

	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(newUUID, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1)) // Using UUID

	userRepo := NewUsersRepository(dbx)
	err = userRepo.Add(context.TODO(), user)

	assert.NoError(t, err)
}

func TestUsersAddFail(t *testing.T) {
	user := &domain.Users{}

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
		VALUES ($1, $2, $3, $4, $5, $6)
		`

	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs("", "", "", "", "").
		WillReturnResult(sqlmock.NewResult(1, 1)) // Using UUID

	userRepo := NewUsersRepository(dbx)
	err = userRepo.Add(context.TODO(), user)

	assert.NotNil(t, err)
}

func TestUsersAddFailEmail(t *testing.T) {
	user := &domain.Users{Name: "xorycx@gmail.com"}
	newUUID := uuid.New()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	queryCheckDuplicate := "SELECT email FROM users WHERE email=$1"
	mock.ExpectQuery(regexp.QuoteMeta(queryCheckDuplicate)).
		WillReturnError(domain.ErrUsersDuplicateEmail)

	query := `
	  INSERT INTO users (
		uuid,
		name,
		email,
		password,
		created_at,
		updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		`

	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs(newUUID, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1)) // Using UUID

	userRepo := NewUsersRepository(dbx)
	err = userRepo.Add(context.TODO(), user)

	assert.NotNil(t, err)
}

func TestUsersUpdate(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	user := &domain.Users{
		UUID:      newUUID,
		Name:      "Cyro Dubeux",
		Email:     "xorycx@gmail.com",
		UpdatedAt: now,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	query := "UPDATE users SET name=\\$1, email=\\$2, updated_at=\\$3 WHERE uuid=\\$4"

	mock.ExpectExec(query).
		WithArgs(user.Name, user.Email, user.UpdatedAt, user.UUID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	userRepo := NewUsersRepository(dbx)
	err = userRepo.Update(context.TODO(), newUUID, user)

	assert.NoError(t, err)
}

func TestUsersUpdateFail(t *testing.T) {
	newUUID := uuid.New()
	user := &domain.Users{}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	query := `
		UPDATE users
		SET
		name=$1,
		email=$2,
		updated_at=$3 
		WHERE uuid=$4
	`

	mock.ExpectExec(regexp.QuoteMeta(query)).
		WithArgs("", "", "", "", "").
		WillReturnResult(sqlmock.NewResult(1, 1))

	userRepo := NewUsersRepository(dbx)
	err = userRepo.Update(context.TODO(), newUUID, user)

	assert.NotNil(t, err)
}

func TestUsersUpdateRowsAffected(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	user := &domain.Users{
		UUID:      newUUID,
		Name:      "Cyro Dubeux",
		Email:     "xorycx@gmail.com",
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
		name=$1,
		email=$2,
		updated_at=$3
		WHERE uuid=$4
	`

	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(
		user.Name,
		user.Email,
		user.UpdatedAt,
		user.UUID,
	).WillReturnResult(sqlmock.NewResult(1, 0))

	userRepo := NewUsersRepository(dbx)
	err = userRepo.Update(context.TODO(), newUUID, user)

	assert.NotNil(t, err)
}

func TestUsersUpdateRowsAffectedFail(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	user := &domain.Users{
		UUID:      newUUID,
		Name:      "Cyro Dubeux",
		Email:     "xorycx@gmail.com",
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
		name=$1,
		email=$2,
		updated_at=$3
		WHERE uuid=$4
	`

	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(
		user.Name,
		user.Email,
		user.UpdatedAt,
		user.UUID,
	).WillReturnResult(sqlmock.NewErrorResult(sql.ErrNoRows))

	userRepo := NewUsersRepository(dbx)
	err = userRepo.Update(context.TODO(), newUUID, user)

	assert.NotNil(t, err)
}

func TestUsersDelete(t *testing.T) {
	newUUID := uuid.New()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	query := "DELETE FROM users WHERE uuid=\\$1"

	mock.ExpectExec(query).
		WithArgs(newUUID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	userRepo := NewUsersRepository(dbx)
	err = userRepo.Delete(context.TODO(), newUUID)

	assert.NoError(t, err)
}

func TestUsersDeleteFailure(t *testing.T) {
	newUUID := uuid.New()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	query := "DELETE FROM users WHERE uuid=\\$1"

	mock.ExpectExec(query).
		WithArgs(0).
		WillReturnResult(sqlmock.NewResult(1, 1))

	userRepo := NewUsersRepository(dbx)
	err = userRepo.Delete(context.TODO(), newUUID)

	assert.NotNil(t, err)
}

func TestUsersDeleteRowsAffected(t *testing.T) {
	newUUID := uuid.New()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	query := "DELETE FROM users WHERE uuid=\\$1"

	mock.ExpectExec(query).
		WithArgs(newUUID).
		WillReturnResult(sqlmock.NewResult(1, 0))

	userRepo := NewUsersRepository(dbx)
	err = userRepo.Delete(context.TODO(), newUUID)

	assert.NotNil(t, err)
}

func TestUsersDeleteRowsAffectedFail(t *testing.T) {
	newUUID := uuid.New()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	query := "DELETE FROM users WHERE uuid=\\$1"

	mock.ExpectExec(query).
		WithArgs(newUUID).
		WillReturnResult(sqlmock.NewErrorResult(sql.ErrNoRows))

	userRepo := NewUsersRepository(dbx)
	err = userRepo.Delete(context.TODO(), newUUID)

	assert.NotNil(t, err)
}

func TestUsersCheckDuplicate(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	user := &domain.Users{
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

	rows := sqlmock.NewRows([]string{"email"}).AddRow("xorycx@gmail.com")

	queryCheckDuplicate := "SELECT email FROM users WHERE email=$1"
	mock.ExpectQuery(regexp.QuoteMeta(queryCheckDuplicate)).
		WillReturnRows(rows).
		WillReturnError(sql.ErrNoRows)

	query := `INSERT INTO 
	users (uuid, name, email, password, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6)`

	mock.ExpectQuery(query).WillReturnRows(rows)

	userRepo := NewUsersRepository(dbx)
	err = userRepo.Add(context.TODO(), user)

	assert.NotNil(t, err)
}

func TestUsersCheckDuplicateFail(t *testing.T) {
	now := time.Now()
	newUUID := uuid.New()
	user := &domain.Users{
		UUID:      newUUID,
		Name:      "John Doe",
		Email:     "john@doe.com",
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

	rows := sqlmock.NewRows([]string{"email"}).AddRow("xorycx@gmail.com")

	queryCheckDuplicate := "SELECT email FROM users WHERE email=$1"
	mock.ExpectQuery(regexp.QuoteMeta(queryCheckDuplicate)).
		WillReturnRows(rows)

	query := `INSERT INTO 
	users (uuid, name, email, password, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6)`

	mock.ExpectQuery(query).WillReturnRows(rows)

	userRepo := NewUsersRepository(dbx)
	err = userRepo.Add(context.TODO(), user)

	assert.NotNil(t, err)
}
