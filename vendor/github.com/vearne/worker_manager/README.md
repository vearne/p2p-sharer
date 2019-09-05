# 可以切换到工程目录下
# 执行
```
go get github.com/vearne/worker_manager
```

```
package main

import (
	"context"
	"github.com/gin-gonic/gin"
	manager "github.com/vearne/worker_manager"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 1. init some worker
	wm := prepareAllWorker()

	// 2. start
	wm.Start()

	// 3. register grace exit
	GracefulExit(wm)

	// 4. block and wait
	wm.Wait()
}

func GracefulExit(wm *manager.WorkerManager) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM)
	switch <-ch {
	case syscall.SIGTERM, syscall.SIGINT:
		log.Println("got signal")
		wm.Stop()
		break
	}
}

func prepareAllWorker() *manager.WorkerManager {
	wm := manager.NewWorkerManager()
	// load worker
	WorkerCount := 2
	for i := 0; i < WorkerCount; i++ {
		wm.AddWorker(NewLoadWorker())
	}
	// web server
	wm.AddWorker(NewWebServer())

	return wm
}

// some worker

type LoadWorker struct {
	RunningFlag bool // is running? true:running false:stoped
	ExitedFlag  bool //  Exit Flag
	ExitChan    chan struct{}
}

func NewLoadWorker() *LoadWorker {
	worker := &LoadWorker{RunningFlag: true, ExitedFlag: false}
	worker.ExitChan = make(chan struct{})
	return worker
}

func (worker *LoadWorker) Start() {
	log.Println("[start]LoadWorker")
	for worker.RunningFlag {
		select {
		case <-time.After(1 * time.Minute):
			//do some thing
			log.Println("LoadWorker do something")
			time.Sleep(time.Second * 3)

		case <-worker.ExitChan:
			log.Println("LoadWorker execute exit logic")
		}

	}
	worker.ExitedFlag = true
}

func (worker *LoadWorker) Stop() {
	log.Println("LoadWorker exit...")
	worker.RunningFlag = false
	close(worker.ExitChan)
	for !worker.ExitedFlag {
		time.Sleep(50 * time.Millisecond)
	}
	log.Println("[end]LoadWorker")
}

type WebServer struct {
	Server *http.Server
}

func NewWebServer() *WebServer {
	return &WebServer{}
}

func (worker *WebServer) Start() {
	log.Println("[start]WebServer")

	ginHandler := gin.Default()
	ginHandler.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/plain", []byte("hello world!"))
	})
	worker.Server = &http.Server{
		Addr:           ":8080",
		Handler:        ginHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	worker.Server.ListenAndServe()
}

func (worker *WebServer) Stop() {
	log.Println("WebServer exit...")
	cxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// gracefull exit web server
	err := worker.Server.Shutdown(cxt)
	if err != nil {
		log.Printf("shutdown error, %v", err)
	}
	log.Println("[end]WebServer exit")
}
```

```
go build main.go
# 启动服务
./main
# 服务退出, 发出SIGTERM信号，服务优雅退出
# 请求自行替换pid的值
kill -15 <pid> 
```
output
```
2019/08/23 14:28:41 [start]LoadWorker
2019/08/23 14:28:41 [start]LoadWorker
2019/08/23 14:28:41 [start]WebServer
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /                         --> main.(*WebServer).Start.func1 (3 handlers)
2019/08/23 14:28:58 got signal
2019/08/23 14:28:58 WebServer exit...
2019/08/23 14:28:58 [end]WebServer exit
2019/08/23 14:28:58 LoadWorker exit...
2019/08/23 14:28:58 LoadWorker execute exit logic
2019/08/23 14:28:58 LoadWorker exit...
2019/08/23 14:28:58 LoadWorker execute exit logic
2019/08/23 14:28:58 [end]LoadWorker
2019/08/23 14:28:58 [end]LoadWorker
```
