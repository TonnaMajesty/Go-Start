package main

import (
	"bytes"
	"context"
	"fmt"
	"sync"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type Manager struct {
	workerList map[AlgType][]*Worker
	lock       sync.RWMutex

	WorkerNumMap       map[AlgType]int
	ImageChanMap       map[AlgType]chan uint64
	ImageChanLengthMap map[AlgType]int
	ImagePriorChanMap  map[AlgType]chan uint64
	WorkerPoolMap      map[AlgType]chan *Worker
}

func NewManger(workerNumMap map[AlgType]int) *Manager {
	manager := &Manager{
		workerList:         make(map[AlgType][]*Worker),
		WorkerNumMap:       workerNumMap,
		ImageChanMap:       make(map[AlgType]chan uint64),
		ImageChanLengthMap: make(map[AlgType]int),
		ImagePriorChanMap:  make(map[AlgType]chan uint64),
		WorkerPoolMap:      make(map[AlgType]chan *Worker),
	}

	for algType, i := range workerNumMap {
		manager.ImageChanMap[algType] = make(chan uint64, 50*i)
		manager.ImageChanLengthMap[algType] = 50 * i
		manager.ImagePriorChanMap[algType] = make(chan uint64, 50*i)
		manager.WorkerPoolMap[algType] = make(chan *Worker, i)
	}

	return manager
}

func (m *Manager) Run(ctx context.Context) {
	for algType, num := range m.WorkerNumMap {
		for i := 0; i < num; i++ {
			m.AddWorker(ctx, algType, i)
		}
		logrus.Infof("Inspection Image Manger Create Analyze Worker, %d", num)
	}

	for algType, _ := range m.WorkerNumMap {
		go m.dispatch(algType)
	}
}

func (m *Manager) AddWorker(ctx context.Context, algType AlgType, id int) {
	work := NewWorker(m.WorkerPoolMap[algType], id, algType)
	m.workerList[algType] = append(m.workerList[algType], &work)
	work.Start(ctx)
	fmt.Printf("alg %d worker %d create, len %d\n", algType, id, len(m.workerList[algType]))
}

func (m *Manager) dispatch(algType AlgType) {
	t := time.NewTicker(time.Millisecond * 10)
	defer t.Stop()
	for {
		select {
		case image := <-m.ImagePriorChanMap[algType]:
			m.dispatchImage(algType, image)
		case <-t.C:
		priority:
			for {
				select {
				// 先把优先级高的队列清理完
				case image := <-m.ImagePriorChanMap[algType]:
					m.dispatchImage(algType, image)
				default:
					break priority
				}
			}

			select {
			case image := <-m.ImageChanMap[algType]:
				m.dispatchImage(algType, image)
			default:
			}
		}
	}
}

func (m *Manager) dispatchImage(algType AlgType, image uint64) {
	for {
		// TODO 有问题，会丢image
		//if len(m.WorkerPoolMap[algType]) == 0 {
		//	return
		//}
		worker := <-m.WorkerPoolMap[algType]
		if worker.stopped {
			close(worker.WImageChan)
			continue
		}
		worker.WImageChan <- image
		break
	}
}

func (m *Manager) AddJob(imageID uint64, algType AlgType) {
	if ch, ok := m.ImageChanMap[algType]; ok {
		ch <- imageID
	}
}

func (m *Manager) AddPriorJob(imageID uint64, algType AlgType) {
	if ch, ok := m.ImagePriorChanMap[algType]; ok {
		ch <- imageID
	}
}

func (m *Manager) IsFull(algType AlgType) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if _, ok := m.WorkerNumMap[algType]; !ok {
		return false
	}

	total := m.ImageChanLengthMap[algType]
	used := len(m.ImageChanMap[algType]) + len(m.ImagePriorChanMap[algType])

	return total-used < 15 // todo 数量大于每次查询待分析图片的数量，可配置
}

func (m *Manager) StopAllWorker() {
	for _, l := range m.workerList {
		for _, worker := range l {
			worker.Stop()
		}
	}
}

func (m *Manager) StopOneWorker(algType AlgType) {
	workers, ok := m.workerList[algType]
	if !ok || len(workers) == 0 {
		return
	}

	worker := workers[len(workers)-1]
	worker.Stop()
	m.workerList[algType] = slice.DropRight(workers, 1)
	fmt.Printf("Alg %d worker %d stopped, len %d \n", algType, worker.id, len(m.workerList[algType]))
}

func (m *Manager) AdjustWorkerNum(ctx context.Context, algType AlgType, num int) {
	m.lock.Lock()
	defer m.lock.Unlock()
	for alg, workers := range m.workerList {
		if alg == algType {
			curWorkerNum := len(workers)
			adjustNum := curWorkerNum - num
			if adjustNum > 0 {
				for i := 0; i < adjustNum; i++ {
					m.StopOneWorker(algType)
				}
			} else {
				for i := 0; i < -adjustNum; i++ {
					m.AddWorker(ctx, algType, curWorkerNum+i)
				}
			}
		}
	}

	m.WorkerNumMap[algType] = num
	m.ImageChanLengthMap[algType] = 50 * num
}

