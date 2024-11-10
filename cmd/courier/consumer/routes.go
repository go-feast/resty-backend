package main

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	pubsubcourier "github.com/go-feast/resty-backend/api/pubsub/courier"
	"github.com/go-feast/resty-backend/infrastructure/pubsub"
	gormcourier "github.com/go-feast/resty-backend/infrastructure/repositories/courier/courier"
	"github.com/go-feast/resty-backend/internal/config"
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
	repository := gormcourier.NewGormCourierRepository(db)

	handler := pubsubcourier.NewHandler(repository, json.Marshal, json.Unmarshal)

	r.AddMiddleware(func(h message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			log.Printf("Received msg: %s, payload: %s", msg.UUID, msg.Payload)
			return h(msg)
		}
	})

	r.AddNoPublisherHandler(
		"restaurant.order.preparing",
		"restaurant.order.preparing",
		pubsub.NewKafkaSub("resty.courier"),
		handler.OrderPreparing(),
	)
}
