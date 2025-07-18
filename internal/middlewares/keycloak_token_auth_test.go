package middlewares_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	keycloakclient "github.com/FischukSergey/chat-service/internal/clients/keycloak"
	"github.com/FischukSergey/chat-service/internal/middlewares"
	middlewaresmocks "github.com/FischukSergey/chat-service/internal/middlewares/mocks"
	"github.com/FischukSergey/chat-service/internal/types"
)

const (
	requiredResource = "chat-ui-client"
	requiredRole     = "support-chat-client"
	bearerPrefix     = "Bearer "
)

func TestNewKeycloakTokenAuth(t *testing.T) {
	suite.Run(t, new(KeycloakTokenAuthSuite))
}

type KeycloakTokenAuthSuite struct {
	suite.Suite
	ctrl         *gomock.Controller
	introspector *middlewaresmocks.MockIntrospector
	authMdlwr    echo.MiddlewareFunc
	req          *http.Request
	resp         *httptest.ResponseRecorder
	ctx          echo.Context
}

func (s *KeycloakTokenAuthSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())

	s.introspector = middlewaresmocks.NewMockIntrospector(s.ctrl)
	s.authMdlwr = middlewares.NewKeycloakTokenAuth(s.introspector, requiredResource, requiredRole)

	s.req = httptest.NewRequest(http.MethodPost, "/getHistory",
		bytes.NewBufferString(`{"pageSize": 100, "cursor": ""}`))
	s.resp = httptest.NewRecorder()
	s.ctx = echo.New().NewContext(s.req, s.resp)
}

func (s *KeycloakTokenAuthSuite) TearDownTest() {
	s.ctrl.Finish()
}

// Positive.

func (s *KeycloakTokenAuthSuite) TestValidToken_AudString() {
	const token = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJIR1lJcHN1UXlsZFNJZTB1T0JaeEpuQjBkZlFuTWI5LUlFcmx6NHk5ek9BIn0.eyJleHAiOjI2NjcxOTk1ODAsImlhdCI6MTY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiNWNiNDBkYzAtYTI0OS00NzgzLWEzMDEtOWUxZjNjZjNlYTQxIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiY2hhdC11aS1jbGllbnQiLCJub25jZSI6ImJhMzdmZDVhLThjMzktNDgxNC1hZmNiLTk1MmExOGI3MjY3ZCIsInNlc3Npb25fc3RhdGUiOiJkODZkMTk4ZS1jMWM1LTRlZGQtODM1MC0zNjFlZTU4MTcxZjIiLCJhY3IiOiIwIiwiYWxsb3dlZC1vcmlnaW5zIjpbIiIsIioiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbIm9mZmxpbmVfYWNjZXNzIiwiZGVmYXVsdC1yb2xlcy1iYW5rIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInJlc291cmNlX2FjY2VzcyI6eyJjaGF0LXVpLWNsaWVudCI6eyJyb2xlcyI6WyJzdXBwb3J0LWNoYXQtY2xpZW50Il19LCJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50IiwibWFuYWdlLWFjY291bnQtbGlua3MiLCJ2aWV3LXByb2ZpbGUiXX19LCJzY29wZSI6Im9wZW5pZCBwcm9maWxlIGVtYWlsIiwic2lkIjoiZDg2ZDE5OGUtYzFjNS00ZWRkLTgzNTAtMzYxZWU1ODE3MWYyIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsInByZWZlcnJlZF91c2VybmFtZSI6ImJvbmQwMDciLCJnaXZlbl9uYW1lIjoiIiwiZmFtaWx5X25hbWUiOiIiLCJlbWFpbCI6ImJvbmQwMDdAdWsuY29tIn0.we-dont-check-signature" //nolint:lll
	s.req.Header.Add(echo.HeaderAuthorization, bearerPrefix+token)

	s.introspector.EXPECT().IntrospectToken(s.req.Context(), token).Return(&keycloakclient.IntrospectTokenResult{Active: true}, nil) //nolint:lll

	var uid types.UserID

	err := s.authMdlwr(func(c echo.Context) error {
		uid = middlewares.MustUserID(c)
		return nil
	})(s.ctx)
	s.Require().NoError(err)
	s.Equal("5cb40dc0-a249-4783-a301-9e1f3cf3ea41", uid.String())
}

