package middlewares

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/weikunlu/go-api-template/api"
	"github.com/weikunlu/go-api-template/core/log"
	"go.uber.org/zap"
	"net/http"
	"sort"
	"time"
)

// SignatureValidationMiddleware is a sample of HMAC validation
func SignatureValidationMiddleware() gin.HandlerFunc {

	var headerXIdentitySignature = "X-Identity-Signature"
	var timeLayout = "20060102150405"
	var durationOfExpired = time.Minute

	return func(c *gin.Context) {
		r := c.Request

		if len(r.Header[headerXIdentitySignature]) == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, api.GetErrorResponse(nil, "identity signature is needed"))
			return
		}

		err := r.ParseForm()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, api.GetErrorResponse(nil, err.Error()))
			return
		}
		form := r.Form

		createdTime := form.Get("t")
		if len(createdTime) == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, api.GetErrorResponse(nil, "incomplete query parameter for signature"))
			return
		}

		keys := make([]string, 0, len(form))
		for key := range form {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		msg := r.Host
		for _, key := range keys {
			msg += key + form[key][0]
		}

		secret := "get_secret_key_by_client_id"
		if err != nil {
			log.WithContext(c).Warn("fail-validate-token", zap.String("invalid_client_id", form.Get("client_id")))
			c.AbortWithStatusJSON(http.StatusForbidden, api.GetErrorResponse(nil, err.Error()))
			return
		}

		// generating signature
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write([]byte(msg))
		expectedSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

		actualSignature := r.Header[headerXIdentitySignature][0]
		if !(actualSignature == expectedSignature) {
			log.WithContext(c).Warn("fail-validate-token", zap.String("raw_message", msg))
			c.AbortWithStatusJSON(http.StatusForbidden, api.GetErrorResponse(nil, "invalid signature"))
			return
		}

		requestTime, err := time.Parse(timeLayout, createdTime)
		if err != nil {
			log.WithContext(c).Warn("fail-validate-token", zap.String("request_query_t", createdTime))
			c.AbortWithStatusJSON(http.StatusBadRequest, api.GetErrorResponse(nil, err.Error()))
			return
		}

		receiveTime := time.Now().UTC()
		diff := receiveTime.Sub(requestTime)
		if diff.Minutes() < 0 || diff.Minutes() > durationOfExpired.Minutes() {
			log.WithContext(c).Warn("fail-validate-token", zap.Any("duration_expired", diff))
			c.AbortWithStatusJSON(http.StatusForbidden, api.GetErrorResponse(nil, "signature expired"))
			return
		}

		c.Request.Form.Set("client_secret", secret)

		c.Next()
	}
}
