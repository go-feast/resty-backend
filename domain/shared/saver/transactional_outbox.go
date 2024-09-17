package saver

import "context"

type TransactionalOutbox[T comparable] interface {
	Save(ctx context.Context, entity T) error
}