func (s *KeycloakTokenAuthSuite) TestValidToken_AudList() {
	const token = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJIR1lJcHN1UXlsZFNJZTB1T0JaeEpuQjBkZlFuTWI5LUlFcmx6NHk5ek9BIn0.eyJleHAiOjI2NjcxOTk1ODAsImlhdCI6MTY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOlsiY2hhdC11aS1jbGllbnQiLCJhY2NvdW50Il0sInN1YiI6IjVjYjQwZGMwLWEyNDktNDc4My1hMzAxLTllMWYzY2YzZWE0MSIsInR5cCI6IkJlYXJlciIsImF6cCI6ImNoYXQtdWktY2xpZW50Iiwibm9uY2UiOiJiYTM3ZmQ1YS04YzM5LTQ4MTQtYWZjYi05NTJhMThiNzI2N2QiLCJzZXNzaW9uX3N0YXRlIjoiZDg2ZDE5OGUtYzFjNS00ZWRkLTgzNTAtMzYxZWU1ODE3MWYyIiwiYWNyIjoiMCIsImFsbG93ZWQtb3JpZ2lucyI6WyIiLCIqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJvZmZsaW5lX2FjY2VzcyIsImRlZmF1bHQtcm9sZXMtYmFuayIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiY2hhdC11aS1jbGllbnQiOnsicm9sZXMiOlsic3VwcG9ydC1jaGF0LWNsaWVudCJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsInNpZCI6ImQ4NmQxOThlLWMxYzUtNGVkZC04MzUwLTM2MWVlNTgxNzFmMiIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJib25kMDA3IiwiZ2l2ZW5fbmFtZSI6IiIsImZhbWlseV9uYW1lIjoiIiwiZW1haWwiOiJib25kMDA3QHVrLmNvbSJ9.we-dont-check-signature" //nolint:lll
	s.req.Header.Add(echo.HeaderAuthorization, bearerPrefix+token)

	s.introspector.EXPECT().IntrospectToken(s.req.Context(), token).
		Return(&keycloakclient.IntrospectTokenResult{Active: true}, nil)

	var uid types.UserID

	err := s.authMdlwr(func(c echo.Context) error {
		uid = middlewares.MustUserID(c)
		return nil
	})(s.ctx)
	s.Require().NoError(err)
	s.Equal("5cb40dc0-a249-4783-a301-9e1f3cf3ea41", uid.String())
}

// Negative.

func (s *KeycloakTokenAuthSuite) TestNoAuthorizationHeader() {
	const token = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJIR1lJcHN1UXlsZFNJZTB1T0JaeEpuQjBkZlFuTWI5LUlFcmx6NHk5ek9BIn0.eyJleHAiOjI2NjcxOTk1ODAsImlhdCI6MTY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOlsiY2hhdC11aS1jbGllbnQiLCJhY2NvdW50Il0sInN1YiI6IjVjYjQwZGMwLWEyNDktNDc4My1hMzAxLTllMWYzY2YzZWE0MSIsInR5cCI6IkJlYXJlciIsImF6cCI6ImNoYXQtdWktY2xpZW50Iiwibm9uY2UiOiJiYTM3ZmQ1YS04YzM5LTQ4MTQtYWZjYi05NTJhMThiNzI2N2QiLCJzZXNzaW9uX3N0YXRlIjoiZDg2ZDE5OGUtYzFjNS00ZWRkLTgzNTAtMzYxZWU1ODE3MWYyIiwiYWNyIjoiMCIsImFsbG93ZWQtb3JpZ2lucyI6WyIiLCIqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJvZmZsaW5lX2FjY2VzcyIsImRlZmF1bHQtcm9sZXMtYmFuayIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiY2hhdC11aS1jbGllbnQiOnsicm9sZXMiOlsic3VwcG9ydC1jaGF0LWNsaWVudCJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsInNpZCI6ImQ4NmQxOThlLWMxYzUtNGVkZC04MzUwLTM2MWVlNTgxNzFmMiIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJib25kMDA3IiwiZ2l2ZW5fbmFtZSI6IiIsImZhbWlseV9uYW1lIjoiIiwiZW1haWwiOiJib25kMDA3QHVrLmNvbSJ9.we-dont-check-signature" //nolint:lll
	s.req.Header.Add("Authentication", bearerPrefix+token)

	err := s.authMdlwr(func(_ echo.Context) error {
		s.Fail("unreachable")
		return nil
	})(s.ctx)
	s.assertHTTPCode(err, http.StatusBadRequest)
}

