package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/weikunlu/go-api-template/core/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// RequestContext is present a request context for logger information
type RequestContext struct {
	Method string                 `json:"method"`
	Path   string                 `json:"path"`
	Query  string                 `json:"query"`
	Body   map[string]interface{} `json:"body"`
	Header map[string]interface{} `json:"headers"`
}

// ResponseContext is present a response context for logger information
type ResponseContext struct {
	Body   string                 `json:"body"`
	Header map[string]interface{} `json:"headers"`
}

func getRequestFields(method string, path string, query string, body map[string]interface{}, headers map[string][]string) (ctx zap.Field, err error) {
	header := make(map[string]interface{})
	for key, h := range headers {
		header[key] = strings.Join(h, ",")
	}
	requestContext := RequestContext{
		Method: method,
		Path:   path,
		Query:  query,
		Body:   body,
		Header: header,
	}
	b, err := json.Marshal(requestContext)
	if err != nil {
		return
	}
	raw := json.RawMessage(b)
	ctx = zap.Any("request", &raw)
	return
}

func getResponseFields(body string, headers map[string][]string) (ctx zap.Field, err error) {
	header := make(map[string]interface{})
	for key, h := range headers {
		header[key] = strings.Join(h, ",")
	}
	requestContext := ResponseContext{
		Body:   body,
		Header: header,
	}
	b, err := json.Marshal(requestContext)
	if err != nil {
		return
	}
	raw := json.RawMessage(b)
	ctx = zap.Any("response", &raw)
	return
}

func AccessLogMiddleware(f func(c *gin.Context, err interface{})) gin.HandlerFunc {

	logger := log.NewLogger(zapcore.InfoLevel, "access-log", false)

	return func(c *gin.Context) {
		start := time.Now()

		u1 := c.MustGet("request_uuid").(string)

		var rawBody []byte
		if c.Request.Body != nil {
			rawBody, _ = ioutil.ReadAll(c.Request.Body)
		}

		var reqBody = make(map[string]interface{})
		if c.Request.ContentLength > 0 {
			switch c.ContentType() {
			case "application/json":
				_ = json.Unmarshal(rawBody, &reqBody)
				break
			case "application/x-www-form-urlencoded":
				params, _ := url.ParseQuery(string(rawBody))
				for key, value := range params {
					reqBody[key] = value[0]
				}
				break
			}
		}

		hostField := zap.String("host_domain", c.Request.Host)

		requestUUIDField := zap.String("request-uuid", u1)

		// wrapper request context
		reqContextField, err := getRequestFields(c.Request.Method, c.Request.URL.RequestURI(), c.Request.URL.RawQuery, reqBody, c.Request.Header)
		if err != nil {
			reqContextField = zap.String("request", err.Error())
		}

		// wrapper error handle
		defer func() {
			if err := recover(); err != nil {
				goErr := errors.Wrap(err, 3)
				logger.Error("request", log.GetErrorMessageField(goErr.Error()), hostField, requestUUIDField, reqContextField)

				f(c, err)
			}
		}()

		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		c.Next()

		statusCodeField := zap.String("status_code", strconv.Itoa(c.Writer.Status()))

		latencyField := zap.String("latency", fmt.Sprintf("%13v", time.Now().Sub(start)))

		// wrapper response context
		resContextField, err := getResponseFields(blw.body.String(), c.Writer.Header())
		if err != nil {
			resContextField = zap.String("response", err.Error())
		}

		logger.Info("request", hostField, requestUUIDField, reqContextField, resContextField, statusCodeField, latencyField)
	}
}
