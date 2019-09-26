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

// Define custom request context format for logger
type RequestContext struct {
	Method string                 `json:"method"`
	Path   string                 `json:"path"`
	Query  string                 `json:"query"`
	Body   map[string]interface{} `json:"body"`
	Header map[string]interface{} `json:"headers"`
}

type ResponseContext struct {
	Method string                 `json:"method"`
	Path   string                 `json:"path"`
	Query  string                 `json:"query"`
	Body   string                 `json:"body"`
	Header map[string]interface{} `json:"headers"`
}

func getRequestField(method string, path string, query string, body map[string]interface{}, headers map[string][]string) (ctx zap.Field, err error) {
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
	byte, err := json.Marshal(requestContext)
	if err != nil {
		return
	}
	raw := json.RawMessage(byte)
	ctx = zap.Any("context", &raw)
	return
}

func getResponseField(method string, path string, query string, body string, headers map[string][]string) (ctx zap.Field, err error) {
	header := make(map[string]interface{})
	for key, h := range headers {
		header[key] = strings.Join(h, ",")
	}

	maxBodyLen := len(body)
	if maxBodyLen > 1024 {
		maxBodyLen = 1024
	}

	requestContext := ResponseContext{
		Method: method,
		Path:   path,
		Query:  query,
		Body:   body[0:maxBodyLen],
		Header: header,
	}
	byte, err := json.Marshal(requestContext)
	if err != nil {
		return
	}
	raw := json.RawMessage(byte)
	ctx = zap.Any("context", &raw)
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
				json.Unmarshal(rawBody, &reqBody)
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

		requestUuidField := zap.String("request-uuid", u1)

		// wrapper request context
		reqContextField, err := getRequestField(c.Request.Method, c.Request.URL.RequestURI(), c.Request.URL.RawQuery, reqBody, c.Request.Header)
		if err != nil {
			fmt.Println(err.Error())
			logger.Info("api_request", hostField, requestUuidField)
		} else {
			logger.Info("api_request", hostField, requestUuidField, reqContextField)
		}

		// wrapper error handle
		defer func() {
			if err := recover(); err != nil {
				goErr := errors.Wrap(err, 3)
				logger.Error(goErr.Error(), hostField, requestUuidField, reqContextField)

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
		resContextField, err := getResponseField(c.Request.Method, c.Request.URL.RequestURI(), c.Request.URL.RawQuery, blw.body.String(), c.Writer.Header())
		if err != nil {
			fmt.Println(err.Error())
			logger.Info("api_response", hostField, requestUuidField, statusCodeField, latencyField)
		} else {
			logger.Info("api_response", hostField, requestUuidField, resContextField, statusCodeField, latencyField)
		}

	}
}
