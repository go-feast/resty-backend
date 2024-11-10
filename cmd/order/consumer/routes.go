package main

import (
	"encoding/json"
	"github.com/ThreeDotsLabs/watermill/message"
	pubsuborder "github.com/go-feast/resty-backend/api/pubsub/order"
	"github.com/go-feast/resty-backend/infrastructure/pubsub"
	gormorder "github.com/go-feast/resty-backend/infrastructure/repositories/order"
	"github.com/go-feast/resty-backend/internal/config"
	"github.com/go-feast/resty-backend/internal/domain/order"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func routes(router *message.Router) {
	dsn := config.DBConn()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	repository := gormorder.NewGormOrderRepository(db)

	handler := pubsuborder.NewHandler(repository, json.Unmarshal, json.Marshal)

	router.AddMiddleware(func(h message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			log.Printf("Received message: %s, Payload: %s", msg.UUID, msg.Payload)
			return h(msg)
		}
	})

	sub := pubsub.NewKafkaSub("resty.order")
	formatter := routeFormatter{sub: sub, r: router, pub: pubsub.NewKafkaPub()}
	formatter.FormatNoPub(order.Created.String(), handler.OrderCreated())
	formatter.FormatNoPub(order.Canceled.String(), handler.OrderCanceled())
	formatter.FormatNoPub(order.Completed.String(), handler.OrderClosed())

	formatter.FormatNoPub(order.RestaurantReceivedOrder.String(), handler.RestaurantReceivedOrder())
	formatter.FormatNoPub(order.RestaurantPreparingOrder.String(), handler.RestaurantPreparingOrder())
	formatter.FormatNoPub(order.RestaurantPreparedOrder.String(), handler.RestaurantPreparedOrder())

	formatter.FormatNoPub(order.PaymentWaiting.String(), handler.PaymentWaiting())
	formatter.FormatNoPub(order.PaymentPaid.String(), handler.PaymentPaid())
	formatter.FormatNoPub(order.PaymentCanceled.String(), handler.PaymentCanceled())

	formatter.FormatNoPub(order.CourierAssigned.String(), handler.AssignedOrder())
	formatter.FormatNoPub(order.CourierTookOrder.String(), handler.CourierTookOrder())
	formatter.FormatNoPub(order.CourierDelivering.String(), handler.CourierDelivering())
	formatter.Format(order.CourierDelivered.String(), order.Completed.String(), handler.CourierDelivered())

}

type routeFormatter struct {
	sub message.Subscriber
	pub message.Publisher
	r   *message.Router
}

func (f *routeFormatter) FormatNoPub(topic string, handlerFunc message.NoPublishHandlerFunc) {
	f.r.AddNoPublisherHandler(topic, topic, f.sub, handlerFunc)
}

func (f *routeFormatter) Format(topic, pubtopic string, handlerFunc message.HandlerFunc) {
	f.r.AddHandler(topic, topic, f.sub, pubtopic, f.pub, handlerFunc)
}
