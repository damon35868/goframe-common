package middleware

import (
	"net/http"

	"github.com/damon35868/goframe-common/vars"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/golang-jwt/jwt/v5"
)

func commonAuth(r *ghttp.Request, jwtKey string) {
	var tokenString = r.Header.Get("Authorization")
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(jwtKey), nil
	})
	if err != nil || token == nil || !token.Valid {
		r.Response.WriteStatus(http.StatusForbidden)
		r.Response.ClearBuffer()

		msg := "token无效或已过期"
		if err != nil {
			msg = err.Error()
		}
		r.Response.WriteJson(g.Map{
			"code":    http.StatusForbidden,
			"message": msg,
		})
		r.Exit()
	}
}

func ClientAuth(r *ghttp.Request) {
	commonAuth(r, g.Cfg().MustGet(r.GetCtx(), vars.ClientJwtKey).String())
	r.Middleware.Next()
}

func AdminAuth(r *ghttp.Request) {
	commonAuth(r, g.Cfg().MustGet(r.GetCtx(), vars.AdminJwtKey).String())
	r.Middleware.Next()
}
