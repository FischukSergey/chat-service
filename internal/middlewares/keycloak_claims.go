package middlewares

import (
	"errors"

	"github.com/golang-jwt/jwt"

	"github.com/FischukSergey/chat-service/internal/types"
)

var (
	ErrNoAllowedResources = errors.New("no allowed resources")
	ErrSubjectNotDefined  = errors.New(`"sub" is not defined`)
)

type claims struct {
	jwt.StandardClaims
	// FIXME: добавь поля, которые нужны для проверки токена
	Subject        string              `json:"sub,omitempty"`
	RealmAccess    map[string][]string `json:"realm_access,omitempty"`
	ResourceAccess map[string]struct {
		Roles []string `json:"roles,omitempty"`
	} `json:"resource_access,omitempty"`
}

// Valid returns errors:
// - from StandardClaims validation;
// - ErrNoAllowedResources, if claims doesn't contain `resource_access` map or it's empty;
// - ErrSubjectNotDefined, if claims doesn't contain `sub` field or subject is zero UUID.
func (c claims) Valid() error {
	// FIXME: реализуй меня

	if err := c.StandardClaims.Valid(); err != nil {
		return err
	}

	if c.Subject == "" {
		return ErrSubjectNotDefined
	}

	if len(c.ResourceAccess) == 0 {
		return ErrNoAllowedResources
	}

	return nil
}

func (c claims) UserID() types.UserID {
	return types.MustParse[types.UserID](c.Subject)
}

// HasResourceRole проверяет наличие указанной роли для указанного ресурса
func (c claims) HasResourceRole(resource, role string) bool {
	resourceRoles, exists := c.ResourceAccess[resource]
	if !exists {
		return false
	}

	for _, r := range resourceRoles.Roles {
		if r == role {
			return true
		}
	}
	return false
}
