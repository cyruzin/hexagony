package usecase

import (
	"context"
	"errors"
	authDomain "hexagony/internal/auth/domain"
	usersDomain "hexagony/internal/users/domain"
	"hexagony/pkg/crypto"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type authUseCase struct {
	authRepo authDomain.AuthRepository
}

func NewAuthUsecase(auth authDomain.AuthRepository) authDomain.AuthUseCase {
	return &authUseCase{
		authRepo: auth,
	}
}

func (a *authUseCase) Authenticate(ctx context.Context, email, password string) (*authDomain.AuthToken, error) {
	user, err := a.authRepo.Authenticate(ctx, email)
	if err != nil {
		return nil, err
	}

	bcrypt := crypto.New()

	if match := bcrypt.CheckPasswordHash(password, user.Password); !match {
		return nil, errors.New("authentication failed")
	}

	customClaims := &usersDomain.User{
		UUID:  user.UUID,
		Name:  user.Name,
		Email: user.Email,
	}

	tokenExpiration := time.Now().Add(1 * time.Hour) // 1 hour from now

	token, err := a.generateToken("user", customClaims, tokenExpiration)
	if err != nil {
		return nil, err
	}

	authToken := authDomain.AuthToken{Token: token}

	return &authToken, nil
}

func (a *authUseCase) generateToken(
	claimKey string,
	claimValue *usersDomain.User,
	expiration time.Time,
) (string, error) {
	if claimKey == "" || claimValue == nil {
		return "", authDomain.ErrEmptyClaim
	}

	signingKey := []byte(os.Getenv("JWT_SECRET"))

	claims := struct {
		jwt.RegisteredClaims
		UUID  uuid.UUID `json:"id"`
		Name  string    `json:"name"`
		Email string    `json:"email"`
	}{
		jwt.RegisteredClaims{
			Issuer:    "Hexagony",
			Subject:   "https://github.com/cyruzin/hexagony",
			Audience:  jwt.ClaimStrings{"Clean Architecture"},
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
		claimValue.UUID,
		claimValue.Name,
		claimValue.Email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	payload, err := token.SignedString(signingKey)
	if err != nil {
		return "", authDomain.ErrSign
	}

	return payload, nil
}
