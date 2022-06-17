package mariadb

import (
	"context"
	"database/sql"
	authDomain "hexagony/auth/domain"
	userDomain "hexagony/users/domain"

	"github.com/jmoiron/sqlx"
)

type mariadbRepository struct {
	Conn *sqlx.DB
}

// NewAuthRepository will create an object that represent
// the auth.Repository interface.
func NewAuthRepository(Conn *sqlx.DB) authDomain.AuthRepository {
	return &mariadbRepository{Conn}
}

func (p *mariadbRepository) Authenticate(ctx context.Context, email string, password string) (*userDomain.User, error) {
	var user userDomain.User

	err := p.Conn.GetContext(ctx, &user, sqlGetUser, email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &user, nil
}
