package engine

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/vearne/p2p-sharer/config"
	"github.com/vearne/p2p-sharer/models"
	"github.com/vearne/p2p-sharer/resource"
	"log"
	"net/http"
	"time"
)

type TrackerWebServer struct {
	Server *http.Server
}

func NewTrackerWebServer() *TrackerWebServer {
	return &TrackerWebServer{}
}

func (worker *TrackerWebServer) Start() {
	log.Println("[start]NodeWebServer")

	appconfig := config.GetOpts()

	ginHandler := NewTrackerRouter()
	worker.Server = &http.Server{
		Addr:           appconfig.Web.ListenAddress,
		Handler:        ginHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	worker.Server.ListenAndServe()
}

func (worker *TrackerWebServer) Stop() {
	cxt, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := worker.Server.Shutdown(cxt)
	if err != nil {
		log.Println("shutdown error", err)
	}
	log.Println("[end]NodeWebServer exit")
}

func NewTrackerRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/v1/heartBeat", heartBeat)
	r.GET("/v1/nodeList", nodeList)
	r.POST("/v1/report", report)
	r.GET("/v1/probe", func(c *gin.Context) {
		//NodeSession * models.Session
		//FilePieceMapper * models.PieceMapper
		//NodeInfoMapper * models.NodeInfoMapper
		resp := make(map[string]interface{})
		resp["NodeSession"] = resource.NodeSession.Clone()
		resp["FilePieceMapper"] = resource.FilePieceMapper.Clone()
		resp["NodeInfoMapper"] = resource.NodeInfoMapper.Clone()
		c.JSON(http.StatusOK, resp)
	})

	return r
}

type heartBeatParam struct {
	NodeId string `form:"nodeId" json:"nodeId" binding:"required"`
	Addr   string `form:"addr" json:"addr" binding:"required"`
}

func heartBeat(c *gin.Context) {
	var param heartBeatParam
	var err error
	err = c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrResponse{"E001", err.Error()})
		return
	}
	// 1. maintain session
	resource.NodeSession.HeartBeat(models.Node(param.NodeId))
	// 2. record node -> addr
	resource.NodeInfoMapper.Add(models.Node(param.NodeId), param.Addr)

	c.JSON(http.StatusOK, models.ErrResponse{"E000", "success"})
}

type nodeListParam struct {
	PieceID string `form:"pieceID" json:"pieceID" binding:"required"`
}

func nodeList(c *gin.Context) {
	var param nodeListParam
	var err error
	err = c.ShouldBindQuery(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrResponse{"E001", err.Error()})
		return
	}

	nodeList := resource.FilePieceMapper.GetNodeList(models.Piece(param.PieceID), 5)
	log.Println("web-nodeList", "PieceID", param.PieceID, "count of nodeList", len(nodeList))
	var resp models.NodeListResp
	resp.Nodes = make([]string, 0)
	for _, node := range nodeList {
		resp.Nodes = append(resp.Nodes, resource.NodeInfoMapper.Get(node))
	}
	c.JSON(http.StatusOK, resp)
}

func report(c *gin.Context) {
	var param models.ReportParam
	var err error
	err = c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrResponse{"E001", err.Error()})
		return
	}

	resource.FilePieceMapper.AddPiece(models.Node(param.NodeID), models.Piece(param.PieceID))
	log.Println("report", param.File,
		resource.NodeInfoMapper.Get(models.Node(param.NodeID)), param.Progress)

	if param.Progress >= 1.0 {
		log.Println("Download Finish", param.File,
			resource.NodeInfoMapper.Get(models.Node(param.NodeID)), param.Progress)
	}
	c.JSON(http.StatusOK, models.ErrResponse{"E000", "success"})
}
