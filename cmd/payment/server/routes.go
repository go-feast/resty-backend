package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	httppay "github.com/go-feast/resty-backend/api/http/payment"
	"github.com/go-feast/resty-backend/infrastructure/pubsub"
	gormpay "github.com/go-feast/resty-backend/infrastructure/repositories/payment"
	"github.com/go-feast/resty-backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func routes(r *gin.Engine) {
	dsn := config.DBConn()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	gormpay.InitializePaymentOrDie(db)

	payRepository := gormpay.NewGormPaymentRepository(db)
	handler := httppay.NewHandler(payRepository, &pubsub.NopPublisher{}, json.Unmarshal, json.Marshal)

	r.GET("/health", func(c *gin.Context) { c.Status(200) })

	v1 := r.Group("/api/v1")
	{
		payments := v1.Group("/payments/:id/")

		payments.POST("/pay", handler.PayForOrder())
		payments.POST("/cancel", handler.CancelPaymentForOrder())
	}
}
