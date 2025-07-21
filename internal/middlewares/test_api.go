package middlewares

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"github.com/FischukSergey/chat-service/internal/types"
)

func SetToken(c echo.Context, uid types.UserID) {
	// Фикс: В контекст по ключу tokenCtxKey необходимо положить jwt.Token с клэймсами:
	// Фикс: - которые всегда валидные
	// Фикс: - из которых можно достать uid
	c.Set(tokenCtxKey, jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		StandardClaims: jwt.StandardClaims{
			Subject: uid.String(),
		},
		RealmAccess: map[string][]string{
			"roles": {"user"},
		},
		ResourceAccess: map[string]struct {
			Roles []string `json:"roles,omitempty"`
		}{
			"chat-service": {
				Roles: []string{"user"},
			},
		},
	}))
}
