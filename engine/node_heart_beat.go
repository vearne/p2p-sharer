package engine

import (
	"fmt"
	"github.com/imroc/req"
	"github.com/vearne/p2p-sharer/config"
	"github.com/vearne/p2p-sharer/models"
	"github.com/vearne/p2p-sharer/resource"
	"log"
	"time"
)

type HeartBeatWorker struct {
	RunningFlag bool // 是否运行 true:运行 false:停止
	ExitedFlag  bool //  已经退出的标识
	ExitChan    chan struct{}

	tracker string
}

func NewHeartBeatWorker(tracker string) *HeartBeatWorker {
	worker := &HeartBeatWorker{RunningFlag: true, ExitedFlag: false}
	worker.ExitChan = make(chan struct{})
	worker.tracker = tracker
	return worker
}

func (worker *HeartBeatWorker) Start() {
	log.Println("[start]HeartBeatWorker")
	// report nodeId and addr
	worker.SendHeartBeat()

	for worker.RunningFlag {
		select {
		case <-time.After(config.GetOpts().SessionExpire / 2):
			worker.SendHeartBeat()
		case <-worker.ExitChan:
			log.Println("HeartBeatWorker execute exit logic")
		}

	}
	worker.ExitedFlag = true
}

func (worker *HeartBeatWorker) Stop() {
	worker.RunningFlag = false
	close(worker.ExitChan)
	for !worker.ExitedFlag {
		time.Sleep(50 * time.Millisecond)
	}
	log.Println("[end]HeartBeatWorker")
}

func (worker *HeartBeatWorker) SendHeartBeat() {
	url := fmt.Sprintf("http://%v/v1/heartBeat", worker.tracker)
	var param heartBeatParam
	param = heartBeatParam{NodeId: resource.NodeId, Addr: resource.Addr}
	r, err := req.Post(url, req.BodyJSON(&param))
	if err != nil {
		log.Println("request tracker error", err)
		return
	}

	var resp models.ErrResponse
	r.ToJSON(&resp)
	log.Println("SendHeartBeat to tracker", resp.Code)
}
