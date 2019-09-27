package apiv1

import (
	"github.com/gin-gonic/gin"
	"github.com/weikunlu/go-api-template/api"
	"github.com/weikunlu/go-api-template/core/log"
	"net/http"
)

func HelloController(r *gin.RouterGroup) {
	r.GET("/", GetHello)
}

func GetHello(c *gin.Context) {

	log.WithContext(c).Info("log for hello")

	c.JSON(http.StatusOK, api.GetSuccessResponse("hello"))
}