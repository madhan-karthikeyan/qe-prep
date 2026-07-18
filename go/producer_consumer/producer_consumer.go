package producer_consumer

import (
	"context"
	"sync"
	"sync/atomic"
)

// RunDemo starts numProducers producers and numConsumers consumers using the
// given queue. It blocks until all producers are done and all items are
// consumed.
func RunDemo[T any](
	queue *BlockingQueue[T],
	numProducers, numConsumers int,
	produce func(idx int) (T, bool),
	consume func(idx int, item T),
) {
	ctx := context.Background()
	var producersDone sync.WaitGroup
	producersDone.Add(numProducers)
	for i := range numProducers {
		go func(idx int) {
			defer producersDone.Done()
			for {
				item, ok := produce(idx)
				if !ok {
					return
				}
				if err := queue.Put(ctx, item); err != nil {
					return
				}
			}
		}(i)
	}

	var consumersDone sync.WaitGroup
	consumersDone.Add(numConsumers)
	for i := range numConsumers {
		go func(idx int) {
			defer consumersDone.Done()
			for {
				item, err := queue.Get(ctx)
				if err != nil {
					return
				}
				consume(idx, item)
			}
		}(i)
	}

	producersDone.Wait()
	queue.Close()

	for {
		if _, err := queue.Get(ctx); err != nil {
			break
		}
	}
	consumersDone.Wait()
}

// IntProducer returns a producer function that yields integers from 0 to n-1.
func IntProducer(n int) func(int) (int, bool) {
	var counter atomic.Int64
	return func(int) (int, bool) {
		v := int(counter.Add(1) - 1)
		if v >= n {
			return 0, false
		}
		return v, true
	}
}
