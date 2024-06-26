package main

import (
	"sync"

	"github.com/hidez8891/shm"
)

type ShmManager struct {
	ShmName    string // 共享内存名称
	ShmBinsNum int64  // 共享内存分区数
	ShmBinSize int64  // 共享内存分区大小

	shm            *shm.Memory
	shmBinIndexMux sync.Mutex
	shmBinIndex    int64
}

func NewShmManager(shmName string, shmBinsNum int64, shmBinSize int64) *ShmManager {
	return &ShmManager{
		ShmName:    shmName,
		ShmBinsNum: shmBinsNum,
		ShmBinSize: shmBinSize,
	}
}

func (s *ShmManager) Open() error {
	r, err := shm.Open(s.ShmName, int32(s.ShmBinsNum*s.ShmBinSize))
	if err != nil {
		return err
	}

	s.shm = r
	return nil
}

func (s *ShmManager) Create() error {
	r, err := shm.Create(s.ShmName, int32(s.ShmBinsNum*s.ShmBinSize))
	if err != nil {
		return err
	}

	s.shm = r
	return nil
}

func (s *ShmManager) Close() error {
	if s.shm == nil {
		return nil
	}
	return s.shm.Close()
}

func (s *ShmManager) WriteIntoShm(data []byte) (int64, error) {
	index := s.GetShmIndex()
	_, err := s.shm.WriteAt(data, s.ShmBinSize*index)

	return index, err
}

func (s *ShmManager) ReadAt(index int64, size int) ([]byte, error) {
	buf := make([]byte, size)
	_, err := s.shm.ReadAt(buf, s.ShmBinSize*index)

	return buf, err
}

func (s *ShmManager) GetShmIndex() int64 {
	s.shmBinIndexMux.Lock()
	defer s.shmBinIndexMux.Unlock()
	if s.shmBinIndex == s.ShmBinsNum-1 {
		s.shmBinIndex = 0
	} else {
		s.shmBinIndex = s.shmBinIndex + 1
	}

	return s.shmBinIndex
}
