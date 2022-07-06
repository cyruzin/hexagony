package mariadb

import (
	"context"
	"database/sql"
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
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &user, nil
}
