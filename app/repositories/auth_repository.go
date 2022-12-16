package postgres

import (
	"context"
	"database/sql"
	"errors"
	"hexagony/app/domain"
	"hexagony/app/repositories/queries"

	"github.com/jmoiron/sqlx"
)

type authRepository struct {
	Conn *sqlx.DB
}

func NewAuthRepository(Conn *sqlx.DB) domain.AuthRepository {
	return &authRepository{Conn}
}

func (p *authRepository) Authenticate(ctx context.Context, email string) (*domain.Users, error) {
	var user domain.Users

	err := p.Conn.GetContext(ctx, &user, queries.SqlAuthGetUser, email)
	noRows := errors.Is(err, sql.ErrNoRows)

	if noRows {
		return nil, domain.ErrAuthUserNotFound
	}

	if err != nil {
		return nil, domain.ErrAuth
	}

	return &user, nil
}
