package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/resty-backend/domain/order"
	"github.com/go-feast/resty-backend/domain/order/event"
	"github.com/pkg/errors"
)

func (h *Handler) FinishedCooking(msg *message.Message) error {
	var (
		ctx = msg.Context()
	)

	eventOrderFinished := &event.JSONOrderFinished{}

	err := h.unmarshaler(msg.Payload, eventOrderFinished)
	if err != nil {
		return errors.Wrap(err, "failed to parse order finished cooking event")
	}

	err = h.repository.Operate(ctx, eventOrderFinished.OrderID, func(o *order.Order) error {
		stateOperator := order.NewStateOperator(o)

		finishedCooking, finishedErr := stateOperator.OrderFinished()
		if finishedErr != nil || !finishedCooking {
			return errors.Wrapf(finishedErr, "can`t set order`s state to finished cooking: order: %s", o.ID())
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to update order finished cooking")
	}

	return nil
}
