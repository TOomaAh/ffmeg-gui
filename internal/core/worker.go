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

func (w *Worker) Start(ready *sync.WaitGroup) {
	ready.Done() // Indiquer que le Worker est prêt
	go func() {
		for job := range w.JobQueue {
			fmt.Println("Worker: received job")
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

func (d *Dispatcher) AddJob(job Job) {
	d.Wg.Add(1) // Ajouter 1 au WaitGroup pour chaque job
	d.JobQueue <- job
}

func (d *Dispatcher) Close() {
	close(d.JobQueue)
}

func (d *Dispatcher) Run() {
	ready := &sync.WaitGroup{}
	ready.Add(cap(d.WorkerPool)) // Ajouter le nombre de Workers au WaitGroup

	for i := 0; i < cap(d.WorkerPool); i++ {
		worker := NewWorker(d.JobQueue, d.Wg)
		worker.Start(ready)
		d.WorkerPool <- worker.JobQueue
	}

	for {
		select {
		case job, ok := <-d.JobQueue:
			if ok {
				fmt.Println("Dispatcher: received job")
				jobChannel := <-d.WorkerPool
				jobChannel <- job
			}
		default:
			// Si la JobQueue est vide, sortir de la boucle
			if len(d.JobQueue) == 0 {
				return
			}
		}
	}
}

func (d *Dispatcher) dispatch() {
	for job := range d.JobQueue {
		fmt.Println("Dispatcher: received job")
		go func(job Job) { // Passer job en tant que paramètre
			jobChannel := <-d.WorkerPool
			jobChannel <- job
		}(job)
	}
}
