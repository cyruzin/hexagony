package mariadb

import (
	"context"
	"database/sql"
	authDomain "hexagony/internal/auth/domain"
	userDomain "hexagony/internal/users/domain"

	"github.com/jmoiron/sqlx"
)

type mariadbRepository struct {
	Conn *sqlx.DB
}

func NewMariaDBRepository(Conn *sqlx.DB) authDomain.AuthRepository {
	return &mariadbRepository{Conn}
}

func (p *mariadbRepository) Authenticate(ctx context.Context, email string) (*userDomain.User, error) {
	var user userDomain.User

	err := p.Conn.GetContext(ctx, &user, sqlGetUser, email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &user, nil
}
