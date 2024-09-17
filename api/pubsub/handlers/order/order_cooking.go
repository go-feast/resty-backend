package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/resty-backend/domain/order"
	"github.com/go-feast/resty-backend/domain/order/event"
	"github.com/pkg/errors"
)

func (h *Handler) CookingOrder(msg *message.Message) error {
	var (
		ctx = msg.Context()
	)

	eventOrderCooking := &event.JSONOrderCooking{}

	err := h.unmarshaler(msg.Payload, eventOrderCooking)
	if err != nil {
		return errors.Wrap(err, "failed to parse order cooking event")
	}

	err = h.repository.Operate(ctx, eventOrderCooking.OrderID, func(o *order.Order) error {
		stateOperator := order.NewStateOperator(o)

		cooking, cookErr := stateOperator.CookOrder()
		if cookErr != nil || !cooking {
			return errors.Wrapf(cookErr, "can`t set order`s state to cooking: order: %s", o.ID())
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to update order cooking")
	}

	return nil
}
