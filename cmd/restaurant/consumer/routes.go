package main

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	pubsubrest "github.com/go-feast/resty-backend/api/pubsub/restaurant"
	"github.com/go-feast/resty-backend/infrastructure/pubsub"
	gormorder "github.com/go-feast/resty-backend/infrastructure/repositories/restaurant/order"
	"github.com/go-feast/resty-backend/internal/config"
	"github.com/go-feast/resty-backend/internal/domain/restaurant"
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

	orderRepository := gormorder.NewGormRepository(db)

	handler := pubsubrest.NewHandler(orderRepository, json.Unmarshal, json.Marshal)
	r.AddMiddleware(func(h message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			log.Printf("Received message: %s, Payload: %s", msg.UUID, msg.Payload)
			return h(msg)
		}
	})

	r.AddMiddleware(func(h message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			v, err := h(msg)
			log.Printf("Posted message: %vz, error: %s", v, err)
			return v, err
		}
	})
	sub := pubsub.NewKafkaSub("resty.restaurant")
	r.AddHandler(
		"payment.paid",
		"payment.paid",
		sub,
		restaurant.ReceivedOrder.String(),
		pubsub.NewKafkaPub(),
		handler.ReceivedOrder(),
	)

	r.AddNoPublisherHandler(
		restaurant.OrderCanceled.String(),
		restaurant.OrderCanceled.String(),
		sub,
		handler.OrderCanceled(),
	)

	r.AddNoPublisherHandler(
		restaurant.OrderCompleted.String(),
		restaurant.OrderCompleted.String(),
		sub,
		handler.OrderCompleted(),
	)
}
