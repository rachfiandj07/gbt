package pool_worker

import "sync"

type PoolWorker interface {
	Start(maxConcurrentWorker int)
}
type poolWorker struct {
	task                chan interface{}
	taskLength          int
	channel             interface{}
	processFunc         func()
	maxConcurrentWorker int
	wg                  sync.WaitGroup
}

func NewPoolWorker(task chan interface{}, processFunc func(), maxConcurrentWorker int) PoolWorker {
	return &poolWorker{
		task:                task,
		processFunc:         processFunc,
		maxConcurrentWorker: maxConcurrentWorker,
	}
}

func (pw *poolWorker) Start(maxConcurrentWorker int) {
	pw.channel = make(chan interface{}, pw.maxConcurrentWorker)
	pw.wg = sync.WaitGroup{}
	for i := 0; i < maxConcurrentWorker; i++ {
		pw.wg.Add(1)
		go func() {
			defer pw.wg.Done()
			for _ = range pw.task {
				pw.processFunc()
			}
		}()
	}

	for task := 0; task < pw.taskLength; task++ {
		pw.task <- task
	}

	pw.wg.Done()
}
