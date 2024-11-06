package order

import "testing"

func TestOrder_SetOrderStatus(t *testing.T) {
	f := func(t *testing.T, wantError bool, orderStatus OrderStatus, setOrderStatus OrderStatus) {
		o := Order{OrderStatus: orderStatus}
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

	f(t, false, Created, Created)
	f(t, false, Created, Canceled)
	f(t, false, Created, Completed)

	f(t, true, Canceled, Created)
	f(t, false, Canceled, Canceled)
	f(t, false, Canceled, Completed)

	f(t, true, Completed, Created)
	f(t, true, Completed, Canceled)
	f(t, false, Completed, Completed)
}

func TestOrder_SetCourierStatus(t *testing.T) {
	f := func(t *testing.T, wantError bool, courierStatus CourierStatus, setCourierStatus CourierStatus) {
		o := Order{CourierStatus: courierStatus}
		err := o.SetCourierStatus(setCourierStatus)
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

	f(t, false, CourierTookOrder, CourierTookOrder)
	f(t, false, CourierTookOrder, CourierDelivering)
	f(t, false, CourierTookOrder, CourierDelivered)

	f(t, true, CourierDelivering, CourierTookOrder)
	f(t, false, CourierDelivering, CourierDelivering)
	f(t, false, CourierDelivering, CourierDelivered)

	f(t, true, CourierDelivered, CourierTookOrder)
	f(t, true, CourierDelivered, CourierDelivering)
	f(t, false, CourierDelivered, CourierDelivered)
}

func TestOrder_SetRestaurantStatus(t *testing.T) {
	f := func(t *testing.T, wantError bool, restaurantStatus RestaurantStatus, setRestaurantStatus RestaurantStatus) {
		o := Order{RestaurantStatus: restaurantStatus}
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

	f(t, false, RestaurantReceivedOrder, RestaurantReceivedOrder)
	f(t, false, RestaurantReceivedOrder, RestaurantPreparingOrder)
	f(t, false, RestaurantReceivedOrder, RestaurantPreparedOrder)

	f(t, true, RestaurantPreparingOrder, RestaurantReceivedOrder)
	f(t, false, RestaurantPreparingOrder, RestaurantPreparingOrder)
	f(t, false, RestaurantPreparingOrder, RestaurantPreparedOrder)

	f(t, true, RestaurantPreparedOrder, RestaurantReceivedOrder)
	f(t, true, RestaurantPreparedOrder, RestaurantPreparingOrder)
	f(t, false, RestaurantPreparedOrder, RestaurantPreparedOrder)
}

func TestOrder_SetTransactionStatus(t *testing.T) {
	f := func(t *testing.T, wantError bool, transactionStatus TransactionStatus, setTransactionStatus TransactionStatus) {
		o := Order{PaymentStatus: transactionStatus}
		err := o.SetPaymentStatus(setTransactionStatus)
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

	f(t, false, PaymentWaiting, PaymentWaiting)
	f(t, false, PaymentWaiting, PaymentPaid)
	f(t, false, PaymentWaiting, PaymentCanceled)

	f(t, true, PaymentPaid, PaymentWaiting)
	f(t, false, PaymentPaid, PaymentPaid)
	f(t, true, PaymentPaid, PaymentCanceled)

	f(t, true, PaymentCanceled, PaymentWaiting)
	f(t, true, PaymentCanceled, PaymentPaid)
	f(t, false, PaymentCanceled, PaymentCanceled)
}
