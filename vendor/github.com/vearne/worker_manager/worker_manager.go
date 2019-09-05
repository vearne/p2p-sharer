package worker_manager

import (
	"fmt"
	"runtime"
	"sync"
)

type Worker interface {
	Start()
	Stop()
}

type WorkerManager struct {
	sync.WaitGroup
	// 保存所有worker
	WorkerSlice []Worker
}

func Stack() []byte {
	buf := make([]byte, 2048)
	n := runtime.Stack(buf, false)
	return buf[:n]
}

func NewWorkerManager() *WorkerManager {
	workerManager := WorkerManager{}
	workerManager.WorkerSlice = make([]Worker, 0, 10)
	return &workerManager
}

func (wm *WorkerManager) AddWorker(w Worker) {
	wm.WorkerSlice = append(wm.WorkerSlice, w)
}

func (wm *WorkerManager) Start() {
	wm.Add(len(wm.WorkerSlice))
	for _, worker := range wm.WorkerSlice {
		go func(w Worker) {
			defer func() {
				err := recover()
				if err != nil {
					fmt.Printf("WorkerManager error, error:%v, stack:%v\n",
						err, string(Stack()))
				}
			}()
			w.Start()
		}(worker)
	}
}

func (wm *WorkerManager) Stop() {
	for _, worker := range wm.WorkerSlice {
		go func(w Worker) {
			defer func() {
				err := recover()
				if err != nil {
					fmt.Printf("WorkerManager error, error:%v, stack:%v\n",
						err, string(Stack()))
				}
			}()

			w.Stop()
			wm.Done()
		}(worker)
	}
}
