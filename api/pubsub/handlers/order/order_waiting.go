package order

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/resty-backend/domain/order"
	"github.com/go-feast/resty-backend/domain/order/event"
	"github.com/pkg/errors"
)

func (h *Handler) OrderWaitingForCourier(msg *message.Message) error {
	var (
		ctx = msg.Context()
	)

	eventOrderWaiting := &event.JSONWaitingForCourier{}

	err := h.unmarshaler(msg.Payload, eventOrderWaiting)
	if err != nil {
		return errors.Wrap(err, "failed to parse order waiting event")
	}

	err = h.repository.Operate(ctx, eventOrderWaiting.OrderID, func(o *order.Order) error {
		stateOperator := order.NewStateOperator(o)

		waiting, waitingErr := stateOperator.WaitForCourier()
		if waitingErr != nil || !waiting {
			return errors.Wrapf(waitingErr, "can`t set order`s state to waiting: order: %s", o.ID())
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to update order waiting")
	}

	return nil
}
