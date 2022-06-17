package crypto

import "golang.org/x/crypto/bcrypt"

type Crypto interface {
	HashPassword(password string, cost int) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type bcryptHash struct{}

// HashPassword encrypts a given password using bcrypt algorithm.
func (b bcryptHash) HashPassword(password string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// CheckPasswordHash checks if the given passwords matches.
func (b bcryptHash) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

func New() Crypto {
	return &bcryptHash{}
}
