package order

import (
	"context"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/go-feast/resty-backend/domain/order"
	"github.com/go-feast/topics"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Outbox struct {
	publisher  message.Publisher
	marshaller func(any) ([]byte, error)
	repository order.OrderRepository
}

func NewOutbox(
	publisher message.Publisher,
	repository order.OrderRepository,
	marshaller func(any) ([]byte, error),
) *Outbox {
	return &Outbox{
		publisher:  publisher,
		marshaller: marshaller,
		repository: repository,
	}
}

func (ob *Outbox) Save(
	ctx context.Context,
	o *order.Order,
) (err error) {
	var createOrderErr, marshallEventErr, publishErr error

	defer func() {
		if marshallEventErr != nil || publishErr != nil {
			err = ob.repository.Delete(ctx, o)
		}
	}()

	createOrderErr = ob.repository.Create(ctx, o)
	if createOrderErr != nil {
		return errors.Wrap(createOrderErr, "outbox: failed to create order")
	}

	bytes, marshallEventErr := ob.marshaller(o.ToEvent().JSONEventOrderCreated())
	if marshallEventErr != nil {
		return errors.Wrap(marshallEventErr, "outbox: saving: failed to marshal event")
	}

	msg := message.NewMessage(uuid.NewString(), bytes)

	msg.SetContext(ctx)

	publishErr = ob.publisher.Publish(topics.OrderCreated.String(), msg)
	if publishErr != nil {
		return errors.Wrap(publishErr, "outbox: saving: failed to publish event")
	}

	return nil
}
