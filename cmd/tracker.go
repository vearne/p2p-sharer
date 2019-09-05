package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vearne/p2p-sharer/config"
	"github.com/vearne/p2p-sharer/engine"
	"github.com/vearne/p2p-sharer/models"
	"github.com/vearne/p2p-sharer/resource"
	manager "github.com/vearne/worker_manager"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var trackerCmd = &cobra.Command{
	Use:   "tracker",
	Short: "tracker service",
	Long:  "tracker service",
	Run:   RunTracker,
}

func init() {
	rootCmd.AddCommand(trackerCmd)
}

func RunTracker(cmd *cobra.Command, args []string) {
	config.InitConfig("tracker")


	initTrackerResource()

	wm := manager.NewWorkerManager()
	wm.AddWorker(engine.NewTrackerWebServer())
	wm.AddWorker(engine.NewCleanerWorker())

	wm.Start()

	// register grace exit
	gracefulExit(wm)

	// block and wait
	wm.Wait()
}

func gracefulExit(wm *manager.WorkerManager) {
	ch := make(chan os.Signal, 1)
	// SIGPIPE信号
	// http://senlinzhan.github.io/2017/03/02/sigpipe/
	// SIGPIPE产生的原因是这样的：如果一个 socket 在接收到了 RST packet 之后，
	// 程序仍然向这个 socket 写入数据， 那么就会产生SIGPIPE信号。
	signal.Ignore(syscall.SIGPIPE)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	sig := <-ch

	log.Printf("got a signal, %v\n", sig)
	wm.Stop()
}


func initTrackerResource(){
	resource.NodeSession = models.NewSession()
	resource.FilePieceMapper = models.NewPieceMapper()
	resource.NodeInfoMapper  = models.NewNodeInfoMapper()
}