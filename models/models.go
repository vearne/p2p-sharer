package models

import (
	gu "github.com/vearne/golib/utils"
	"sync"
	"time"
)

type SeedInfo struct {
	FileName    string `json:"fileName"`
	TrackerAddr string `json:"trackerAddr"`
	// total length
	Length    int          `json:"length"`
	Pieces    []*PieceInfo `json:"pieces"`
	CreatedAt time.Time    `json:"createdAt"`
}

type PieceInfo struct {
	Index    int    `json:"index"`
	Length   int    `json:"length"`
	Checksum string `json:"checksum"`
}

type NodeSet *gu.StringSet
type PieceSet *gu.StringSet

// Piece is SHA1 of file-piece
type Piece  string
type Node  string

type ErrResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}


type ConcurentLimiter struct{
	Locker sync.Mutex
	Num int
}

func NewConcurentLimiter(n int) *ConcurentLimiter{
	return &ConcurentLimiter{Num:n}
}

func (c *ConcurentLimiter) TryEnter() bool{
	c.Locker.Lock()
	defer c.Locker.Unlock()
	if c.Num > 0{
		c.Num--
		return true
	}
	return false
}

func (c *ConcurentLimiter) Exit(){
	c.Locker.Lock()
	defer c.Locker.Unlock()
	c.Num++
}


type DownloadTask struct{
	FileName string
	Total int
	SuccessCount int
	WaitForDeal *gu.IntSet
	Pieces []*PieceInfo
}


type NodeListResp struct {
	Nodes []string `json:"nodes"`
}


type ReportParam struct {
	PieceID  string  `form:"pieceID" json:"pieceID" binding:"required"`
	NodeID   string  `form:"nodeID" json:"nodeID" binding:"required"`
	Progress float64 `form:"progress" json:"progress" binding:"required"`
	File string  `form:"file" json:"file" binding:"required"`
}