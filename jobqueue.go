package jobqueue

import (
	"log"
	"sync"
)

type Job interface {
	Todo() error
}

type Worker struct {
	workerPool chan chan *Job
	jobChannel chan *Job
	doneChan   chan bool
	closeChan chan struct{}
}

func newWorker(workerPool chan chan *Job, doneChan chan bool, closeChan chan struct{}) *Worker {
	return &Worker{
		workerPool: workerPool,
		jobChannel: make(chan *Job),
		doneChan:   doneChan,
		closeChan:	closeChan,
	}
}

func (w *Worker) start()  {
	go func() {
		for {
			w.workerPool <- w.jobChannel

			select {
			case job := <- w.jobChannel:
				if err := (*job).Todo(); err != nil {
					log.Printf("%v Todo : %v", job, err)
				}
				w.doneChan <- true
			case <- w.closeChan:
				return
			}
		}
	}()
}


type Factory struct {
	workerPool	chan chan* Job
	jobQueue	chan* Job
	doneChan	chan bool
	closeChan	chan struct{}
	maxWorkers	int
	sync.WaitGroup
}

func NewFactory(maxWorkers int) *Factory {
	workerPool := make(chan chan* Job, maxWorkers)
	jobQueue := make(chan* Job)
	doneChan := make(chan bool)
	closeChan := make(chan struct{})
	return &Factory{
		workerPool: workerPool,
		jobQueue:   jobQueue,
		doneChan:   doneChan,
		maxWorkers: maxWorkers,
		closeChan:	closeChan,
	}
}

func (f *Factory) Run()  {
	for i := 0; i < f.maxWorkers; i++ {
		worker := newWorker(f.workerPool, f.doneChan, f.closeChan)
		worker.start()
	}

	go f.dispatch()
}

func (f *Factory) Close() {
	close(f.closeChan)
}

func (f *Factory) dispatch() {
	for  {
		select {
		case job := <- f.jobQueue:
			go func(job *Job) {
				jobChannel := <- f.workerPool
				jobChannel <- job
			}(job)
		case <-f.doneChan:
			f.Done()
		case <- f.closeChan:
			return
		}
	}
}

func (f *Factory)SetAJob(job *Job)  {
	f.Add(1)
	go func() {
		f.jobQueue <- job
	}()
}
