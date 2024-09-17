package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/resty-backend/domain/order"
	"github.com/go-feast/resty-backend/domain/order/event"
	"github.com/pkg/errors"
)

func (h *Handler) OrderDelivered(msg *message.Message) error {
	var (
		ctx = msg.Context()
	)

	eventOrderDelivered := &event.JSONDelivered{}

	err := h.unmarshaler(msg.Payload, eventOrderDelivered)
	if err != nil {
		return errors.Wrap(err, "failed to parse order delivered event")
	}

	err = h.repository.Operate(ctx, eventOrderDelivered.OrderID, func(o *order.Order) error {
		stateOperator := order.NewStateOperator(o)

		delivered, deliveredErr := stateOperator.OrderDelivered()
		if deliveredErr != nil || !delivered {
			return errors.Wrapf(deliveredErr, "can`t set order`s state to delivered: order: %s", o.ID())
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to update order delivered")
	}

	return nil
}
