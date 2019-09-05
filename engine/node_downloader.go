package engine

import (
	"fmt"
	"github.com/imroc/req"
	gu "github.com/vearne/golib/utils"
	"github.com/vearne/p2p-sharer/config"
	"github.com/vearne/p2p-sharer/consts"
	"github.com/vearne/p2p-sharer/models"
	"github.com/vearne/p2p-sharer/resource"
	"github.com/vearne/p2p-sharer/utils"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type DownloadWorker struct {
	RunningFlag bool // 是否运行 true:运行 false:停止
	ExitedFlag  bool //  已经退出的标识
	ExitChan    chan struct{}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewDownloadWorker() *DownloadWorker {
	worker := &DownloadWorker{RunningFlag: true, ExitedFlag: false}
	worker.ExitChan = make(chan struct{})
	return worker
}

func (worker *DownloadWorker) Start() {
	log.Println("[start]DownloadWorker")

	for worker.RunningFlag {
		select {
		case task := <-resource.TaskChan:
			worker.download(task)
		case <-worker.ExitChan:
			log.Println("DownloadWorker execute exit logic")
		}

	}
	worker.ExitedFlag = true
}

func (worker *DownloadWorker) Stop() {
	worker.RunningFlag = false
	close(worker.ExitChan)
	for !worker.ExitedFlag {
		time.Sleep(50 * time.Millisecond)
	}
	log.Println("[end]DownloadWorker")
}

func (worker *DownloadWorker) download(info *models.SeedInfo) {
	hbworker := NewHeartBeatWorker(info.TrackerAddr)
	// 0. a new goroutine to send heartbeat
	go hbworker.Start()
	defer func() {
		// Need to keep the heartbeat for a while so that other nodes can download
		// worker's gracefull-exit is not considered
		time.AfterFunc(1*time.Minute, func() {
			hbworker.Stop()
		})
	}()

	// 1. init empty file, but it has allocate disk space.
	log.Println("1. init empty file")
	localFilepath := filepath.Join(config.GetOpts().DownloadDir, info.FileName)
	fp, err := os.Create(localFilepath)
	defer fp.Close()
	if err != nil {
		log.Println("filepath create failed")
		return
	}

	temp := make([]byte, BuffSize*1024)
	for i := 0; i*BuffSize*1024 < info.Length; i++ {
		size := min(BuffSize*1024, info.Length-i*BuffSize*1024)
		fp.Write(temp[0:size])
	}

	// 2  init data struct to represent Task
	var task models.DownloadTask
	task.FileName = info.FileName
	task.Total = len(info.Pieces)
	task.SuccessCount = 0
	task.Pieces = info.Pieces
	task.WaitForDeal = gu.NewIntSet()
	for _, item := range info.Pieces {
		task.WaitForDeal.Add(item.Index)
	}
	log.Println("2. init DownloadTask", "piece count", task.Total, "SuccessCount", task.SuccessCount)
	// 3. execute
	log.Println("3. execute")
	var gotFlag bool = false
	for task.WaitForDeal.Size() > 0 && worker.RunningFlag {
		gotFlag = false

		n := task.WaitForDeal.Size()
		rIdx := rand.Intn(n)
		idx := task.WaitForDeal.ToSlice()[rIdx]
		piece := task.Pieces[idx]
		addrList := getNodes(info.TrackerAddr, piece.Checksum)

		for _, addr := range addrList {
			// 3.1 getPieceData
			data, ok := getPieceData(addr, piece.Checksum)
			log.Println("download-getPieceData", "from", addr, ok)
			if ok {

				log.Println("3.1 getPieceData", piece.Index, piece.Checksum)
				// 3.2 verify checksum

				got := utils.BlockHash(data)
				log.Println("3.2 verify checksum", piece.Index,
					"targetCheckSum", piece.Checksum, "realCheckSum", got)
				if got != piece.Checksum {
					continue
				}

				// 3.3 write file
				log.Println("3.3 write file", piece.Index, piece.Checksum)
				fp.Seek(int64(idx)*config.GetOpts().PieceSize, 0)
				fp.Write(data)

				// 3.4 update PieceLocalMapper
				log.Println("3.4 update PieceLocalMapper", piece.Index, piece.Checksum)
				localInfo := models.LocalInfo{PieceInfo: *piece, FilePath: localFilepath}
				resource.PieceLocalMapper.Add(models.Piece(piece.Checksum), &localInfo)

				// 3.5 report to tracker
				log.Println("3.5 report to tracker", piece.Index, piece.Checksum)
				ReportToTracker(info.TrackerAddr, piece.Checksum, &task)

				// 3.6 update progress
				log.Println("3.6 update progress", piece.Index, piece.Checksum)
				gotFlag = true
				task.SuccessCount++
				task.WaitForDeal.Remove(piece.Index)
				if task.SuccessCount == task.Total {
					tagDownloadFinish(&task)
				}

				break
			}
		}

		if !gotFlag {
			time.Sleep(200 * time.Millisecond)
		}
	}

	log.Println("ok. Dowload is finished")
}

func min(a, b int) int {
	if a <= b {
		return a
	} else {
		return b
	}
}

func tagDownloadFinish(task *models.DownloadTask) {
	okFilePath := filepath.Join(config.GetOpts().DownloadDir,
		task.FileName+consts.DownloadFinishPostfix)
	sf, err := os.Create(okFilePath)
	defer sf.Close()
	if err != nil {
		log.Println(err)
	}
}

func getNodes(tracker string, pieceId string) []string {
	url := fmt.Sprintf("http://%v/v1/nodeList", tracker)
	param := req.Param{}
	param["pieceID"] = pieceId
	r, err := req.Get(url, param)
	if err != nil {
		log.Println("getNodes", err)
		return make([]string, 0)
	}

	if r.Response().StatusCode == 200 {
		var resp models.NodeListResp
		r.ToJSON(&resp)

		// shuffle
		if len(resp.Nodes) > 0 {
			rand.Shuffle(len(resp.Nodes), ReSortNodes(resp.Nodes))
		}
		return resp.Nodes
	}
	return make([]string, 0)
}

func ReSortNodes(nodes []string) func(i, j int) {
	return func(i, j int) {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	}
}

func getPieceData(addr string, pieceId string) (data []byte, ok bool) {
	url := fmt.Sprintf("http://%v/v1/pieceData", addr)
	log.Println("getPieceData", url)

	param := req.Param{}
	param["pieceID"] = pieceId

	r, err := req.Get(url, param)
	if err != nil {
		log.Println("getPieceData", err)
		return nil, false
	}

	data, err = r.ToBytes()
	if err != nil {
		log.Println("getPieceData", err)
		return nil, false
	}
	log.Println("getPieceData", r.Response().StatusCode, len(data))
	return data, true
}

func ReportToTracker(tracker string, pieceId string, task *models.DownloadTask) bool {
	var param models.ReportParam
	param.PieceID = pieceId
	param.NodeID = resource.NodeId
	param.File = task.FileName
	param.Progress = float64(task.SuccessCount) / float64(task.Total)

	url := fmt.Sprintf("http://%v/v1/report", tracker)
	_, err := req.Post(url, req.BodyJSON(&param))
	if err != nil {
		log.Println("reportToTracker error", err)
		return false
	}
	return true
}
