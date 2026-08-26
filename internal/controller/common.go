package controller

import (
	"github.com/Knoblauchpilze/backend-toolkit/pkg/rest"
	"github.com/gin-gonic/gin"
)

type Routes []*rest.Route

type serviceAwareHttpHandler[T any] func(*gin.Context, T)

func createServiceAwareHttpHandler[T any](handler serviceAwareHttpHandler[T], service T) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(c, service)
	}
}
