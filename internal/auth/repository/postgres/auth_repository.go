package postgres

import (
	"context"
	"database/sql"
	"errors"
	authDomain "hexagony/internal/auth/domain"
	userDomain "hexagony/internal/users/domain"

	"github.com/jmoiron/sqlx"
)

type postgresRepository struct {
	Conn *sqlx.DB
}

func NewPostgresRepository(Conn *sqlx.DB) authDomain.AuthRepository {
	return &postgresRepository{Conn}
}

func (p *postgresRepository) Authenticate(ctx context.Context, email string) (*userDomain.User, error) {
	var user userDomain.User

	err := p.Conn.GetContext(ctx, &user, sqlGetUser, email)
	noRows := errors.Is(err, sql.ErrNoRows)

	if noRows {
		return nil, authDomain.ErrUserNotFound
	}

	if err != nil {
		return nil, authDomain.ErrAuth
	}

	return &user, nil
}
