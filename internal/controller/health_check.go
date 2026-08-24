package controller

import (
	"net/http"

	"github.com/Knoblauchpilze/backend-toolkit/pkg/db"
	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/gin-gonic/gin"
)

func HealthCheckEndpoints(pool db.Connection) Routes {
	var out Routes

	getHandler := createServiceAwareHttpHandler(healthcheck, pool)
	get := rest.NewRoute(http.MethodGet, "/healthcheck", getHandler)
	out = append(out, get)

	return out
}

// healthcheck godoc
//
// @Summary Health check
// @Description Verifies that the service can reach its database.
// @Tags health
// @Produce json
// @Success 200 {object} rest.ResponseEnvelope[string]
// @Failure 503 {object} rest.ResponseEnvelope[string] "Database unavailable"
// @Router /healthcheck [get]
func healthcheck(c *gin.Context, pool db.Connection) {
	err := pool.Ping(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, err)
		return
	}

	c.JSON(http.StatusOK, "OK")
}
