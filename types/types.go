package types

import (
	"github.com/golang-jwt/jwt/v5"
)

type (
	JwtClaims struct {
		Id uint
		jwt.RegisteredClaims
	}
)
