package keycloakclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type IntrospectTokenResult struct {
	Exp    int      `json:"exp"`
	Iat    int      `json:"iat"`
	Aud    []string `json:"aud,omitempty"` // может быть как строкой, так и массивом
	Active bool     `json:"active"`
}

// UnmarshalJSON реализует пользовательское декодирование для IntrospectTokenResult.
func (r *IntrospectTokenResult) UnmarshalJSON(data []byte) error {
	// Промежуточная структура с Aud типа interface{}
	type Alias struct {
		Exp    int  `json:"exp"`
		Iat    int  `json:"iat"`
		Aud    any  `json:"aud"`
		Active bool `json:"active"`
	}

	var alias Alias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	// Копируем простые поля
	r.Exp = alias.Exp
	r.Iat = alias.Iat
	r.Active = alias.Active

	// Преобразуем Aud в []string в зависимости от типа
	switch v := alias.Aud.(type) {
	case string:
		r.Aud = []string{v}
	case []string:
		r.Aud = v
	case []any:
		r.Aud = make([]string, len(v))
		for i, item := range v {
			if s, ok := item.(string); ok {
				r.Aud[i] = s
			}
		}
	case nil:
		r.Aud = nil
	default:
		r.Aud = []string{}
	}

	return nil
}

// IntrospectToken implements
// https://www.keycloak.org/docs/latest/authorization_services/index.html#obtaining-information-about-an-rpt
func (c *Client) IntrospectToken(ctx context.Context, token string) (*IntrospectTokenResult, error) {
	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token/introspect", c.basePath, c.realm)

	resp, err := c.auth(ctx).
		SetFormData(map[string]string{
			"token":           token,
			"token_type_hint": "requesting_party_token",
		}).
		Post(url)
	if err != nil {
		return nil, fmt.Errorf("send request to keycloak: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("errored keycloak response: %v", resp.Status())
	}

	var result IntrospectTokenResult
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return nil, fmt.Errorf("unmarshal keycloak response: %v", err)
	}

	return &result, nil
}

func (c *Client) auth(ctx context.Context) *resty.Request {
	// Используем client_id и client_secret для Basic Authentication
	// согласно OAuth 2.0 Client Credentials Grant
	return c.cli.R().
		SetContext(ctx).
		SetBasicAuth(c.clientID, c.clientSecret)
}
