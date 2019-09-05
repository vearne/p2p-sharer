package utils

import (
	"context"
	"fmt"
	"sync"
)

// 一个简易的协程池实现
type JobContextFunc func(ctx context.Context, key interface{}) *GPResult

type GContextPool struct {
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
	// Context
	Ctx context.Context
}

func NewGContextPool(ctx context.Context, size int) *GContextPool {
	pool := GContextPool{}
	pool.JobChan = make(chan interface{}, SIZE)
	pool.ResultChan = make(chan *GPResult, SIZE)
	pool.Size = size
	pool.IsClose = false
	pool.Ctx = ctx
	return &pool
}

func (p *GContextPool) ApplyAsync(f JobContextFunc, slice []interface{}) <-chan *GPResult {

	p.TargetCount = len(slice)
	// Producer
	go p.Produce(slice)
	// consumer
	for i := 0; i < p.Size; i++ {
		go p.Consume(f)
	}

	return p.ResultChan
}

func (p *GContextPool) Produce(slice []interface{}) {
	for _, key := range slice {
		p.JobChan <- key
	}
	close(p.JobChan)
}

func doCtxOne(ctx context.Context, job interface{}, f JobContextFunc) (result *GPResult) {
	defer func() {
		r := recover()
		if r != nil {
			err := fmt.Errorf("execute job error, recover%v, job:%v", r, job)
			result = &GPResult{Value: nil, Err: err}
		}
	}()
	return f(ctx, job)
}

func (p *GContextPool) Consume(f JobContextFunc) {
	for job := range p.JobChan {

		select {
		case <-p.Ctx.Done():
			result := GPResult{Value: nil,
				Err: fmt.Errorf("execute was canceled, job:%v", job)}
			p.ResultChan <- &result
		default:
			// 没有结束 ... 执行 ...
			p.ResultChan <- doCtxOne(p.Ctx, job, f)
		}
		p.FinishOne()
	}
	p.TryClose()
}

// 记录完成了一个任务
func (p *GContextPool) FinishOne() {
	p.Lock()
	p.FinishCount++
	p.Unlock()
}

// 关闭结果Channel
func (p *GContextPool) TryClose() {
	p.Lock()
	if p.FinishCount == p.TargetCount && !p.IsClose {
		close(p.ResultChan)
		p.IsClose = true
	}
	p.Unlock()
}
