package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/vearne/p2p-sharer/config"
	"github.com/vearne/p2p-sharer/consts"
	"github.com/vearne/p2p-sharer/models"
	"github.com/vearne/p2p-sharer/resource"
	"github.com/vearne/p2p-sharer/utils"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const BuffSize = 1024

type NodeWebServer struct {
	Server *http.Server
}

func NewNodeWebServer() *NodeWebServer {
	return &NodeWebServer{}
}

func (worker *NodeWebServer) Start() {
	log.Println("[start]NodeWebServer")

	appconfig := config.GetOpts()

	ginHandler := NewNodeRouter()
	worker.Server = &http.Server{
		Addr:           appconfig.Web.ListenAddress,
		Handler:        ginHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	worker.Server.ListenAndServe()
}

func (worker *NodeWebServer) Stop() {
	cxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := worker.Server.Shutdown(cxt)
	if err != nil {
		log.Println("shutdown   error", err)
	}
	log.Println("[end]NodeWebServer exit")
}

func NewNodeRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/v1/task", recvTask)
	//r.GET("/v1/pieceData", middleware.ConcurrentLimit(3), dealPieceData)
	r.GET("/v1/pieceData", dealPieceData)
	r.GET("/v1/probe", func(c *gin.Context) {
		//PieceLocalMapper * models.PieceLocalMapper
		resp := make(map[string]interface{})
		resp["PieceLocalMapper"] = resource.PieceLocalMapper.Clone()
		c.JSON(http.StatusOK, resp)
	})

	return r
}

type taskParam struct {
	SeedFile string `form:"seedFile" json:"seedFile" binding:"required"`
}

func recvTask(c *gin.Context) {
	var param taskParam
	var err error
	err = c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrResponse{"E001", err.Error()})
		return
	}
	// 1. get seed file
	response, err := http.Get(param.SeedFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrResponse{"E001",
			"send file can't download"})
		return
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	var info models.SeedInfo
	json.Unmarshal(body, &info)

	// save seed file in download dir
	seedFilePath := filepath.Join(config.GetOpts().DownloadDir, info.FileName+consts.SeedPostfix)
	sf, err := os.Create(seedFilePath)
	defer sf.Close()
	if err != nil {
		log.Println(err)
	}
	bt, _ := json.MarshalIndent(info, "", "\t")
	sf.Write(bt)

	// 2. async execute
	resource.TaskChan <- &info
	c.JSON(http.StatusOK, models.ErrResponse{"E000", "success"})
}

type pieceDataParam struct {
	PieceID string `form:"pieceID" json:"pieceID" binding:"required"`
}

func dealPieceData(c *gin.Context) {
	var param pieceDataParam
	var err error
	err = c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrResponse{"E001", err.Error()})
		return
	}

	info, ok := resource.PieceLocalMapper.Get(models.Piece(param.PieceID))
	log.Println("dealPieceData", info.FilePath, info.Index, info.Length)
	data := bytes.NewBuffer(make([]byte, 0, config.GetOpts().PieceSize))
	if ok {
		fi, err := os.Open(info.FilePath)
		defer fi.Close()

		log.Println("dealPieceData", "---1---")
		if err != nil {
			log.Println("dealPieceData", err)
			goto DefaultError
		}
		pos := config.GetOpts().PieceSize * int64(info.Index)
		_, err = fi.Seek(pos, 0)
		log.Println("dealPieceData", "---2---")
		if err != nil {
			log.Println("dealPieceData", err)
			goto DefaultError
		}
		buffer := make([]byte, BuffSize)
		total := info.Length
		for total > 0 {
			n, err := fi.Read(buffer)
			if err != nil {
				log.Println("dealPieceData", err)
				if err == io.EOF {
					break
				} else {
					goto DefaultError
				}
			}
			size := utils.Min(n, total)
			data.Write(buffer[0:size])
			total = total - n
		}
		c.Data(http.StatusOK, "application/octet-stream", data.Bytes())
		return
	}

DefaultError:
	c.Data(http.StatusNoContent, "application/octet-stream", nil)

}
