package middlewares

import (
	"context"
	"errors"

	"github.com/labstack/echo/v4"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4/middleware"

	keycloakclient "github.com/FischukSergey/chat-service/internal/clients/keycloak"
	"github.com/FischukSergey/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/introspector_mock.gen.go -package=middlewaresmocks Introspector

const tokenCtxKey = "user-token"

var ErrNoRequiredResourceRole = errors.New("no required resource role")
var ErrTokenNotActive = errors.New("token is not active")

type Introspector interface {
	IntrospectToken(ctx context.Context, token string) (*keycloakclient.IntrospectTokenResult, error)
}

// NewKeycloakTokenAuth returns a middleware that implements "active" authentication:
// each request is verified by the Keycloak server.
func NewKeycloakTokenAuth(introspector Introspector, resource, role string) echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup:  "header:Authorization", // Вставил
		AuthScheme: "Bearer",
		Validator: func(tokenStr string, eCtx echo.Context) (bool, error) {
			// 1. FIXME: интроспектим токен
			// 2. FIXME: проверяем, что он Active
			// 3. FIXME: парсим токен, используя наши claims (без проверки подписи, это уже сделал Keycloak)
			// 4. FIXME: проверяем, что клеймы валидные
			// 5. FIXME: проверяем, что среди них есть нужная роль для нужного ресурса
			// 6. FIXME: сохраняем токен в контекст запроса

			// 1. интроспектим токен
			introspectResult, err := introspector.IntrospectToken(eCtx.Request().Context(), tokenStr)
			if err != nil {
				return false, err
			}
			// 2. проверяем, что он Active
			if !introspectResult.Active {
				return false, ErrTokenNotActive
			}
			// 3. парсим токен, используя наши claims (без проверки подписи, это уже сделал Keycloak)
			claims := &claims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
				return nil, errors.New("не проверяем подпись")
			})

			// Проверяем наличие Subject и что он не Zero UUID
			if claims.Subject == "" || claims.Subject == "00000000-0000-0000-0000-000000000000" {
				return false, ErrSubjectNotDefined
			}

			// 4. проверяем, что клеймы валидные
			if err := claims.Valid(); err != nil {
				return false, err
			}
			// 5. проверяем, что среди них есть нужная роль для нужного ресурса
			if !claims.HasResourceRole(resource, role) {
				return false, ErrNoRequiredResourceRole
			}
			// 6. сохраняем токен в контекст запроса
			if token == nil {
				token = &jwt.Token{Claims: claims}
			}
			eCtx.Set(tokenCtxKey, token)
			return true, nil
		},
	})
}

// возвращает userID из контекста запроса
func MustUserID(eCtx echo.Context) types.UserID {
	uid, ok := userID(eCtx)
	if !ok {
		panic("no user token in request context")
	}
	return uid
}

// возвращает userID из контекста запроса
func userID(eCtx echo.Context) (types.UserID, bool) {
	t := eCtx.Get(tokenCtxKey)
	if t == nil {
		return types.UserIDNil, false
	}

	tt, ok := t.(*jwt.Token)
	if !ok {
		return types.UserIDNil, false
	}

	userIDProvider, ok := tt.Claims.(interface{ UserID() types.UserID })
	if !ok {
		return types.UserIDNil, false
	}
	return userIDProvider.UserID(), true
}
