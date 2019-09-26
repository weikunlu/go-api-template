package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/weikunlu/go-api-template/core/log"
	"gopkg.in/satori/go.uuid.v1"
)

func RequestUuidMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		// generate UUID for each request
		u1 := uuid.NewV1().String()
		c.Set("request_uuid", u1)

		// register logger to gin.context
		log.NewContext(c, log.GetStringField("request_uuid", u1))

		c.Next()
	}
}
