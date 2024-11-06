package restaurant

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrder_SetOrderStatus(t *testing.T) {
	f := func(t *testing.T, wantError bool, status OrderStatus, setOrderStatus OrderStatus) {
		o := &Order{OrderStatus: status}
		err := o.SetOrderStatus(setOrderStatus)
		if wantError {
			if err == nil {
				t.Errorf("want error, got nil")
			}
		} else {
			if err != nil {
				t.Errorf("want nil, got %v", err)
			}
		}
	}

	f(t, false, OrderCreated, OrderCreated)
	f(t, false, OrderCreated, OrderCanceled)
	f(t, false, OrderCreated, OrderCompleted)

	f(t, true, OrderCanceled, OrderCreated)
	f(t, false, OrderCanceled, OrderCanceled)
	f(t, true, OrderCanceled, OrderCompleted)

	f(t, true, OrderCompleted, OrderCreated)
	f(t, true, OrderCompleted, OrderCanceled)
	f(t, false, OrderCompleted, OrderCompleted)
}

func TestOrder_SetRestaurantStatus_WithOrderCreated(t *testing.T) {
	f := func(t *testing.T, wantError bool, status RestaurantStatus, setRestaurantStatus RestaurantStatus) {
		o := &Order{RestaurantStatus: status, OrderStatus: OrderCreated}
		err := o.SetRestaurantStatus(setRestaurantStatus)
		if wantError {
			if err == nil {
				t.Errorf("want error, got nil")
			}
		} else {
			if err != nil {
				t.Errorf("want nil, got %v", err)
			}
		}
	}

	f(t, false, ReceivedOrder, ReceivedOrder)
	f(t, false, ReceivedOrder, PreparingOrder)
	f(t, false, ReceivedOrder, PreparedOrder)

	f(t, true, PreparingOrder, ReceivedOrder)
	f(t, false, PreparingOrder, PreparingOrder)
	f(t, false, PreparingOrder, PreparedOrder)

	f(t, true, PreparedOrder, ReceivedOrder)
	f(t, true, PreparedOrder, PreparingOrder)
	f(t, false, PreparedOrder, PreparedOrder)
}

func TestOrder_SetRestaurantStatus_WithOrderArguments(t *testing.T) {
	f := func(t *testing.T, wantError bool, os OrderStatus, status RestaurantStatus, setRestaurantStatus RestaurantStatus) {
		o := &Order{RestaurantStatus: status, OrderStatus: os}
		err := o.SetRestaurantStatus(setRestaurantStatus)
		if wantError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}

	{
		f(t, false, OrderCreated, ReceivedOrder, ReceivedOrder)
		f(t, false, OrderCreated, ReceivedOrder, PreparingOrder)
		f(t, false, OrderCreated, ReceivedOrder, PreparedOrder)

		f(t, false, OrderCanceled, ReceivedOrder, ReceivedOrder)
		f(t, true, OrderCanceled, ReceivedOrder, PreparingOrder)
		f(t, false, OrderCanceled, ReceivedOrder, PreparedOrder)

		f(t, false, OrderCompleted, ReceivedOrder, ReceivedOrder)
		f(t, true, OrderCompleted, ReceivedOrder, PreparingOrder)
		f(t, false, OrderCompleted, ReceivedOrder, PreparedOrder)
	}

	{
		f(t, true, OrderCreated, PreparingOrder, ReceivedOrder)
		f(t, false, OrderCreated, PreparingOrder, PreparingOrder)
		f(t, false, OrderCreated, PreparingOrder, PreparedOrder)

		f(t, true, OrderCanceled, PreparingOrder, ReceivedOrder)
		f(t, true, OrderCanceled, PreparingOrder, PreparingOrder)
		f(t, false, OrderCanceled, PreparingOrder, PreparedOrder)

		f(t, true, OrderCompleted, PreparingOrder, ReceivedOrder)
		f(t, true, OrderCompleted, PreparingOrder, PreparingOrder)
		f(t, false, OrderCompleted, PreparingOrder, PreparedOrder)

	}

	{
		f(t, true, OrderCreated, PreparedOrder, ReceivedOrder)
		f(t, true, OrderCreated, PreparedOrder, PreparingOrder)
		f(t, false, OrderCreated, PreparedOrder, PreparedOrder)

		f(t, true, OrderCanceled, PreparedOrder, ReceivedOrder)
		f(t, true, OrderCanceled, PreparedOrder, PreparingOrder)
		f(t, false, OrderCanceled, PreparedOrder, PreparedOrder)

		f(t, true, OrderCompleted, PreparedOrder, ReceivedOrder)
		f(t, true, OrderCompleted, PreparedOrder, PreparingOrder)
		f(t, false, OrderCompleted, PreparedOrder, PreparedOrder)
	}
}

func TestOrder_IsCanceledOrCompleted(t *testing.T) {
	f := func(t *testing.T, want bool, status OrderStatus) {
		o := &Order{OrderStatus: status}
		assert.Equal(t, want, o.IsCanceledOrCompleted())
	}

	f(t, false, OrderCreated)
	f(t, true, OrderCanceled)
	f(t, true, OrderCompleted)
}
