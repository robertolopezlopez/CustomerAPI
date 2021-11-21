package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {

	tests := map[string]struct {
		headers      map[string]string
		expectedCode int
		expectedBody string
	}{
		"200 ok": {
			headers:      map[string]string{"X-Token": "test"},
			expectedCode: http.StatusOK,
			expectedBody: "pong",
		},
		"401 nok": {
			headers:      map[string]string{"X-Token": ""},
			expectedCode: http.StatusUnauthorized,
		},
		"401 empty": {
			headers:      map[string]string{},
			expectedCode: http.StatusUnauthorized,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			router := SetupRouter()
			w := httptest.NewRecorder()

			req, _ := http.NewRequest("GET", "/ping", nil)
			for key, value := range test.headers {
				req.Header.Add(key, value)
			}
			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedCode, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
		})
	}
}

// TODO MORE TESTS