func (s *KeycloakTokenAuthSuite) TestNotBearerAuth() {
	const token = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJIR1lJcHN1UXlsZFNJZTB1T0JaeEpuQjBkZlFuTWI5LUlFcmx6NHk5ek9BIn0.eyJleHAiOjI2NjcxOTk1ODAsImlhdCI6MTY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOlsiY2hhdC11aS1jbGllbnQiLCJhY2NvdW50Il0sInN1YiI6IjVjYjQwZGMwLWEyNDktNDc4My1hMzAxLTllMWYzY2YzZWE0MSIsInR5cCI6IkJlYXJlciIsImF6cCI6ImNoYXQtdWktY2xpZW50Iiwibm9uY2UiOiJiYTM3ZmQ1YS04YzM5LTQ4MTQtYWZjYi05NTJhMThiNzI2N2QiLCJzZXNzaW9uX3N0YXRlIjoiZDg2ZDE5OGUtYzFjNS00ZWRkLTgzNTAtMzYxZWU1ODE3MWYyIiwiYWNyIjoiMCIsImFsbG93ZWQtb3JpZ2lucyI6WyIiLCIqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJvZmZsaW5lX2FjY2VzcyIsImRlZmF1bHQtcm9sZXMtYmFuayIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiY2hhdC11aS1jbGllbnQiOnsicm9sZXMiOlsic3VwcG9ydC1jaGF0LWNsaWVudCJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsInNpZCI6ImQ4NmQxOThlLWMxYzUtNGVkZC04MzUwLTM2MWVlNTgxNzFmMiIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJib25kMDA3IiwiZ2l2ZW5fbmFtZSI6IiIsImZhbWlseV9uYW1lIjoiIiwiZW1haWwiOiJib25kMDA3QHVrLmNvbSJ9.we-dont-check-signature" //nolint:lll
	s.req.Header.Add(echo.HeaderAuthorization, "Basic "+token)

	err := s.authMdlwr(func(_ echo.Context) error {
		s.Fail("unreachable")
		return nil
	})(s.ctx)
	s.assertHTTPCode(err, http.StatusBadRequest)
}

func (s *KeycloakTokenAuthSuite) TestIntrospectError() {
	const token = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJIR1lJcHN1UXlsZFNJZTB1T0JaeEpuQjBkZlFuTWI5LUlFcmx6NHk5ek9BIn0.eyJleHAiOjI2NjcxOTk1ODAsImlhdCI6MTY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiNWNiNDBkYzAtYTI0OS00NzgzLWEzMDEtOWUxZjNjZjNlYTQxIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiY2hhdC11aS1jbGllbnQiLCJub25jZSI6ImJhMzdmZDVhLThjMzktNDgxNC1hZmNiLTk1MmExOGI3MjY3ZCIsInNlc3Npb25fc3RhdGUiOiJkODZkMTk4ZS1jMWM1LTRlZGQtODM1MC0zNjFlZTU4MTcxZjIiLCJhY3IiOiIwIiwiYWxsb3dlZC1vcmlnaW5zIjpbIiIsIioiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbIm9mZmxpbmVfYWNjZXNzIiwiZGVmYXVsdC1yb2xlcy1iYW5rIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInJlc291cmNlX2FjY2VzcyI6eyJjaGF0LXVpLWNsaWVudCI6eyJyb2xlcyI6WyJzdXBwb3J0LWNoYXQtY2xpZW50Il19LCJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50IiwibWFuYWdlLWFjY291bnQtbGlua3MiLCJ2aWV3LXByb2ZpbGUiXX19LCJzY29wZSI6Im9wZW5pZCBwcm9maWxlIGVtYWlsIiwic2lkIjoiZDg2ZDE5OGUtYzFjNS00ZWRkLTgzNTAtMzYxZWU1ODE3MWYyIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsInByZWZlcnJlZF91c2VybmFtZSI6ImJvbmQwMDciLCJnaXZlbl9uYW1lIjoiIiwiZmFtaWx5X25hbWUiOiIiLCJlbWFpbCI6ImJvbmQwMDdAdWsuY29tIn0.we-dont-check-signature" //nolint:lll
	s.req.Header.Add(echo.HeaderAuthorization, bearerPrefix+token)

	s.introspector.EXPECT().IntrospectToken(s.req.Context(), token).Return(nil, context.Canceled)

	err := s.authMdlwr(func(_ echo.Context) error {
		s.Fail("unreachable")
		return nil
	})(s.ctx)
	s.assertHTTPCode(err, http.StatusUnauthorized)
	s.Require().ErrorIs(err, context.Canceled)
}

