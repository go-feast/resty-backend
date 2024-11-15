package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	apiorder "github.com/go-feast/resty-backend/api/http/order"
	"github.com/go-feast/resty-backend/infrastructure/pubsub"
	gormorder "github.com/go-feast/resty-backend/infrastructure/repositories/order"
	"github.com/go-feast/resty-backend/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func routes(e *gin.Engine) {
	dsn := config.DBConn()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	gormorder.InitializerOrderOrDie(db)

	orderRepository := gormorder.NewGormOrderRepository(db)
	handler := apiorder.NewHandler(orderRepository, pubsub.NewKafkaPub(), json.Marshal)

	e.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	v1 := e.Group("/api/v1")

	{
		orders := v1.Group("orders")
		orders.POST("/", handler.TakeOrder())
		orders.GET("/:id", handler.GetOrder())
		orders.POST("/:id/cancel", handler.CancelOrder())
		orders.POST("/:id/complete", handler.CloseOrder())
	}
}
