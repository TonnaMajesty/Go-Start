package main

import (
	"context"

	"github.com/sirupsen/logrus"
)

type Manager struct {
	WorkerNum  int
	workerList []Worker
	MImageChan chan *Pic
	WorkerPool chan chan *Pic
}

func NewManger(workerNum int) *Manager {
	workPool := make(chan chan *Pic, workerNum)
	return &Manager{
		WorkerNum:  workerNum,
		MImageChan: make(chan *Pic, 500),
		WorkerPool: workPool,
	}
}

func (m *Manager) Run(ctx context.Context) {
	for i := 0; i < m.WorkerNum; i++ {
		work := NewWorker(m.WorkerPool, i)
		m.workerList = append(m.workerList, work)
		work.Start(ctx)
	}

	go m.dispatch()
}

func (m *Manager) dispatch() {
	for {
		select {
		case image := <-m.MImageChan:
			go func(image *Pic) {
				imageChan := <-m.WorkerPool
				imageChan <- image
			}(image)
		}
	}
}

func (m *Manager) AddJob(image *Pic) {
	m.MImageChan <- image
}

func (m *Manager) Stop() {
	for _, v := range m.workerList {
		v.Stop()
	}
}

type Worker struct {
	id         int
	quit       chan bool
	WImageChan chan *Pic
	WorkerPool chan chan *Pic
}

func NewWorker(workerPool chan chan *Pic, id int) Worker {
	return Worker{
		id:         id,
		WorkerPool: workerPool,
		WImageChan: make(chan *Pic, 500),
		quit:       make(chan bool),
	}
}

func (w Worker) Start(ctx context.Context) {
	go func() {
		for {
			w.WorkerPool <- w.WImageChan
			select {
			case image := <-w.WImageChan:
				DealPicture(ctx, image)
			case <-w.quit:
				logrus.Info("worker ", w.id, " stopping")
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}