func (s *KeycloakTokenAuthSuite) TestInactiveToken() {
	const token = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJIR1lJcHN1UXlsZFNJZTB1T0JaeEpuQjBkZlFuTWI5LUlFcmx6NHk5ek9BIn0.eyJleHAiOjI2NjcxOTk1ODAsImlhdCI6MTY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiNWNiNDBkYzAtYTI0OS00NzgzLWEzMDEtOWUxZjNjZjNlYTQxIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiY2hhdC11aS1jbGllbnQiLCJub25jZSI6ImJhMzdmZDVhLThjMzktNDgxNC1hZmNiLTk1MmExOGI3MjY3ZCIsInNlc3Npb25fc3RhdGUiOiJkODZkMTk4ZS1jMWM1LTRlZGQtODM1MC0zNjFlZTU4MTcxZjIiLCJhY3IiOiIwIiwiYWxsb3dlZC1vcmlnaW5zIjpbIiIsIioiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbIm9mZmxpbmVfYWNjZXNzIiwiZGVmYXVsdC1yb2xlcy1iYW5rIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInJlc291cmNlX2FjY2VzcyI6eyJjaGF0LXVpLWNsaWVudCI6eyJyb2xlcyI6WyJzdXBwb3J0LWNoYXQtY2xpZW50Il19LCJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50IiwibWFuYWdlLWFjY291bnQtbGlua3MiLCJ2aWV3LXByb2ZpbGUiXX19LCJzY29wZSI6Im9wZW5pZCBwcm9maWxlIGVtYWlsIiwic2lkIjoiZDg2ZDE5OGUtYzFjNS00ZWRkLTgzNTAtMzYxZWU1ODE3MWYyIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsInByZWZlcnJlZF91c2VybmFtZSI6ImJvbmQwMDciLCJnaXZlbl9uYW1lIjoiIiwiZmFtaWx5X25hbWUiOiIiLCJlbWFpbCI6ImJvbmQwMDdAdWsuY29tIn0.we-dont-check-signature" //nolint:lll
	s.req.Header.Add(echo.HeaderAuthorization, bearerPrefix+token)

	s.introspector.EXPECT().IntrospectToken(s.req.Context(), token).
		Return(&keycloakclient.IntrospectTokenResult{Active: false}, nil)

	err := s.authMdlwr(func(_ echo.Context) error {
		s.Fail("unreachable")
		return nil
	})(s.ctx)
	s.assertHTTPCode(err, http.StatusUnauthorized)
}

func (s *KeycloakTokenAuthSuite) TestInvalidExpiresAt() {
	const token = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJIR1lJcHN1UXlsZFNJZTB1T0JaeEpuQjBkZlFuTWI5LUlFcmx6NHk5ek9BIn0.eyJleHAiOjE2NjcxOTk1ODAsImlhdCI6MTY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOlsiY2hhdC11aS1jbGllbnQiLCJhY2NvdW50Il0sInN1YiI6IjVjYjQwZGMwLWEyNDktNDc4My1hMzAxLTllMWYzY2YzZWE0MSIsInR5cCI6IkJlYXJlciIsImF6cCI6ImNoYXQtdWktY2xpZW50Iiwibm9uY2UiOiJiYTM3ZmQ1YS04YzM5LTQ4MTQtYWZjYi05NTJhMThiNzI2N2QiLCJzZXNzaW9uX3N0YXRlIjoiZDg2ZDE5OGUtYzFjNS00ZWRkLTgzNTAtMzYxZWU1ODE3MWYyIiwiYWNyIjoiMCIsImFsbG93ZWQtb3JpZ2lucyI6WyIiLCIqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJvZmZsaW5lX2FjY2VzcyIsImRlZmF1bHQtcm9sZXMtYmFuayIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiY2hhdC11aS1jbGllbnQiOnsicm9sZXMiOlsic3VwcG9ydC1jaGF0LWNsaWVudCJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsInNpZCI6ImQ4NmQxOThlLWMxYzUtNGVkZC04MzUwLTM2MWVlNTgxNzFmMiIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJib25kMDA3IiwiZ2l2ZW5fbmFtZSI6IiIsImZhbWlseV9uYW1lIjoiIiwiZW1haWwiOiJib25kMDA3QHVrLmNvbSJ9.we-dont-check-signature" //nolint:lll
	s.req.Header.Add(echo.HeaderAuthorization, bearerPrefix+token)

	s.introspector.EXPECT().IntrospectToken(s.req.Context(), token).
		Return(&keycloakclient.IntrospectTokenResult{Active: true}, nil)

	err := s.authMdlwr(func(_ echo.Context) error {
		s.Fail("unreachable")
		return nil
	})(s.ctx)
	s.assertHTTPCode(err, http.StatusUnauthorized)

	var jwtErr *jwt.ValidationError
	s.Require().ErrorAs(err, &jwtErr)
	s.Empty(jwtErr.Errors ^ jwt.ValidationErrorExpired)
}

