package main

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)
// https://www.jianshu.com/p/3bb67e53e11f

func main() {
	q := NewDelayQueue(&LocalMsgQueueConf{10, 10*time.Second, 2})
	q.initWorker()

}

type LocalMsg struct {
	message             string
	ExpectExecuteTimeMs int64
}

type LocalMsgQueueConf struct {
	BufferSize  int           // chan 缓冲区大小
	ChanTimeOut time.Duration // 加入chan的超时时间
	WorkerCount int64         // 消费worker 数
}

type LocalMsgQueue struct {
	messages           chan LocalMsg
	conf               *LocalMsgQueueConf
	currentWorkerCount int64 // 当前运行的worker数量
}

func NewDelayQueue(config *LocalMsgQueueConf) *LocalMsgQueue {
	delayQueue := &LocalMsgQueue{
		messages:           make(chan LocalMsg, config.BufferSize),
		conf:               config,
		currentWorkerCount: 0,
	}
	delayQueue.initWorker()
	return delayQueue
}

func (d *LocalMsgQueue) initWorker() {
	go func() {
		defer func() {
			if errPanic := recover(); errPanic != nil {
				// 打印堆栈，异常上报
			}
		}()

		for {
			if d.currentWorkerCount < d.conf.WorkerCount {
				atomic.AddInt64(&d.currentWorkerCount, 1)
				go d.process()
			}
		}
	}()
}

func (d *LocalMsgQueue) process() (err error) {
	msg := LocalMsg{}
	defer func() {
		if errPanic := recover(); errPanic != nil {
			// 打印堆栈，异常上报
		}
		if err != nil {
			// 打印日志，异常上报
		}
		atomic.AddInt64(&d.currentWorkerCount, -1)
	}()
	for {
		select {
		case msg = <-d.messages:
			cTime := time.Now().UnixNano() / 1e6
			if cTime < msg.ExpectExecuteTimeMs {
				time.Sleep(time.Duration(msg.ExpectExecuteTimeMs-cTime) * time.Millisecond)
			}
		}
	}
}

func (d *LocalMsgQueue) Send(req LocalMsg) (err error) {
	chanCtx, cancel := context.WithTimeout(context.Background(), d.conf.ChanTimeOut)
	defer func() {
		if err != nil {
			// 打印日志，异常上报
		}
		cancel()
	}()
	select {
	case <-chanCtx.Done():
		if chanCtx.Err() == context.DeadlineExceeded {
			return errors.New("deadline_exceeded")
		}
		return chanCtx.Err()
	case d.messages <- req:
		return nil
	}
}


