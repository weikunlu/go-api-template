package apptest

import (
	"io"
	"net/http"
	"net/http/httptest"
)

type Header struct {
	Name  string
	Value string
}

func PerformHttpRequest(r http.Handler, method, path string, headers []Header, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)

	for _, h := range headers {
		req.Header.Set(h.Name, h.Value)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func GetHeader(name string, value string) Header {
	return Header{
		Name:  name,
		Value: value,
	}
}
