package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	httprestaurant "github.com/go-feast/resty-backend/api/http/restaurant"
	"github.com/go-feast/resty-backend/infrastructure/pubsub"
	gormorder "github.com/go-feast/resty-backend/infrastructure/repositories/restaurant/order"
	gormrestaurant "github.com/go-feast/resty-backend/infrastructure/repositories/restaurant/restaurnt"
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

	gormrestaurant.InitializeRestaurantOrDie(db)

	restaurantRepository := gormrestaurant.NewGormRestaurantRepository(db)
	orderRepository := gormorder.NewGormRepository(db)
	handler := httprestaurant.NewHandler(restaurantRepository, pubsub.NewKafkaPub(), orderRepository, json.Marshal, json.Unmarshal)

	e.GET("/health", func(c *gin.Context) { c.Status(http.StatusOK) })

	v1 := e.Group("/api/v1")
	{
		v1.POST("/restaurants", handler.CreateRestaurant())
		v1.GET("/restaurants/:id", handler.GetRestaurant())
		v1.POST("/orders/:id/preparing", handler.SetPreparingOrder())
		v1.POST("/orders/:id/prepared", handler.SetPreparedOrder())
	}
}
