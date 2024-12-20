package workerpool

import (
	"log"
	"sync"
)

type WorkerPool struct {
	workers  []*Worker
	taskChan chan<- func()
	wgGen    *sync.WaitGroup
	wgCloser *sync.WaitGroup
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
		wgGen:    &sync.WaitGroup{},
	}
}

func (wp *WorkerPool) RunPool() {
	for _, worker := range wp.workers {
		wp.wgGen.Add(1)
		go func(worker *Worker) {
			defer wp.wgGen.Done()
			worker.Run()
		}(worker)
	}
}

func (wp *WorkerPool) StopPool() {
	close(wp.taskChan)
}
func (wp *WorkerPool) Wait() {
	wp.wgGen.Wait()
}
func (wp *WorkerPool) AddTask(task func()) {
	wp.wgCloser.Add(1)
	go func() {
		wp.taskChan <- task
		log.Print("AddTask: task sended")
		wp.wgCloser.Done()
	}()
	wp.wgCloser.Wait()
	close(wp.taskChan)
}
