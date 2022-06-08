package multiple

import (
	"sync"
)

type Dispatcher interface {
	Dispatch(producer Producer, worker Worker)
}

type defaultDispatcher struct {
	workerCount int
}

func DefaultDispatcher(workerCount int) Dispatcher {
	return &defaultDispatcher{
		workerCount,
	}
}

func (d *defaultDispatcher) Dispatch(producer Producer, worker Worker) {
	var workerCount = d.workerCount
	var producers []Producer
	if producer.CanShard() {
		producers = producer.Shard()
	} else {
		producers = []Producer{producer}
	}

	// init Queue
	dataQueue := make(chan interface{}, workerCount)
	doneQueue := make(chan bool)

	var WWG sync.WaitGroup
	var PWG sync.WaitGroup
	var count = 0
	// start workers
	WWG.Add(1)
	go func() {
		exit := false
		done := false
		tokenQueue := make(chan int, workerCount)
		for i := 0; i < workerCount; i++ {
			tokenQueue <- i
		}

		data := <-dataQueue
		for !exit {
			<-tokenQueue
			count++
			go func(d interface{}) {
				defer func() { tokenQueue <- 1 }()
				worker.Apply(d)
			}(data)

			if done {
				select {
				case data = <-dataQueue:
				default:
					exit = true
				}
			} else {
				select {
				case data = <-dataQueue:
				case done = <-doneQueue:
					select {
					case data = <-dataQueue:
					default:
						exit = true
					}
				}
			}
		}

		// wait all work done
		for i := 0; i < workerCount; i++ {
			<-tokenQueue
		}
		WWG.Done()
	}()

	// start producer
	PWG.Add(len(producers))
	for i := range producers {
		p := producers[i]
		go func(id int) {
			for !p.IsDone() {
				data := p.Yield()
				dataQueue <- data
			}
			PWG.Done()
		}(i)
	}

	// wait Producer Done
	PWG.Wait()
	doneQueue <- true

	// wait worker down
	WWG.Wait()
}
