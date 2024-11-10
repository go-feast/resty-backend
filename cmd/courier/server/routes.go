package main

import (
	"github.com/gin-gonic/gin"
	apicourier "github.com/go-feast/resty-backend/api/http/courier"
	"github.com/go-feast/resty-backend/infrastructure/pubsub"
	gormcourier "github.com/go-feast/resty-backend/infrastructure/repositories/courier/courier"
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

	gormcourier.InitializeCourierOrDie(db)

	courierRepository := gormcourier.NewGormCourierRepository(db)
	handler := apicourier.NewHandler(courierRepository, pubsub.NewKafkaPub())

	e.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	v1 := e.Group("/api/v1")
	{
		couriers := v1.Group("/couriers")
		couriers.POST("/", handler.CreateCourier())
		couriers.GET("/:id", handler.GetCourier())
		couriers.POST("/:cid/:oid", handler.AssignCourier())

		orders := v1.Group("/orders/:id")
		orders.GET("/", handler.GetOrder())
		orders.POST("/took", handler.TookOrder())
		orders.POST("/delivering", handler.DeliveringOrder())
		orders.POST("/delivered", handler.DeliveredOrder())
	}
}
