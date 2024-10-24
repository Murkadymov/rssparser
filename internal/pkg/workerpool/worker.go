package workerpool

import "log"

type Worker struct {
	WorkerID int
	taskChan <-chan func()
}

func NewWorker(id int, taskChan chan func()) *Worker {
	return &Worker{
		WorkerID: id,
		taskChan: taskChan,
	}
}

func (w *Worker) Run() {
	log.Printf("worker id: %d started reading channel", w.WorkerID)
	for task := range w.taskChan {

		log.Printf("worker id: %d got task", w.WorkerID)
		task()
		log.Printf("worker id: %d ended task", w.WorkerID)
	}
	log.Printf("worker id: %d finished work", w.WorkerID)
}
