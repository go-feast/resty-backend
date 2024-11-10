package main

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	pubsubpay "github.com/go-feast/resty-backend/api/pubsub/payment"
	"github.com/go-feast/resty-backend/infrastructure/pubsub"
	gormpay "github.com/go-feast/resty-backend/infrastructure/repositories/payment"
	"github.com/go-feast/resty-backend/internal/config"
	"github.com/go-feast/resty-backend/internal/domain/payment"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func routes(r *message.Router) {
	dsn := config.DBConn()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	r.AddMiddleware(func(h message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			log.Printf("Received message: %s, Payload: %s", msg.UUID, msg.Payload)
			return h(msg)
		}
	})

	repository := gormpay.NewGormPaymentRepository(db)
	handler := pubsubpay.NewHandler(repository, json.Unmarshal, json.Marshal)
	sub := pubsub.NewKafkaSub("resty.payment")
	r.AddHandler(
		payment.OrderCreated,
		payment.OrderCreated,
		sub,
		payment.Waiting.String(),
		pubsub.NewKafkaPub(),
		handler.CreatePayment(),
	)

	r.AddNoPublisherHandler(
		payment.OrderCanceled,
		payment.OrderCanceled,
		sub,
		handler.OrderCanceled(),
	)

	r.AddNoPublisherHandler(
		payment.OrderCompleted,
		payment.OrderCompleted,
		sub,
		handler.OrderCompleted(),
	)
}
