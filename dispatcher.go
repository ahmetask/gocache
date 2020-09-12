package gocache

import (
	"log"
	"sync"
)

type workerPool struct {
	internalQueue     chan work
	readyPool         chan chan work //boss says hey i have a new job at my desk workers who available can get it in this way he does not have to ask current status of workers
	workers           []*worker
	dispatcherStopped sync.WaitGroup
	workersStopped    *sync.WaitGroup
	quit              chan bool
}

func newWorkerPool(maxWorkers int) *workerPool {
	workersStopped := sync.WaitGroup{}

	readyPool := make(chan chan work, maxWorkers)
	workers := make([]*worker, maxWorkers, maxWorkers)

	// create workers
	for i := 0; i < maxWorkers; i++ {
		workers[i] = newWorker(i+1, readyPool, &workersStopped)
	}

	return &workerPool{
		internalQueue:     make(chan work),
		readyPool:         readyPool,
		workers:           workers,
		dispatcherStopped: sync.WaitGroup{},
		workersStopped:    &workersStopped,
		quit:              make(chan bool),
	}
}

func (q *workerPool) Start() {
	//tell workers to get ready
	for i := 0; i < len(q.workers); i++ {
		q.workers[i].Start()
	}
	// open factory
	go q.dispatch()
}

func (q *workerPool) Stop() {
	q.quit <- true
	q.dispatcherStopped.Wait()
	log.Println("Stopping Working Pool")
}

func (q *workerPool) dispatch() {
	//open factory gate
	q.dispatcherStopped.Add(1)
	for {
		select {
		case job := <-q.internalQueue:
			workerXChannel := <-q.readyPool //free worker x founded
			workerXChannel <- job           // here is your job worker x
		case <-q.quit:
			// free all workers
			for i := 0; i < len(q.workers); i++ {
				q.workers[i].Stop()
			}
			// wait for all workers to finish their job
			q.workersStopped.Wait()
			//close factory gate
			q.dispatcherStopped.Done()
			return
		}
	}
}

func (q *workerPool) Submit(job work) {
	// daily - fill the board with new works
	q.internalQueue <- job
}
