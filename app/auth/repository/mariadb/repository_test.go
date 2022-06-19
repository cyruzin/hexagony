package mariadb

import (
	"context"
	domainUsers "hexagony/app/users/domain"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	mockUser := domainUsers.User{
		UUID:      uuid.New(),
		Name:      "Cyro Dubeux",
		Email:     "xorycx@gmail.com",
		Password:  "12345678",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	row := sqlmock.NewRows([]string{
		"uuid",
		"name",
		"email",
		"password",
		"created_at",
		"updated_at",
	}).AddRow(
		mockUser.UUID,
		mockUser.Name,
		mockUser.Email,
		mockUser.Password,
		mockUser.CreatedAt,
		mockUser.UpdatedAt,
	)

	query := "SELECT \\* from users WHERE email = \\?"

	mock.ExpectQuery(query).WillReturnRows(row)

	authRepo := NewMariaDBRepository(dbx)
	user, err := authRepo.Authenticate(context.TODO(), "xorycx@gmail.com")

	assert.NoError(t, err)
	assert.Equal(t, "Cyro Dubeux", user.Name)
}

func TestAuthenticateFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	defer db.Close()

	dbx := sqlx.NewDb(db, "sqlmock")

	row := sqlmock.NewRows([]string{
		"uuid",
		"name",
		"email",
		"password",
		"created_at",
		"updated_at",
	}).AddRow("", "", "", "", "", "")

	query := "SELECT \\* from users WHERE email = \\?"

	mock.ExpectQuery(query).WillReturnRows(row)

	authRepo := NewMariaDBRepository(dbx)
	user, err := authRepo.Authenticate(context.TODO(), "xorycx@gmail.com")

	assert.Nil(t, user)
	assert.Error(t, err)
}
