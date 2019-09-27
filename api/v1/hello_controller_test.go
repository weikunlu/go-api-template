package apiv1_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/weikunlu/go-api-template/api/v1"
	"github.com/weikunlu/go-api-template/core/apptest"
	"net/http"
	"testing"
)

func TestGetHello(t *testing.T) {
	router, group := apptest.SetupRouter("hello")
	apiv1.HelloController(group)

	w := apptest.PerformHttpRequest(router, "GET", "hello/", nil, nil)

	assert.Equal(t, http.StatusOK, w.Code)
}