func (s *KeycloakTokenAuthSuite) TestInvalidIssuedAt() {
	const token = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJIR1lJcHN1UXlsZFNJZTB1T0JaeEpuQjBkZlFuTWI5LUlFcmx6NHk5ek9BIn0.eyJleHAiOjI2NjcxOTk1ODAsImlhdCI6MjY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOlsiY2hhdC11aS1jbGllbnQiLCJhY2NvdW50Il0sInN1YiI6IjVjYjQwZGMwLWEyNDktNDc4My1hMzAxLTllMWYzY2YzZWE0MSIsInR5cCI6IkJlYXJlciIsImF6cCI6ImNoYXQtdWktY2xpZW50Iiwibm9uY2UiOiJiYTM3ZmQ1YS04YzM5LTQ4MTQtYWZjYi05NTJhMThiNzI2N2QiLCJzZXNzaW9uX3N0YXRlIjoiZDg2ZDE5OGUtYzFjNS00ZWRkLTgzNTAtMzYxZWU1ODE3MWYyIiwiYWNyIjoiMCIsImFsbG93ZWQtb3JpZ2lucyI6WyIiLCIqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJvZmZsaW5lX2FjY2VzcyIsImRlZmF1bHQtcm9sZXMtYmFuayIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiY2hhdC11aS1jbGllbnQiOnsicm9sZXMiOlsic3VwcG9ydC1jaGF0LWNsaWVudCJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsInNpZCI6ImQ4NmQxOThlLWMxYzUtNGVkZC04MzUwLTM2MWVlNTgxNzFmMiIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJib25kMDA3IiwiZ2l2ZW5fbmFtZSI6IiIsImZhbWlseV9uYW1lIjoiIiwiZW1haWwiOiJib25kMDA3QHVrLmNvbSJ9.we-dont-check-signature" //nolint:lll
	s.req.Header.Add(echo.HeaderAuthorization, bearerPrefix+token)

	s.introspector.EXPECT().IntrospectToken(s.req.Context(), token).
		Return(&keycloakclient.IntrospectTokenResult{Active: true}, nil)

	err := s.authMdlwr(func(_ echo.Context) error {
		s.Fail("unreachable")
		return nil
	})(s.ctx)
	s.assertHTTPCode(err, http.StatusUnauthorized)

	var jwtErr *jwt.ValidationError
	s.Require().ErrorAs(err, &jwtErr)
	s.Empty(jwtErr.Errors ^ jwt.ValidationErrorIssuedAt)
}

func (s *KeycloakTokenAuthSuite) TestNoSubject() {
	const token = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJIR1lJcHN1UXlsZFNJZTB1T0JaeEpuQjBkZlFuTWI5LUlFcmx6NHk5ek9BIn0.eyJleHAiOjI2NjcxOTk1ODAsImlhdCI6MTY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOlsiY2hhdC11aS1jbGllbnQiLCJhY2NvdW50Il0sInR5cCI6IkJlYXJlciIsImF6cCI6ImNoYXQtdWktY2xpZW50Iiwibm9uY2UiOiJiYTM3ZmQ1YS04YzM5LTQ4MTQtYWZjYi05NTJhMThiNzI2N2QiLCJzZXNzaW9uX3N0YXRlIjoiZDg2ZDE5OGUtYzFjNS00ZWRkLTgzNTAtMzYxZWU1ODE3MWYyIiwiYWNyIjoiMCIsImFsbG93ZWQtb3JpZ2lucyI6WyIiLCIqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJvZmZsaW5lX2FjY2VzcyIsImRlZmF1bHQtcm9sZXMtYmFuayIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiY2hhdC11aS1jbGllbnQiOnsicm9sZXMiOlsic3VwcG9ydC1jaGF0LWNsaWVudCJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsInNpZCI6ImQ4NmQxOThlLWMxYzUtNGVkZC04MzUwLTM2MWVlNTgxNzFmMiIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJib25kMDA3IiwiZ2l2ZW5fbmFtZSI6IiIsImZhbWlseV9uYW1lIjoiIiwiZW1haWwiOiJib25kMDA3QHVrLmNvbSJ9.we-dont-check-signature" //nolint:lll
	s.req.Header.Add(echo.HeaderAuthorization, bearerPrefix+token)

	s.introspector.EXPECT().IntrospectToken(s.req.Context(), token).
		Return(&keycloakclient.IntrospectTokenResult{Active: true}, nil)

	err := s.authMdlwr(func(_ echo.Context) error {
		s.Fail("unreachable")
		return nil
	})(s.ctx)
	s.assertHTTPCode(err, http.StatusUnauthorized)
	s.Require().ErrorIs(err, middlewares.ErrSubjectNotDefined)
}

