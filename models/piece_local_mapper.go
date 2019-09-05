package models

import (
	"github.com/mohae/deepcopy"
	"sync"
)

type LocalInfo struct {
	PieceInfo

	FilePath string
}

type PieceLocalMapper struct {
	Inner  map[Piece]*LocalInfo `json:"inner"`
	Locker sync.RWMutex         `json:"-"`
}

func NewPieceLocalMapper() *PieceLocalMapper {
	var m PieceLocalMapper
	m.Inner = make(map[Piece]*LocalInfo)
	return &m
}

func (m *PieceLocalMapper) Get(piece Piece) (info *LocalInfo, ok bool) {
	m.Locker.RLock()
	defer m.Locker.RUnlock()
	info, ok = m.Inner[piece]
	return info, ok
}

func (m *PieceLocalMapper) Add(piece Piece, info *LocalInfo) {
	m.Locker.Lock()
	defer m.Locker.Unlock()
	m.Inner[piece] = info
}

// deepcopy
func (m *PieceLocalMapper) Clone() *PieceLocalMapper {
	m.Locker.RLock()
	defer m.Locker.RUnlock()
	return deepcopy.Copy(m).(*PieceLocalMapper)
}
