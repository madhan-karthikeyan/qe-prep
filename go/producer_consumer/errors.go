package producer_consumer

import "errors"

// ErrQueueClosed is returned when attempting to Put to or Get from a closed
// queue.
var ErrQueueClosed = errors.New("queue is closed")
