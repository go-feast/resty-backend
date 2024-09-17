package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/resty-backend/domain/order"
	"github.com/go-feast/resty-backend/domain/order/event"
	"github.com/pkg/errors"
)

func (h *Handler) OrderClosed(msg *message.Message) error {
	var (
		ctx = msg.Context()
	)

	eventOrderClosed := &event.JSONClosed{}

	err := h.unmarshaler(msg.Payload, eventOrderClosed)
	if err != nil {
		return errors.Wrap(err, "failed to parse order closed event")
	}

	err = h.repository.Operate(ctx, eventOrderClosed.OrderID, func(o *order.Order) error {
		stateOperator := order.NewStateOperator(o)

		closed, closeErr := stateOperator.CloseOrder()
		if closeErr != nil || !closed {
			return errors.Wrapf(closeErr, "can`t set order`s state to closed: order: %s", o.ID())
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to update order closed")
	}

	return nil
}
