package workerpool

import (
	"log"
	"sync"
)

type WorkerPool struct {
	workers  []*Worker
	taskChan chan<- func()
	wg       *sync.WaitGroup
}

func NewWorkerPool(workerCnt int) *WorkerPool {
	taskChan := make(chan func())
	workers := make([]*Worker, workerCnt)
	for i := range workerCnt {
		workers[i] = NewWorker(i, taskChan)
	}
	return &WorkerPool{
		workers:  workers,
		taskChan: taskChan,
		wg:       &sync.WaitGroup{},
	}
}

func (wp *WorkerPool) RunPool() {
	for _, worker := range wp.workers {
		wp.wg.Add(1)
		go func(worker *Worker) {
			defer wp.wg.Done()
			worker.Run()
		}(worker)
	}
}

func (wp *WorkerPool) StopPool() {
	close(wp.taskChan)
}
func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}
func (wp *WorkerPool) AddTask(task func()) {
	wp.taskChan <- task
	log.Print("AddTask: task sended")
}
