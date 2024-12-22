package workerpool

import (
	"github.com/mmcdole/gofeed"
	"log"
	"rssparser/internal/service/feed"
	"sync"
)

type e struct {
	workers  []*feed.FeedWorker
	taskChan chan *gofeed.Item
	wgGen    *sync.WaitGroup
	wgCloser *sync.WaitGroup
}

func NewWorkerPool(workerCnt int) *WorkerPool {
	taskChan := make(chan *gofeed.Item)
	workers := make([]*feed.FeedWorker, workerCnt)
	for i := range workerCnt {
		workers[i] = feed.NewFeedWorker()
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
		go func(worker *feed.FeedWorker) {
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
func (wp *WorkerPool) AddTask(item *gofeed.Item) {
	wp.wgCloser.Add(1)
	go func() {
		wp.taskChan <- item
		log.Print("Add Item to read: item sended")
		wp.wgCloser.Done()
	}()
	wp.wgCloser.Wait()
	close(wp.taskChan)
}
