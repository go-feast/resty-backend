package middleware_test

import (
	"github.com/go-feast/resty-backend/internal/http/middleware"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestMiddleware(t *testing.T) {
	t.Parallel()
	t.Run("healthz", testRoute(middleware.Healthz))
	t.Run("readyz", testRoute(middleware.Readyz))
	t.Run("ping", testRoute(middleware.Ping))
}

func testRoute(h http.HandlerFunc) func(t *testing.T) {
	return func(t *testing.T) {
		assert.HTTPStatusCode(t, h,
			http.MethodGet, "http://localhost:8080/", nil /*url.Values*/, http.StatusOK)
	}
}
