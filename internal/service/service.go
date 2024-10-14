package service

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-feast/resty-backend/config"
	"github.com/go-feast/resty-backend/internal/metrics"
	"github.com/go-feast/resty-backend/internal/tracing"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"net"
	"net/http"
	"sync"
	"time"
)

func RunService(ctx context.Context, serviceName, version string, c config.ServiceConfig, routes func(chi.Router)) {
	log.Trace().Str("version", version).Msgf("starting %s", serviceName)

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(version),
		),
		resource.WithProcess(),
		resource.WithOS(),
	)
	if err != nil {
		log.Panic().Err(err).Msg("filed to create resource")
	}

	if err = tracing.RegisterTracerProvider(ctx, res); err != nil {
		log.Panic().Err(err).Msg("failed to register tracer provider")
	}

	metrics.RegisterServiceName(serviceName)

	router := chi.NewRouter()

	middlewares(router)
	routes(router)

	service := &http.Server{
		Addr:         net.JoinHostPort(c.Server.Host, c.Server.Port),
		Handler:      router,
		WriteTimeout: c.Server.WriteTimeout,
		ReadTimeout:  c.Server.ReadTimeout,
		IdleTimeout:  c.Server.IdleTimeout,
	}

	metric := &http.Server{ //nolint:gosec
		Addr:    net.JoinHostPort(c.MetricServer.Host, c.MetricServer.Port),
		Handler: promhttp.Handler(),
	}

	go func() {
		<-ctx.Done()

		closeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if shutdownErr := metric.Shutdown(closeCtx); shutdownErr != nil {
			log.Err(shutdownErr).Msg("failed to shutdown metric server")
		}

		log.Info().Msg("metric server shut down successfully")

		if shutdownErr := service.Shutdown(closeCtx); shutdownErr != nil {
			log.Err(shutdownErr).Msg("failed to shutdown service")
		}

		log.Info().Msg("service shut down successfully")
	}()

	var wg sync.WaitGroup

	wg.Add(2)

	go func(server *http.Server) {
		defer wg.Done()

		log.Info().Str("addr", server.Addr).Msg("starting metric server")

		if metricErr := metric.ListenAndServe(); metricErr != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Err(metricErr).Msg("failed to start metric server")
			return
		}
	}(metric)

	go func(server *http.Server) {
		defer wg.Done()

		log.Info().Str("addr", server.Addr).Msg("starting http server")

		if serviceErr := server.ListenAndServe(); serviceErr != nil && !errors.Is(serviceErr, http.ErrServerClosed) {
			log.Err(serviceErr).Msg("failed to start service")
			return
		}
	}(service)

	wg.Wait()
}

func middlewares(r chi.Router) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(hlog.NewHandler(log.Logger))
	r.Use(middleware.Recoverer)
}
