package core

import (
	"fmt"
	"sync"
)

type Job struct {
	Name string
	Work func() (bool, error)
}

type Worker struct {
	JobQueue chan Job
	Wg       *sync.WaitGroup
}

func NewWorker(jobQueue chan Job, wg *sync.WaitGroup) *Worker {
	return &Worker{JobQueue: jobQueue, Wg: wg}
}

func (w *Worker) Start() {
	go func() {
		for job := range w.JobQueue {
			fmt.Println("Worker: received job")
			w.Wg.Add(1)
			job.Work()
			w.Wg.Done()
		}
	}()
}

type Dispatcher struct {
	WorkerPool chan chan Job
	JobQueue   chan Job
	Wg         *sync.WaitGroup
}

func NewDispatcher(maxWorkers int) *Dispatcher {
	workerPool := make(chan chan Job, maxWorkers)
	jobQueue := make(chan Job)
	wg := &sync.WaitGroup{}

	return &Dispatcher{WorkerPool: workerPool, JobQueue: jobQueue, Wg: wg}
}

func (d *Dispatcher) Run() {
	for i := 0; i < cap(d.WorkerPool); i++ {
		worker := NewWorker(d.JobQueue, d.Wg)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for job := range d.JobQueue {
		fmt.Println("Dispatcher: received job")
		job := job // Create a new variable and assign the value of job to it
		go func() {
			worker := <-d.WorkerPool
			worker <- job
		}()
	}
}
