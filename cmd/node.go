package cmd

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/vearne/p2p-sharer/config"
	"github.com/vearne/p2p-sharer/consts"
	"github.com/vearne/p2p-sharer/engine"
	"github.com/vearne/p2p-sharer/models"
	"github.com/vearne/p2p-sharer/resource"
	"github.com/vearne/p2p-sharer/utils"
	manager "github.com/vearne/worker_manager"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "node service",
	Long:  "node service",
	Run:   RunNode,
}

func init() {
	rootCmd.AddCommand(nodeCmd)
	utils.SetConnPool()
}

func RunNode(cmd *cobra.Command, args []string) {
	config.InitConfig("node")

	initNodeResource()

	// scan download dir to find files which has been downloaded
	scanAndReport(config.GetOpts().DownloadDir)

	log.Println("start worker")
	wm := manager.NewWorkerManager()
	wm.AddWorker(engine.NewNodeWebServer())
	wm.AddWorker(engine.NewDownloadWorker())

	wm.Start()

	// register grace exit
	gracefulExit(wm)

	// block and wait
	wm.Wait()
}

func initNodeResource() {
	resource.InitBaseNodeInfo()
	resource.TaskChan = make(chan *models.SeedInfo, 10)
	resource.PieceLocalMapper = models.NewPieceLocalMapper()
}

func scanAndReport(downloadDir string) {
	fileList, err := ioutil.ReadDir(downloadDir)
	if err != nil {
		log.Fatal("can't open downloadDir")
	}
	for _, item := range fileList {
		if !item.IsDir() {
			name := item.Name()
			if strings.HasSuffix(name, consts.DownloadFinishPostfix) {
				// 1. check if data file exist
				dataFileName := name[0 : len(name)-len(consts.DownloadFinishPostfix)]
				if !Exists(filepath.Join(downloadDir, dataFileName)) {
					continue
				}
				// 2. check if seed file exist
				seedFileName := dataFileName + consts.SeedPostfix
				if !Exists(filepath.Join(downloadDir, seedFileName)) {
					continue
				}

				// 3. represent the file is downloaded
				seedFilePath := filepath.Join(downloadDir, seedFileName)
				dataFilePath := filepath.Join(downloadDir, dataFileName)

				data, err := ioutil.ReadFile(seedFilePath)
				if err != nil {
					log.Println("read file error", seedFilePath)
				}
				var seedInfo models.SeedInfo
				json.Unmarshal(data, &seedInfo)
				// 3.1 fill resource.PieceLocalMapper
				for _, piece := range seedInfo.Pieces {
					temp := models.LocalInfo{*piece, dataFilePath}
					resource.PieceLocalMapper.Add(models.Piece(piece.Checksum), &temp)
				}
				// 3.2 report to tracker
				makeOurselfAsSource(resource.NodeId, &seedInfo)
			}
		}
	}
}

func makeOurselfAsSource(nodeId string, info *models.SeedInfo) {
	// 1.
	hbworker := engine.NewHeartBeatWorker(info.TrackerAddr)
	// 0. a new goroutine to send heartbeat
	go hbworker.Start()
	defer func() {
		// Need to keep the heartbeat for a while so that other nodes can download
		// worker's gracefull-exit is not considered
		time.AfterFunc(3*time.Hour, func() {
			hbworker.Stop()
		})

	}()

	// 2. report to tracker
	//ReportToTracker(tracker string, pieceId string, task *models.DownloadTask)
	task := models.DownloadTask{}
	task.Total = len(info.Pieces)
	task.SuccessCount = task.Total
	task.FileName = info.FileName
	successCount := 0
	failedCount := 0
	for _, piece := range info.Pieces {
		ok := engine.ReportToTracker(info.TrackerAddr, piece.Checksum, &task)
		if ok {
			successCount++
		} else {
			failedCount++
		}
	}
	log.Println("_report", "fileName", info.FileName, "successCount", successCount,
		"failedCount", failedCount)
}

// Exists reports whether the named file or directory exists.
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
