package payment

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPayment_IsOrderCanceledOrCompleted(t *testing.T) {
	o := Payment{OrderStatus: OrderCanceled}
	assert.True(t, o.IsOrderCanceledOrCompleted())
	o = Payment{OrderStatus: OrderCompleted}
	assert.True(t, o.IsOrderCanceledOrCompleted())
	o = Payment{OrderStatus: OrderCreated}
	assert.False(t, o.IsOrderCanceledOrCompleted())
}

func TestPayment_SetPaymentStatus(t *testing.T) {
	f := func(t *testing.T, wantError bool, transactionStatus PaymentStatus, setTransactionStatus PaymentStatus) {
		o := Payment{PaymentStatus: transactionStatus}
		err := o.SetPaymentStatus(setTransactionStatus)
		if wantError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}

	f(t, false, Waiting, Waiting)
	f(t, false, Waiting, Paid)
	f(t, false, Waiting, Canceled)

	f(t, true, Paid, Waiting)
	f(t, false, Paid, Paid)
	f(t, false, Paid, Canceled)

	f(t, true, Canceled, Waiting)
	f(t, false, Canceled, Paid)
	f(t, false, Canceled, Canceled)
}

func TestPayment_SetPaymentStatus_WithOrderCanceled(t *testing.T) {
	f := func(t *testing.T, wantError bool, transactionStatus PaymentStatus, setTransactionStatus PaymentStatus) {
		o := Payment{PaymentStatus: transactionStatus, OrderStatus: OrderCanceled}
		err := o.SetPaymentStatus(setTransactionStatus)
		if wantError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}

	f(t, false, Waiting, Waiting)
	f(t, true, Waiting, Paid)
	f(t, true, Waiting, Canceled)

	f(t, true, Paid, Waiting)
	f(t, true, Paid, Paid)
	f(t, true, Paid, Canceled)

	f(t, true, Canceled, Waiting)
	f(t, true, Canceled, Paid)
	f(t, true, Canceled, Canceled)
}

func TestPayment_SetPaymentStatus_WithOrderCompleted(t *testing.T) {
	f := func(t *testing.T, wantError bool, transactionStatus PaymentStatus, setTransactionStatus PaymentStatus) {
		o := Payment{PaymentStatus: transactionStatus, OrderStatus: OrderCompleted}
		err := o.SetPaymentStatus(setTransactionStatus)
		if wantError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}

	f(t, false, Waiting, Waiting)
	f(t, true, Waiting, Paid)
	f(t, true, Waiting, Canceled)

	f(t, true, Paid, Waiting)
	f(t, true, Paid, Paid)
	f(t, true, Paid, Canceled)

	f(t, true, Canceled, Waiting)
	f(t, true, Canceled, Paid)
	f(t, true, Canceled, Canceled)
}
