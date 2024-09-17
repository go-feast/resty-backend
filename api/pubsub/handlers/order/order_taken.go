package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/resty-backend/domain/order"
	"github.com/go-feast/resty-backend/domain/order/event"
	"github.com/pkg/errors"
)

func (h *Handler) CookingTaken(msg *message.Message) error {
	var (
		ctx = msg.Context()
	)

	eventOrderTaken := &event.JSONCourierTook{}

	err := h.unmarshaler(msg.Payload, eventOrderTaken)
	if err != nil {
		return errors.Wrap(err, "failed to parse order taken event")
	}

	err = h.repository.Operate(ctx, eventOrderTaken.OrderID, func(o *order.Order) error {
		stateOperator := order.NewStateOperator(o)

		taken, takingErr := stateOperator.CourierTookOrder(eventOrderTaken.CourierID)
		if takingErr != nil || !taken {
			return errors.Wrapf(takingErr, "can`t set order`s state to taken: order: %s", o.ID())
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to update order taken")
	}

	return nil
}
