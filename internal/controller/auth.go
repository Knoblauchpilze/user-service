package controller

import (
	"net/http"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/Knoblauchpilze/user-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const apiKeyHeaderKey = "X-Api-Key"

func AuthEndpoints(service service.AuthService) Routes {
	var out Routes

	authHandler := createServiceAwareHttpHandler(authUser, service)
	auth := rest.NewRoute(http.MethodGet, "/auth", authHandler)
	out = append(out, auth)

	return out
}

// authUser godoc
//
// @Summary Authenticate API key
// @Description Validates the API key provided in the request header.
// @Tags auth
// @Produce json
// @Param X-Api-Key header string true "API key"
// @Success 204
// @Failure 400 {object} rest.ResponseEnvelope[string] "Invalid API key"
// @Failure 403 {object} rest.ResponseEnvelope[string] "User is not authenticated"
// @Failure 500 {object} rest.ResponseEnvelope[string] "Internal server error"
// @Router /users/auth [get]
func authUser(c *gin.Context, s service.AuthService) {
	apiKey, exists := tryGetApiKeyHeader(c.Request)
	if !exists {
		c.AbortWithStatusJSON(http.StatusBadRequest, "Invalid API key")
		return
	}

	_, err := s.Authenticate(c.Request.Context(), apiKey)
	if err != nil {
		if isUserNotAuthenticated(err) {
			c.AbortWithStatusJSON(http.StatusForbidden, err)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
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
	return err == service.ErrUserNotAuthenticated || err == service.ErrAuthenticationExpired
}
