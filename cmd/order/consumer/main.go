package main

import (
	"context"
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/resty-backend/api/pubsub/handlers/order"
	"github.com/go-feast/resty-backend/config"
	repository "github.com/go-feast/resty-backend/infrastructure/repositories/order/gorm"
	"github.com/go-feast/resty-backend/internal/closer"
	"github.com/go-feast/resty-backend/internal/consumer"
	"github.com/go-feast/resty-backend/internal/logging"
	"github.com/go-feast/resty-backend/internal/pubsub"

	"github.com/go-feast/topics"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"
	"go.opentelemetry.io/otel"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"os/signal"
)

const (
	version     = "v1.0"
	serviceName = "order_consumer"
	driverName  = "pgx/v5"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Kill, os.Interrupt)
	defer stop()

	c := &config.ConsumerConfig{}

	err := envconfig.Process(ctx, c)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to process env vars")
	}

	log.Info().Any("config", c).Send()

	Closer := closer.NewCloser()
	defer Closer.Close()

	db, err := gorm.Open(postgres.New(
		postgres.Config{
			DriverName: driverName,
			DSN:        c.DB.DSN(),
		}), &gorm.Config{})

	if err != nil {
		log.Fatal().Err(err).
			Str("dsn", c.DB.DSN()).
			Str("driver", driverName).
			Msg("failed to connect to database")
	}
	subscriberSQL, err := pubsub.NewSQLSubscriber(db, logging.WatermillLoggerAdapter{Log: log.Logger})
	if err != nil {
		panic(err)
	}

	publisherKafka, err := pubsub.NewKafkaPublisher(c.Kafka.KafkaURL, logging.WatermillLoggerAdapter{Log: log.Logger})
	if err != nil {
		panic(err)
	}

	subscriberKafka, err := pubsub.NewKafkaSubscriber(c.Kafka.KafkaURL, logging.WatermillLoggerAdapter{Log: log.Logger})
	if err != nil {
		panic(err)
	}
	var (
		orderRepository = repository.NewOrderRepository(db)
	)

	var (
		handler = order.NewHandler(
			json.Unmarshal,
			otel.GetTracerProvider().Tracer(serviceName),
			orderRepository,
		)
	)

	consumer.RunConsumer(ctx, serviceName, version, *c, func(r *message.Router) {
		r.AddHandler(
			"handler.order.created",
			topics.OrderCreated.String(),
			subscriberSQL,
			topics.OrderCreated.String(),
			publisherKafka,
			handler.OrderCreated,
		)

		r.AddNoPublisherHandler(
			"handler.order.paid",
			topics.Paid.String(),
			subscriberKafka,
			handler.OrderPaid,
		)

		r.AddNoPublisherHandler(
			"order.cooking",
			topics.Cooking.String(),
			subscriberKafka,
			handler.CookingOrder,
		)

		r.AddNoPublisherHandler(
			"order.cooking.finished",
			topics.CookingFinished.String(),
			subscriberKafka,
			handler.FinishedCooking,
		)

		r.AddNoPublisherHandler(
			"order.waiting",
			topics.WaitingForCourier.String(),
			subscriberKafka,
			handler.OrderWaitingForCourier,
		)

		r.AddNoPublisherHandler(
			"order.taken",
			topics.CourierTook.String(),
			subscriberKafka,
			handler.CookingTaken,
		)

		r.AddNoPublisherHandler(
			"order.delivering",
			topics.Delivering.String(),
			subscriberKafka,
			handler.OrderDelivering,
		)

		r.AddNoPublisherHandler(
			"order.delivered",
			topics.Delivered.String(),
			subscriberKafka,
			handler.OrderDelivered,
		)

		r.AddNoPublisherHandler(
			"order.closed",
			topics.Closed.String(),
			subscriberKafka,
			handler.OrderClosed,
		)

		r.AddNoPublisherHandler(
			"order.canceled",
			topics.Canceled.String(),
			subscriberKafka,
			handler.OrderCanceled,
		)
	})
}
