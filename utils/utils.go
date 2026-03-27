package utils

import (
	"context"

	"github.com/damon35868/goframe-common/types"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashedPassword(password string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hash)
}

func GetTokenUserId(ctx context.Context, jwtKey string) uint {
	tokenString := g.RequestFromCtx(ctx).Request.Header.Get("Authorization")
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}
	tokenClaims, _ := jwt.ParseWithClaims(tokenString, &types.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if claims, ok := tokenClaims.Claims.(*types.JwtClaims); ok && tokenClaims.Valid {
		return claims.Id
	}

	return 0
}

func GenToken(ctx context.Context, jwtKey string, id uint, registeredClaims jwt.RegisteredClaims) (token string, err error) {
	uc := &types.JwtClaims{
		Id:               id,
		RegisteredClaims: registeredClaims,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, uc)
	token, err = t.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}
	return
}
