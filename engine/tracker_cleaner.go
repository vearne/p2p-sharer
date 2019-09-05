package engine

import (
	"github.com/vearne/p2p-sharer/config"
	"github.com/vearne/p2p-sharer/resource"
	"log"
	"time"
)


type CleanerWorker struct {
	RunningFlag bool // 是否运行 true:运行 false:停止
	ExitedFlag  bool //  已经退出的标识
	ExitChan    chan struct{}
}

func NewCleanerWorker() *CleanerWorker {
	worker := &CleanerWorker{RunningFlag: true, ExitedFlag: false}
	worker.ExitChan = make(chan struct{})
	return worker
}

func (worker *CleanerWorker) Start() {
	log.Println("[start]CleanerWorker")
	for worker.RunningFlag {
		select {
		case <-time.After(config.GetOpts().SessionCheckInterval):
			clean()
		case <-worker.ExitChan:
			log.Println("CleanerWorker execute exit logic")
		}

	}
	worker.ExitedFlag = true
}

func (worker *CleanerWorker) Stop() {
	worker.RunningFlag = false
	close(worker.ExitChan)
	for !worker.ExitedFlag {
		time.Sleep(50 * time.Millisecond)
	}
	log.Println("[end]CleanerWorker")
}

func clean() {
	nodeList := resource.NodeSession.FindAndPurgeInvalid(config.GetOpts().SessionExpire)
	log.Println("tracker-clean", nodeList)
	for _, node := range nodeList {
		resource.FilePieceMapper.PurgeNode(node)
	}
}
