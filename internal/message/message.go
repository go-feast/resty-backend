package message

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
)

type Event map[string]interface{}

func NewMessage(v Event, marshaler func(v interface{}) ([]byte, error)) *message.Message {
	payload, err := marshaler(v)
	if err != nil {
		panic(err)
	}

	return message.NewMessage(
		uuid.NewString(),
		payload,
	)
}
