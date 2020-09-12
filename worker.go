package gocache

import (
	"sync"
)

type worker struct {
	id        int
	done      *sync.WaitGroup
	readyPool chan chan work //get work from the boss
	work      chan work
	quit      chan bool
}

func newWorker(id int, readyPool chan chan work, done *sync.WaitGroup) *worker {
	return &worker{
		id:        id,
		done:      done,
		readyPool: readyPool,
		work:      make(chan work),
		quit:      make(chan bool),
	}
}

func (w *worker) Process(work work) {
	//Do the work
	work.do()
}

func (w *worker) Start() {
	go func() {
		w.done.Add(1) // wait for me
		for {
			w.readyPool <- w.work //hey i am ready to work on new job
			select {
			case work := <-w.work: // hey i am waiting for new job
				w.Process(work) // ok i am on it
			case <-w.quit:
				w.done.Done() // ok i am here i finished my all jobs
				return
			}
		}
	}()
}

func (w *worker) Stop() {
	//tell worker to stop after current process
	w.quit <- true
}
