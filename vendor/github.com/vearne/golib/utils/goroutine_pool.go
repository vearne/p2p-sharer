package utils

import (
	"fmt"
	"sync"
)

const (
	SIZE int = 50
)

type GPResult struct {
	Value interface{}
	Err   error
}

// 一个简易的协程池实现
type JobFunc func(param interface{}) *GPResult

type GPool struct {
	sync.Mutex

	// 任务队列
	JobChan chan interface{}
	// 结果队列
	ResultChan chan *GPResult
	// 协程池的大小
	Size int
	// 已经完成的任务量
	FinishCount int
	// 目标任务量
	TargetCount int
	// ResultChan 是否Close
	IsClose bool
}

func NewGPool(size int) *GPool {
	pool := GPool{}
	pool.JobChan = make(chan interface{}, SIZE)
	pool.ResultChan = make(chan *GPResult, SIZE)
	pool.Size = size
	pool.IsClose = false
	return &pool
}

func (p *GPool) ApplyAsync(f JobFunc, slice []interface{}) <-chan *GPResult {

	p.TargetCount = len(slice)
	// Producer
	go p.Produce(slice)
	// consumer
	for i := 0; i < p.Size; i++ {
		go p.Consume(f)
	}

	return p.ResultChan
}

func (p *GPool) Produce(slice []interface{}) {
	for _, key := range slice {
		p.JobChan <- key
	}
	close(p.JobChan)
}

func doOne(job interface{}, f JobFunc) (result *GPResult) {
	defer func() {
		r := recover()
		if r != nil {
			err := fmt.Errorf("execute job error, recover%v, job:%v", r, job)
			result = &GPResult{Value: nil, Err: err}
		}
	}()
	return f(job)
}

func (p *GPool) Consume(f JobFunc) {
	for job := range p.JobChan {
		p.ResultChan <- doOne(job, f)
		p.FinishOne()
	}
	p.TryClose()
}

// 记录完成了一个任务
func (p *GPool) FinishOne() {
	p.Lock()
	p.FinishCount++
	p.Unlock()
}

// 关闭结果Channel
func (p *GPool) TryClose() {
	p.Lock()
	if p.FinishCount == p.TargetCount && !p.IsClose {
		close(p.ResultChan)
		p.IsClose = true
	}
	p.Unlock()
}
