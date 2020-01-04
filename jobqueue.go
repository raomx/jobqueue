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
}

func newWorker(workerPool chan chan *Job, done chan bool) *Worker {
	return &Worker{
		workerPool: workerPool,
		jobChannel: make(chan *Job),
		doneChan:   done,
	}
}

func (w *Worker) start()  {
	go func() {
		for {
			w.workerPool <- w.jobChannel

			select {
			case job := <- w.jobChannel:
				if err := (*job).Todo(); err != nil {
					log.Fatalf("%v Todo : %v", job, err)
				}
				w.doneChan <- true
			}
		}
	}()
}


type Factory struct {
	workerPool chan chan* Job
	jobQueue   chan* Job
	doneChan   chan bool
	maxWorkers int
	sync.WaitGroup
}

func NewFactory(maxWorkers int) *Factory {
	workerPool := make(chan chan* Job, maxWorkers)
	jobQueue := make(chan* Job)
	doneChan := make(chan bool)
	return &Factory{
		workerPool: workerPool,
		jobQueue:   jobQueue,
		doneChan:   doneChan,
		maxWorkers: maxWorkers,
	}
}

func (f *Factory)Run()  {
	for i := 0; i < f.maxWorkers; i++ {
		worker := newWorker(f.workerPool, f.doneChan)
		worker.start()
	}

	go f.dispatch()
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
		}
	}
}

func (f *Factory)SetAJob(job *Job)  {
	go func() {
		f.Add(1)
		f.jobQueue <- job
	}()
}