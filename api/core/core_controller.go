package apicore

import (
	"github.com/gin-gonic/gin"
	"github.com/weikunlu/go-api-template/api"
	"github.com/weikunlu/go-api-template/config"
	"net/http"
)

func CoreController(r *gin.RouterGroup) {
	r.GET("/health_check", HealthCheck)
	r.POST("/version", BuildInformation)
}

func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, "")
}

func BuildInformation(c *gin.Context) {
	c.JSON(http.StatusOK, api.GetSuccessResponse(config.GetBuildConfig()))
}
