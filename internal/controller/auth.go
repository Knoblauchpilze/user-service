package controller

import (
	"net/http"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/errors"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/Knoblauchpilze/user-service/internal/service"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const apiKeyHeaderKey = "X-Api-Key"

func AuthEndpoints(service service.AuthService) rest.Routes {
	var out rest.Routes

	authHandler := createServiceAwareHttpHandler(authUser, service)
	auth := rest.NewRoute(http.MethodGet, "/auth", authHandler)
	out = append(out, auth)

	return out
}

func authUser(c echo.Context, s service.AuthService) error {
	apiKey, exists := tryGetApiKeyHeader(c.Request())
	if !exists {
		return c.JSON(http.StatusBadRequest, "Invalid API key")
	}

	_, err := s.Authenticate(c.Request().Context(), apiKey)
	if err != nil {
		if isUserNotAuthenticated(err) {
			return c.JSON(http.StatusForbidden, err)
		}

		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}

func tryGetApiKeyHeader(req *http.Request) (uuid.UUID, bool) {
	apiKeys, ok := req.Header[apiKeyHeaderKey]
	if !ok {
		return uuid.UUID{}, false
	}
	if len(apiKeys) != 1 {
		return uuid.UUID{}, false
	}

	apiKey, err := uuid.Parse(apiKeys[0])
	if err != nil {
		return uuid.UUID{}, false
	}

	return apiKey, true
}

func isUserNotAuthenticated(err error) bool {
	return errors.IsErrorWithCode(err, service.UserNotAuthenticated) || errors.IsErrorWithCode(err, service.AuthenticationExpired)
}
