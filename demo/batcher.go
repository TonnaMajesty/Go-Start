package main

import (
	"fmt"
	"time"
)

// 在C/S模型下，需要服务端主动推送大量数据时，通过定时或按数量到达一定条件时，将数据批量取出，实现合并推送，减少推送次数
// 支持按数量、长度、时间等方式对数据进行批量获取
func main() {
	b := NewBatcher(10, WithMaxItems(100))
	for i := 0; i < 200; i++ {
		b.Put(i)
	}

	// 获取数量为100的数组
	b.Get()
}

// Batcher 批处理接口
type Batcher interface {
	// 数据入列
	Put(interface{}) error
	// 数据出列
	Get() ([]interface{}, error)
	// 更新缓冲区
	Flush() error
	// 弃用
	Dispose()
	// 是否已弃用
	IsDisposed() bool
}

// mutex 锁
// 1.18之前，利用channel实现TryLock方法，不阻塞的加锁操作
// 1.18后，官方自带TryLock方法
type mutex struct {
	// This is really more of a semaphore design, but eh
	// Full -> locked, empty -> unlocked
	lock chan struct{}
}

func newMutex() *mutex {
	return &mutex{lock: make(chan struct{}, 1)}
}

func (m *mutex) Lock() {
	m.lock <- struct{}{}
}

func (m *mutex) Unlock() {
	<-m.lock
}

func (m *mutex) TryLock() bool {
	select {
	case m.lock <- struct{}{}:
		return true
	default:
		return false
	}
}

type BatcherNew struct {
	items     []interface{}      // 元素数组
	batchChan chan []interface{} // 缓冲区
	lock      *mutex             // 锁
	disposed  bool               // 弃用
	arrayLen  uint               // 数组长度
	option    Option             // 选项
}

type CalculateBytes func(interface{}) uint

type ErrDisposed struct {

}


func (e ErrDisposed) Error() string {
	return fmt.Sprintf("已弃用")
}

type Option struct {
	maxTime        time.Duration  // 按时间
	maxItems       uint           // 按数量
	maxBytes       uint           // 按长度
	availableBytes uint           // 可用长度
	calculateBytes CalculateBytes // 计算长度的函数
}

func (b *BatcherNew) Put(item interface{}) error {
	b.lock.Lock()
	if b.disposed { // 判断队列是否可用
		b.lock.Unlock()
		return ErrDisposed{}
	}

	// 添加元素
	b.items = append(b.items, item)
	if b.option.calculateBytes != nil { // 计算长度
		b.option.availableBytes += b.option.calculateBytes(item)
	}
	if b.ready() { // 如果满足条件，将数据刷入缓冲区
		b.flush()
	}
	b.lock.Unlock()
	return nil
}
func (b *BatcherNew) ready() bool {
	// 按数量判断
	if b.option.maxItems != 0 && uint(len(b.items)) >= b.option.maxItems {
		return true
	}
	// 按字节数判断
	if b.option.maxBytes != 0 && b.option.availableBytes >= b.option.maxBytes {
		return true
	}
	return false
}

func (b *BatcherNew) Get() ([]interface{}, error) {
	// 定时器
	var timeout <-chan time.Time
	if b.option.maxTime > 0 {
		timeout = time.After(b.option.maxTime)
	}

	select {
	case items, ok := <-b.batchChan:
		if !ok {
			return nil, ErrDisposed{}
		}
		return items, nil
	case <-timeout: // 定时获取是阻塞式的
		for {
			if b.lock.TryLock() { // 尝试加锁
				select {
				case items, ok := <-b.batchChan:
					b.lock.Unlock()
					if !ok {
						return nil, ErrDisposed{}
					}
					return items, nil
				default:
				}
				// 直接取当前数据项
				items := b.items
				b.items = make([]interface{}, 0, b.arrayLen)
				b.option.availableBytes = 0
				b.lock.Unlock()
				return items, nil
			} else { // 加锁失败,说明可能正在Put或Flush
				select {
				case items, ok := <-b.batchChan:
					if !ok {
						return nil, ErrDisposed{}
					}
					// 从缓冲区读取到数据，直接返回
					return items, nil
				default:
					// 继续循环，尝试取数据
				}
			}
		}

	}
}

func (b *BatcherNew) Flush() error {
	// This is the same pattern as a Put
	b.lock.Lock()
	if b.disposed {
		b.lock.Unlock()
		return ErrDisposed{}
	}
	b.flush()
	b.lock.Unlock()
	return nil
}

// flush 将数组输出到缓冲区
func (b *BatcherNew) flush() {
	b.batchChan <- b.items
	// 重新初始化
	b.items = make([]interface{}, 0, b.arrayLen)
	b.option.availableBytes = 0
}

func (b *BatcherNew) Dispose() {
	for {
		if b.lock.TryLock() {
			if b.disposed {
				b.lock.Unlock()
				return
			}
			b.disposed = true
			b.items = nil
			b.drainBatchChan()
			close(b.batchChan)
			b.lock.Unlock()
		} else {
			b.drainBatchChan()
		}
	}
}

func (b *BatcherNew) IsDisposed() bool {
	b.lock.Lock()
	disposed := b.disposed
	b.lock.Unlock()
	return disposed
}

func (b *BatcherNew) drainBatchChan() {
	for {
		select {
		case <-b.batchChan:
		default:
			return
		}
	}
}

// WithMaxTime 按时间取数据
func WithMaxTime(maxTime time.Duration) Option {
	return Option{maxTime: maxTime}
}

// WithMaxItems 按数量取数据
func WithMaxItems(maxItems uint) Option {
	return Option{maxItems: maxItems}
}

// WithMaxBytes 按字节长度取数据
func WithMaxBytes(maxBytes uint, calculateBytes CalculateBytes) Option {
	return Option{maxBytes: maxBytes, calculateBytes: calculateBytes}
}

// NewBatcher 初始化
func NewBatcher(queueLen uint, option Option) Batcher {
	var arrayLen uint = 1024
	if option.maxItems > 0 {
		arrayLen = option.maxItems
	}
	return &BatcherNew{
		option:    option,
		items:     make([]interface{}, 0, arrayLen),
		batchChan: make(chan []interface{}, queueLen),
		lock:      newMutex(),
		disposed:  false,
		arrayLen:  arrayLen,
	}
}
