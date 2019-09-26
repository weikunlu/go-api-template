package apptest

import "github.com/gin-gonic/gin"

func SetupRouter(uri string) (r *gin.Engine, group *gin.RouterGroup) {
	gin.SetMode(gin.TestMode)
	r = gin.New()
	group = r.Group(uri)
	return
}
