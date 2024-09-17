package metrics_test

import (
	"github.com/go-feast/resty-backend/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewCounter(t *testing.T) {
	t.Run("assert NewCounter registers in handler", func(t *testing.T) {
		// register in default registry
		_ = metrics.NewCounter("test", "new_counter")

		handler := promhttp.Handler()

		assert.HTTPBodyContains(t, handler.ServeHTTP, http.MethodGet,
			"/", nil, `test_new_counter 0`)
	})
}
