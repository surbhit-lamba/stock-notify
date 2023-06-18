package router

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func SetupRouter(ctx context.Context, nrApp *newrelic.Application) *gin.Engine {
	r := gin.New()

	r.Use(nrgin.Middleware(nrApp))

	setRoutes(r.Group("stocks"), ctx)

	return r
}

func setRoutes(r *gin.RouterGroup, ctx context.Context) {
	setHealthRoutes(r, ctx)
}

func setHealthRoutes(r *gin.RouterGroup, ctx context.Context) {
	r.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "yoo healthy!!")
	})
}

