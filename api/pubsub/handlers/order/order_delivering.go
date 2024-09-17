package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/resty-backend/domain/order"
	"github.com/go-feast/resty-backend/domain/order/event"
	"github.com/pkg/errors"
)

func (h *Handler) OrderDelivering(msg *message.Message) error {
	var (
		ctx = msg.Context()
	)

	eventOrderDelivering := &event.JSONDelivering{}

	err := h.unmarshaler(msg.Payload, eventOrderDelivering)
	if err != nil {
		return errors.Wrap(err, "failed to parse order delivering event")
	}

	err = h.repository.Operate(ctx, eventOrderDelivering.OrderID, func(o *order.Order) error {
		stateOperator := order.NewStateOperator(o)

		delivering, deliveryErr := stateOperator.DeliveringOrder()
		if deliveryErr != nil || !delivering {
			return errors.Wrapf(deliveryErr, "can`t set order`s state to delivering: order: %s", o.ID())
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to update order delivering")
	}

	return nil
}