func (s *KeycloakTokenAuthSuite) TestSubjectIsZeroUUID() {
	const token = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJIR1lJcHN1UXlsZFNJZTB1T0JaeEpuQjBkZlFuTWI5LUlFcmx6NHk5ek9BIn0.eyJleHAiOjI2NjcxOTk1ODAsImlhdCI6MTY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOlsiY2hhdC11aS1jbGllbnQiLCJhY2NvdW50Il0sInN1YiI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCIsInR5cCI6IkJlYXJlciIsImF6cCI6ImNoYXQtdWktY2xpZW50Iiwibm9uY2UiOiJiYTM3ZmQ1YS04YzM5LTQ4MTQtYWZjYi05NTJhMThiNzI2N2QiLCJzZXNzaW9uX3N0YXRlIjoiZDg2ZDE5OGUtYzFjNS00ZWRkLTgzNTAtMzYxZWU1ODE3MWYyIiwiYWNyIjoiMCIsImFsbG93ZWQtb3JpZ2lucyI6WyIiLCIqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJvZmZsaW5lX2FjY2VzcyIsImRlZmF1bHQtcm9sZXMtYmFuayIsInVtYV9hdXRob3JpemF0aW9uIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiY2hhdC11aS1jbGllbnQiOnsicm9sZXMiOlsic3VwcG9ydC1jaGF0LWNsaWVudCJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsInNpZCI6ImQ4NmQxOThlLWMxYzUtNGVkZC04MzUwLTM2MWVlNTgxNzFmMiIsImVtYWlsX3ZlcmlmaWVkIjp0cnVlLCJwcmVmZXJyZWRfdXNlcm5hbWUiOiJib25kMDA3IiwiZ2l2ZW5fbmFtZSI6IiIsImZhbWlseV9uYW1lIjoiIiwiZW1haWwiOiJib25kMDA3QHVrLmNvbSJ9.we-dont-check-signature" //nolint:lll
	s.req.Header.Add(echo.HeaderAuthorization, bearerPrefix+token)

	s.introspector.EXPECT().IntrospectToken(s.req.Context(), token).
		Return(&keycloakclient.IntrospectTokenResult{Active: true}, nil)

	err := s.authMdlwr(func(_ echo.Context) error {
		s.Fail("unreachable")
		return nil
	})(s.ctx)
	s.assertHTTPCode(err, http.StatusUnauthorized)
	s.Require().ErrorIs(err, middlewares.ErrSubjectNotDefined)
}

func (s *KeycloakTokenAuthSuite) TestNoResourceAccess_EmptyMap() {
	const token = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJIR1lJcHN1UXlsZFNJZTB1T0JaeEpuQjBkZlFuTWI5LUlFcmx6NHk5ek9BIn0.eyJleHAiOjI2NjcxOTk1ODAsImlhdCI6MTY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiNWNiNDBkYzAtYTI0OS00NzgzLWEzMDEtOWUxZjNjZjNlYTQxIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiY2hhdC11aS1jbGllbnQiLCJub25jZSI6ImJhMzdmZDVhLThjMzktNDgxNC1hZmNiLTk1MmExOGI3MjY3ZCIsInNlc3Npb25fc3RhdGUiOiJkODZkMTk4ZS1jMWM1LTRlZGQtODM1MC0zNjFlZTU4MTcxZjIiLCJhY3IiOiIwIiwiYWxsb3dlZC1vcmlnaW5zIjpbIiIsIioiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbIm9mZmxpbmVfYWNjZXNzIiwiZGVmYXVsdC1yb2xlcy1iYW5rIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInJlc291cmNlX2FjY2VzcyI6e30sInNjb3BlIjoib3BlbmlkIHByb2ZpbGUgZW1haWwiLCJzaWQiOiJkODZkMTk4ZS1jMWM1LTRlZGQtODM1MC0zNjFlZTU4MTcxZjIiLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwicHJlZmVycmVkX3VzZXJuYW1lIjoiYm9uZDAwNyIsImdpdmVuX25hbWUiOiIiLCJmYW1pbHlfbmFtZSI6IiIsImVtYWlsIjoiYm9uZDAwN0B1ay5jb20ifQ.we-dont-check-signature" //nolint:lll
	s.req.Header.Add(echo.HeaderAuthorization, bearerPrefix+token)

	s.introspector.EXPECT().IntrospectToken(s.req.Context(), token).
		Return(&keycloakclient.IntrospectTokenResult{Active: true}, nil)

	err := s.authMdlwr(func(_ echo.Context) error {
		s.Fail("unreachable")
		return nil
	})(s.ctx)
	s.assertHTTPCode(err, http.StatusUnauthorized)
	s.Require().ErrorIs(err, middlewares.ErrNoAllowedResources)
}

