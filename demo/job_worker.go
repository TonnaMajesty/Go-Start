package main

import (
	"fmt"
	"sync"
	"time"
)

// https://mp.weixin.qq.com/s/aR5s936ihICiZYOQyMoi5A

type SubWorkerNew struct {
	Id      int
	JobChan chan string
}

type W2New struct {
	SubWorkers []SubWorkerNew
	MaxNum     int
	ChPool     chan chan string
	QuitChan   chan struct{}
	Wg         *sync.WaitGroup
}

func NewW2(maxNum int) *W2New {
	chPool := make(chan chan string, maxNum)
	subWorkers := make([]SubWorkerNew, maxNum)
	for i := 0; i < maxNum; i++ {
		subWorkers[i] = SubWorkerNew{Id: i, JobChan: make(chan string)}
		chPool <- subWorkers[i].JobChan
	}
	wg := new(sync.WaitGroup)
	wg.Add(maxNum)

	return &W2New{
		MaxNum:     maxNum,
		SubWorkers: subWorkers,
		ChPool:     chPool,
		QuitChan:   make(chan struct{}),
		Wg:         wg,
	}
}

func (w *W2New) StartPool() {
	for i := 0; i < w.MaxNum; i++ {
		go func(wg *sync.WaitGroup, subWorker *SubWorkerNew) {
			defer wg.Done()
			for {
				select {
				case job := <-subWorker.JobChan:
					fmt.Printf("SubWorker %d processing job %s\n", subWorker.Id, job)
					time.Sleep(time.Second) // 模拟任务处理过程
				case <-w.QuitChan:
					return
				}
			}
		}(w.Wg, &w.SubWorkers[i])
	}
}

func (w *W2New) Stop() {
	close(w.QuitChan)
	for i := 0; i < w.MaxNum; i++ {
		close(w.SubWorkers[i].JobChan)
	}
	w.Wg.Wait()
}

func (w *W2New) Dispatch(job string) {
	select {
	case jobChan := <-w.ChPool:
		jobChan <- job
	default:
		// todo 阻塞到<-w.ChPool 或 返回error调用方重试
		fmt.Println("All workers busy")
	}
}
func (w *W2New) AddWorker() {
	newWorker := SubWorkerNew{Id: w.MaxNum, JobChan: make(chan string)}
	w.SubWorkers = append(w.SubWorkers, newWorker)
	w.ChPool <- newWorker.JobChan
	w.MaxNum++
	w.Wg.Add(1)
	go func(subWorker *SubWorkerNew) {
		defer w.Wg.Done()

		for {
			select {
			case job := <-subWorker.JobChan:
				fmt.Printf("SubWorker %d processing job %s\n", subWorker.Id, job)
				time.Sleep(time.Second) // 模拟任务处理过程
			case <-w.QuitChan:
				return
			}
		}
	}(&newWorker)
}

func (w *W2New) RemoveWorker() {
	if w.MaxNum > 1 {
		worker := w.SubWorkers[w.MaxNum-1]
		close(worker.JobChan)
		w.MaxNum--
		w.SubWorkers = w.SubWorkers[:w.MaxNum]
	}
}

func main() {
	pool := NewW2(3)
	pool.StartPool()

	pool.Dispatch("task 1")
	pool.Dispatch("task 2")
	pool.Dispatch("task 3")

	pool.AddWorker()
	pool.AddWorker()
	pool.RemoveWorker()

	pool.Stop()
}
