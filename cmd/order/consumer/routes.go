package main

import (
	"github.com/ThreeDotsLabs/watermill/message"
	gormorder "github.com/go-feast/resty-backend/infrastructure/repositories/order"
	"github.com/go-feast/resty-backend/internal/config"
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

	gormorder.InitializerOrderOrDie(db)

}
