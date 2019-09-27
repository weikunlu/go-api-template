package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/weikunlu/go-api-template/api"
	"github.com/weikunlu/go-api-template/api/core"
	"github.com/weikunlu/go-api-template/api/v1"
	"github.com/weikunlu/go-api-template/config"
	"github.com/weikunlu/go-api-template/middlewares"
	"net/http"
)

func recoveryHandler(c *gin.Context, err interface{}) {
	c.JSON(http.StatusInternalServerError, api.GetErrorResponse(err, ""))
}

func SetupRouter() (r *gin.Engine) {
	cfg := config.GetAppConfig()

	r = gin.New()

	if cfg.AppEnv == "local" {
		r.Use(gin.Logger())
	}

	r.Use(cors.Default())
	r.Use(middlewares.RequestUuidMiddleware())

	root := r.Group("")
	{
		apicore.CoreController(root)
	}

	api := r.Group("api")
	api.Use(middlewares.AccessLogMiddleware(recoveryHandler))

	hello := api.Group("v1/hello")
	{
		apiv1.HelloController(hello)
	}

	return
}
