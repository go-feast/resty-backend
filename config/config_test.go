package config_test

import (
	"github.com/go-feast/resty-backend/config"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestServerConfig(t *testing.T) {
	testCases := []struct {
		st any
	}{
		{&config.MainServiceServerConfig{}},
		{&config.MetricServerConfig{}},
	}

	for _, testCase := range testCases {
		tc := testCase
		v := reflect.ValueOf(tc.st)

		_, ok := v.Interface().(config.ServerConfig)
		assert.True(t, ok)
	}
}

func TestMainServiceServerConfig(t *testing.T) {
	c := config.MainServiceServerConfig{
		Port:         "80",
		Host:         "127.0.0.1",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	assert.Equal(t, c.HostPort(), "127.0.0.1:80")
	assert.Equal(t, c.WriteTimeoutDur(), 10*time.Second)
	assert.Equal(t, c.ReadTimeoutDur(), 10*time.Second)
	assert.Equal(t, c.IdleTimeoutDur(), 10*time.Second)
}

func TestMetricServerConfig(t *testing.T) {
	c := config.MetricServerConfig{
		Port:         "80",
		Host:         "127.0.0.1",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	assert.Equal(t, c.HostPort(), "127.0.0.1:80")
	assert.Equal(t, c.WriteTimeoutDur(), 10*time.Second)
	assert.Equal(t, c.ReadTimeoutDur(), 10*time.Second)
	assert.Equal(t, c.IdleTimeoutDur(), 10*time.Second)
}
