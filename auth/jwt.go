package auth

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"strings"
)

// IdentityClaims identity JWT claims
type IdentityClaims struct {
	Username string  `json:"username"`
	Name     string  `json:"name"`
	LastName string  `json:"last_name"`
	Picture  *string `json:"picture"`
	Email    string  `json:"email"`
	Locale   string  `json:"locale"`
	Role     string  `json:"role"`
	jwt.StandardClaims
}

// ParseBearerJWT returns JWT claims from JWT bearer token
func ParseBearerJWT(bearer string) (*IdentityClaims, error) {
	spToken := strings.Split(bearer, " ")

	if len(spToken) < 1 {
		return nil, errors.New("invalid bearer token")
	}

	tokenStr := spToken[1]
	token, err := jwt.ParseWithClaims(tokenStr, &IdentityClaims{}, func(token *jwt.Token) (interface{}, error) {
		return token, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*IdentityClaims)
	if !ok || !token.Valid {
		return nil, errors.New("failed to map jwt claims")
	}

	return claims, nil
}
