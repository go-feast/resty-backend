package consumer

import (
	"context"
	"errors"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/resty-backend/config"
	"github.com/go-feast/resty-backend/internal/logging"
	"github.com/go-feast/resty-backend/internal/metrics"
	"github.com/go-feast/resty-backend/internal/tracing"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"net"
	"net/http"
	"sync"
	"time"
)

func RunConsumer(ctx context.Context, serviceName, version string, c config.ConsumerConfig, routes func(router *message.Router)) {
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
		log.Err(err).Msg("filed to create resource")
		return
	}

	if err = tracing.RegisterTracerProvider(ctx, res); err != nil {
		log.Err(err).Msg("failed to register tracer provider")
		return
	}

	metrics.RegisterServiceName(serviceName)

	metricServer := &http.Server{
		Addr:    net.JoinHostPort(c.MetricServer.Host, c.MetricServer.Port),
		Handler: promhttp.Handler(),
	}

	router, err := message.NewRouter(
		message.RouterConfig{},
		&logging.WatermillLoggerAdapter{Log: log.Logger},
	)
	if err != nil {
		return
	}

	routes(router)

	var wg sync.WaitGroup
	wg.Add(2) // 2 servers: metric and consumer

	go func() {
		defer wg.Done()

		if e := router.Run(ctx); e != nil {
			log.Error().Err(err).Msgf("failed to run consumer: %s", e.Error())
			return
		}

		log.Info().Msg("exiting consumer")
	}()
	go func() {
		defer wg.Done()
		log.Info().Msgf("listening on %s", metricServer.Addr)
		e := metricServer.ListenAndServe()
		if e != nil &&
			!errors.Is(e, http.ErrServerClosed) {
			log.Error().Err(e).Msgf("failed to start metrics server: %s", e.Error())
			return
		}
		log.Info().Msg("closed metrics server")
	}()

	go func() {
		<-ctx.Done()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if e := metricServer.Shutdown(shutdownCtx); e != nil {
			log.Error().Err(e).Msgf("failed to shutdown metrics server: %s", e.Error())
			return
		}
		log.Info().Msg("closed metrics server")
	}()
	go func() {
		<-ctx.Done()
		if router.IsRunning() {
			if e := router.Close(); e != nil {
				log.Error().Err(err).Msgf("failed to close consumer: %s", e.Error())
			}
			log.Info().Msg("closed consumer")
		}
	}()

	wg.Wait()
}