type Worker struct {
	id         int
	quit       chan bool
	stopped    bool
	algType    AlgType
	WImageChan chan uint64
	WorkerPool chan *Worker
}

func NewWorker(workerPool chan *Worker, id int, algType AlgType) Worker {
	return Worker{
		id:         id,
		algType:    algType,
		WorkerPool: workerPool,
		WImageChan: make(chan uint64, 500),
		quit:       make(chan bool),
	}
}

func (w Worker) Start(ctx context.Context) {
	go func() {
		for {
			w.WorkerPool <- &w
			select {
			case imageAnalyzeID := <-w.WImageChan:

				var entry *base.SentinelEntry
				var blockErr *base.BlockError

				for {
					entry, blockErr = sentinel.Entry(w.algType.String())
					if blockErr != nil {
						time.Sleep(100 * time.Millisecond)
						continue
					}
					break
				}

				err := DoAnalyze(ctx, imageAnalyzeID)
				if err != nil {
					// TODO 判断 err 类型
					sentinel.TraceError(entry, err)
					logrus.Errorf("AlgType %s ImageAnalyzeID %d DoAnalyzeAndUpdateImage failed, err:%s", w.algType, imageAnalyzeID, err)
					entry.Exit()
				} else {
					logrus.Infof("AlgType %s ImageAnalyzeID %d DoAnalyzeAndUpdateImage success", w.algType, imageAnalyzeID)
					entry.Exit()
				}
			case <-w.quit:
				logrus.Info("worker ", w.algType.String(), w.id, " stopping")
				return
			}
		}
	}()
}

func (w Worker) Stop() {
	go func() {
		w.quit <- true
		w.stopped = true
	}()
}

var client = resty.New().SetBaseURL("http://127.0.0.1:8080/transmissionChannel/8501").SetTimeout(1 * time.Second)

func DoAnalyze(ctx context.Context, imageId uint64) error {
	res := map[string]interface{}{}

	_, err := client.R().SetBody(map[string]interface{}{"test": "test"}).SetResult(&res).Post("/transmissionChannel/groundObjDetect")

	return err
}

// +gengo:enum
type AlgType uint8

// AI分析优先级
const (
	ALG_TYPE_UNKNOWN        AlgType = iota
	ALG_TYPE__AERIALDETECT          // 高空异物
	ALG_TYPE__WILDFIREWORKS         // 野外烟火
	ALG_TYPE__FOREIGNMATTER         // 线下异物
)

func (AlgType) EnumValues() []any {
	return []any{
		ALG_TYPE__AERIALDETECT, ALG_TYPE__WILDFIREWORKS, ALG_TYPE__FOREIGNMATTER,
	}
}
func (v AlgType) MarshalText() ([]byte, error) {
	str := v.String()
	return []byte(str), nil
}

func (v *AlgType) UnmarshalText(data []byte) error {
	vv, err := ParseAlgTypeFromString(string(bytes.ToUpper(data)))
	if err != nil {
		return err
	}
	*v = vv
	return nil
}

func ParseAlgTypeFromString(s string) (AlgType, error) {
	switch s {
	case "AERIALDETECT":
		return ALG_TYPE__AERIALDETECT, nil
	case "WILDFIREWORKS":
		return ALG_TYPE__WILDFIREWORKS, nil
	case "FOREIGNMATTER":
		return ALG_TYPE__FOREIGNMATTER, nil

	}
	return ALG_TYPE_UNKNOWN, nil
}

func (v AlgType) String() string {
	switch v {
	case ALG_TYPE__AERIALDETECT:
		return "AERIALDETECT"
	case ALG_TYPE__WILDFIREWORKS:
		return "WILDFIREWORKS"
	case ALG_TYPE__FOREIGNMATTER:
		return "FOREIGNMATTER"

	}
	return "UNKNOWN"
}

func ParseAlgTypeLabelString(label string) (AlgType, error) {
	switch label {
	case "高空异物":
		return ALG_TYPE__AERIALDETECT, nil
	case "野外烟火":
		return ALG_TYPE__WILDFIREWORKS, nil
	case "线下异物":
		return ALG_TYPE__FOREIGNMATTER, nil

	}
	return ALG_TYPE_UNKNOWN, nil
}

func (v AlgType) Label() string {
	switch v {
	case ALG_TYPE__AERIALDETECT:
		return "高空异物"
	case ALG_TYPE__WILDFIREWORKS:
		return "野外烟火"
	case ALG_TYPE__FOREIGNMATTER:
		return "线下异物"

	}
	return "UNKNOWN"
}
