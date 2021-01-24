package jwt

import (
	"github.com/dgrijalva/jwt-go"
)

type JwtCustomClaims struct {
	UserID string `json:"userID"`
	jwt.StandardClaims
}