func (s *KeycloakTokenAuthSuite) TestNoResourceAccess_NoKey() {
	const token = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJIR1lJcHN1UXlsZFNJZTB1T0JaeEpuQjBkZlFuTWI5LUlFcmx6NHk5ek9BIn0.eyJleHAiOjI2NjcxOTk1ODAsImlhdCI6MTY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiNWNiNDBkYzAtYTI0OS00NzgzLWEzMDEtOWUxZjNjZjNlYTQxIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiY2hhdC11aS1jbGllbnQiLCJub25jZSI6ImJhMzdmZDVhLThjMzktNDgxNC1hZmNiLTk1MmExOGI3MjY3ZCIsInNlc3Npb25fc3RhdGUiOiJkODZkMTk4ZS1jMWM1LTRlZGQtODM1MC0zNjFlZTU4MTcxZjIiLCJhY3IiOiIwIiwiYWxsb3dlZC1vcmlnaW5zIjpbIiIsIioiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbIm9mZmxpbmVfYWNjZXNzIiwiZGVmYXVsdC1yb2xlcy1iYW5rIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInNjb3BlIjoib3BlbmlkIHByb2ZpbGUgZW1haWwiLCJzaWQiOiJkODZkMTk4ZS1jMWM1LTRlZGQtODM1MC0zNjFlZTU4MTcxZjIiLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwicHJlZmVycmVkX3VzZXJuYW1lIjoiYm9uZDAwNyIsImdpdmVuX25hbWUiOiIiLCJmYW1pbHlfbmFtZSI6IiIsImVtYWlsIjoiYm9uZDAwN0B1ay5jb20ifQ.we-dont-check-signature" //nolint:lll
	s.req.Header.Add(echo.HeaderAuthorization, bearerPrefix+token)

	s.introspector.EXPECT().IntrospectToken(s.req.Context(), token).
		Return(&keycloakclient.IntrospectTokenResult{Active: true}, nil)

	err := s.authMdlwr(func(_ echo.Context) error {
		s.Fail("unreachable")
		return nil
	})(s.ctx)
	s.assertHTTPCode(err, http.StatusUnauthorized)
	s.Require().ErrorIs(err, middlewares.ErrNoAllowedResources)
}

func (s *KeycloakTokenAuthSuite) TestNoResourceRole_NoNeededResource() {
	const token = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJIR1lJcHN1UXlsZFNJZTB1T0JaeEpuQjBkZlFuTWI5LUlFcmx6NHk5ek9BIn0.eyJleHAiOjI2NjcxOTk1ODAsImlhdCI6MTY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiNWNiNDBkYzAtYTI0OS00NzgzLWEzMDEtOWUxZjNjZjNlYTQxIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiY2hhdC11aS1jbGllbnQiLCJub25jZSI6ImJhMzdmZDVhLThjMzktNDgxNC1hZmNiLTk1MmExOGI3MjY3ZCIsInNlc3Npb25fc3RhdGUiOiJkODZkMTk4ZS1jMWM1LTRlZGQtODM1MC0zNjFlZTU4MTcxZjIiLCJhY3IiOiIwIiwiYWxsb3dlZC1vcmlnaW5zIjpbIiIsIioiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbIm9mZmxpbmVfYWNjZXNzIiwiZGVmYXVsdC1yb2xlcy1iYW5rIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInJlc291cmNlX2FjY2VzcyI6eyJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50IiwibWFuYWdlLWFjY291bnQtbGlua3MiLCJ2aWV3LXByb2ZpbGUiXX19LCJzY29wZSI6Im9wZW5pZCBwcm9maWxlIGVtYWlsIiwic2lkIjoiZDg2ZDE5OGUtYzFjNS00ZWRkLTgzNTAtMzYxZWU1ODE3MWYyIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsInByZWZlcnJlZF91c2VybmFtZSI6ImJvbmQwMDciLCJnaXZlbl9uYW1lIjoiIiwiZmFtaWx5X25hbWUiOiIiLCJlbWFpbCI6ImJvbmQwMDdAdWsuY29tIn0.we-dont-check-signature" //nolint:lll
	s.req.Header.Add(echo.HeaderAuthorization, bearerPrefix+token)

	s.introspector.EXPECT().IntrospectToken(s.req.Context(), token).
		Return(&keycloakclient.IntrospectTokenResult{Active: true}, nil)

	err := s.authMdlwr(func(_ echo.Context) error {
		s.Fail("unreachable")
		return nil
	})(s.ctx)
	s.assertHTTPCode(err, http.StatusUnauthorized)
	s.Require().ErrorIs(err, middlewares.ErrNoRequiredResourceRole)
}

