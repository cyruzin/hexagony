package crypto

import "golang.org/x/crypto/bcrypt"

// Crypto interface to implement new algorithms.
type Crypto interface {
	HashPassword(password string, cost int) (string, error)
	CheckPasswordHash(password, hash string) bool
}

// bcryptHash is the type used in bcrypt.
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

// New creates a new instance to work with the bcrypt algorithm functions.
func New() Crypto {
	return &bcryptHash{}
}
