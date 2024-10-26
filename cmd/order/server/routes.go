package main

import (
	"github.com/gin-gonic/gin"
	apiorder "github.com/go-feast/resty-backend/api/http/order"
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
	handler := apiorder.NewHandler(orderRepository)

	e.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	v1 := e.Group("/api/v1")

	{
		v1.POST("/orders", handler.TakeOrder())
	}
}
