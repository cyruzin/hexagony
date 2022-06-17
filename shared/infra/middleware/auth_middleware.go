package middleware

import (
	"errors"
	"fmt"
	"hexagony/lib/rest"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// AuthMiddleware checks if the request contains Bearer Token
// on the headers and if it is valid.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Capturing Authorizathion header.
		tokenHeader := r.Header.Get("Authorization")

		// Checking if the value is empty.
		if tokenHeader == "" {
			rest.DecodeError(w, r, errors.New("empty token"), http.StatusBadRequest)
			return
		}

		// Checking if the header contains Bearer string and if the token exists.
		if !strings.Contains(tokenHeader, "Bearer") || len(strings.Split(tokenHeader, "Bearer ")) == 1 {
			rest.DecodeError(w, r, errors.New("malformed token"), http.StatusBadRequest)
			return
		}

		// Capturing the token.
		jwtString := strings.Split(tokenHeader, "Bearer ")[1]

		// Parsing the token to verify its authenticity.
		token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("secret"), nil
		})

		// Returning parsing errors.
		if err != nil {
			rest.DecodeError(w, r, errors.New("unathorized"), http.StatusUnauthorized)
			return
		}

		// If the token is valid.
		if token.Valid {
			next.ServeHTTP(w, r)
		} else {
			rest.DecodeError(w, r, errors.New("invalid jwt token"), http.StatusUnauthorized)
			return
		}
	})
}