func (s *KeycloakTokenAuthSuite) TestNoResourceRole_NoNeededRole() {
	const token = "eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJIR1lJcHN1UXlsZFNJZTB1T0JaeEpuQjBkZlFuTWI5LUlFcmx6NHk5ek9BIn0.eyJleHAiOjI2NjcxOTk1ODAsImlhdCI6MTY2NzE5OTI4MCwiYXV0aF90aW1lIjoxNjY3MTk4OTI4LCJqdGkiOiI5NGQ3ZDBkNS0zZTZmLTQ5NGItYTkzYy1hYjliMDkxMzQ3YmEiLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjMwMTAvcmVhbG1zL0JhbmsiLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiNWNiNDBkYzAtYTI0OS00NzgzLWEzMDEtOWUxZjNjZjNlYTQxIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoiY2hhdC11aS1jbGllbnQiLCJub25jZSI6ImJhMzdmZDVhLThjMzktNDgxNC1hZmNiLTk1MmExOGI3MjY3ZCIsInNlc3Npb25fc3RhdGUiOiJkODZkMTk4ZS1jMWM1LTRlZGQtODM1MC0zNjFlZTU4MTcxZjIiLCJhY3IiOiIwIiwiYWxsb3dlZC1vcmlnaW5zIjpbIiIsIioiXSwicmVhbG1fYWNjZXNzIjp7InJvbGVzIjpbIm9mZmxpbmVfYWNjZXNzIiwiZGVmYXVsdC1yb2xlcy1iYW5rIiwidW1hX2F1dGhvcml6YXRpb24iXX0sInJlc291cmNlX2FjY2VzcyI6eyJjaGF0LXVpLWNsaWVudCI6eyJyb2xlcyI6WyJhYnJhY2FkYWJyYSJdfSwiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIiwic3VwcG9ydC1jaGF0LWNsaWVudCJdfX0sInNjb3BlIjoib3BlbmlkIHByb2ZpbGUgZW1haWwiLCJzaWQiOiJkODZkMTk4ZS1jMWM1LTRlZGQtODM1MC0zNjFlZTU4MTcxZjIiLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwicHJlZmVycmVkX3VzZXJuYW1lIjoiYm9uZDAwNyIsImdpdmVuX25hbWUiOiIiLCJmYW1pbHlfbmFtZSI6IiIsImVtYWlsIjoiYm9uZDAwN0B1ay5jb20ifQ.we-dont-check-signature" //nolint:lll
	s.req.Header.Add(echo.HeaderAuthorization, bearerPrefix+token)

	s.introspector.EXPECT().IntrospectToken(s.req.Context(), token).
		Return(&keycloakclient.IntrospectTokenResult{Active: true}, nil)

	err := s.authMdlwr(func(_ echo.Context) error {
		s.Fail("unreachable")
		return nil
	})(s.ctx)
	s.assertHTTPCode(err, http.StatusUnauthorized)
	s.Require().ErrorIs(err, middlewares.ErrNoRequiredResourceRole)
}

func (s *KeycloakTokenAuthSuite) assertHTTPCode(err error, code int) {
	var httpErr *echo.HTTPError
	s.Require().ErrorAs(err, &httpErr)
	s.Equal(httpErr.Code, code)
}

func TestMustUserID_NoUID(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(t, err)

	assert.Panics(t, func() {
		middlewares.MustUserID(echo.New().NewContext(req, httptest.NewRecorder()))
	})
}
