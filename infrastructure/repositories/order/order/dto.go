package order

import (
	domainorder "github.com/go-feast/resty-backend/domain/order"
	"github.com/go-feast/resty-backend/domain/shared/geo"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

func InitializeOrderScheme(db *gorm.DB) {
	err := db.AutoMigrate(&DatabaseOrderDTO{})
	if err != nil {
		panic(errors.Wrap(err, "failed to migrate database"))
	}
}

type DatabaseOrderDTO struct { //nolint:govet
	ID            uuid.UUID `order:"type:uuid;primaryKey"`
	RestaurantID  uuid.UUID `order:"type:uuid"`
	CustomerID    uuid.UUID `order:"type:uuid"`
	CourierID     uuid.UUID `order:"type:uuid"`
	State         string    `order:"type:text"`
	TransactionID uuid.UUID `order:"type:uuid"`
	Latitude      float64   `order:"type:numeric"`
	Longitude     float64   `order:"type:numeric"`
	CreatedAt     time.Time
	Meals         []uuid.UUID `order:"type:uuid"`
}

func toOrder(d *DatabaseOrderDTO) *domainorder.Order {
	dst, _ := geo.NewDestination(d.Latitude, d.Longitude)
	return domainorder.NewOrderFactory().
		WithID(d.ID).
		WithRestaurantID(d.RestaurantID).
		WithCustomerID(d.CustomerID).
		WithCourierID(d.CourierID).
		WithMeals(d.Meals).
		WithState(d.State).
		WithTransactionID(d.TransactionID).
		WithDestination(dst).
		WithCreatedAt(d.CreatedAt).
		Build()
}

func toDatabaseDTO(o *domainorder.Order) *DatabaseOrderDTO {
	return &DatabaseOrderDTO{
		ID:            o.ID(),
		RestaurantID:  o.RestaurantID(),
		CustomerID:    o.CustomerID(),
		CourierID:     o.CourierID(),
		State:         o.State().String(),
		TransactionID: o.TransactionID(),
		Latitude:      o.Destination().Latitude(),
		Longitude:     o.Destination().Longitude(),
		CreatedAt:     o.CreateAt(),
		Meals:         o.Meals(),
	}
}
