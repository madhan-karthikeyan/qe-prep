package producer_consumer

import (
	"context"
)

type BlockingQueue[T any] struct {
	ch     chan T
	closed chan struct{}
}

func NewBlockingQueue[T any](capacity int) *BlockingQueue[T] {
	return &BlockingQueue[T]{
		ch:     make(chan T, capacity),
		closed: make(chan struct{}),
	}
}

func (q *BlockingQueue[T]) Put(ctx context.Context, item T) error {
	select {
	case <-q.closed:
		return ErrQueueClosed
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	select {
	case q.ch <- item:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (q *BlockingQueue[T]) Get(ctx context.Context) (T, error) {
	var zero T
	select {
	case <-ctx.Done():
		return zero, ctx.Err()
	default:
	}
	select {
	case item, ok := <-q.ch:
		if !ok {
			return zero, ErrQueueClosed
		}
		return item, nil
	case <-q.closed:
		select {
		case item, ok := <-q.ch:
			if !ok {
				return zero, ErrQueueClosed
			}
			return item, nil
		default:
			return zero, ErrQueueClosed
		}
	case <-ctx.Done():
		return zero, ctx.Err()
	}
}

func (q *BlockingQueue[T]) Close() {
	close(q.closed)
}
