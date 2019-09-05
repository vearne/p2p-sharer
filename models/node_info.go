package models

import (
	"github.com/mohae/deepcopy"
	"sync"
)

type NodeInfoMapper struct {
	// nodeId -> addr
	Inner  map[Node]string `json:"inner"`
	Locker sync.RWMutex    `json:"-"`
}

func NewNodeInfoMapper() *NodeInfoMapper {
	m := NodeInfoMapper{}
	m.Inner = make(map[Node]string)
	return &m
}

func (m *NodeInfoMapper) Add(node Node, addr string) {
	m.Locker.Lock()
	defer m.Locker.Unlock()
	m.Inner[node] = addr
}

func (m *NodeInfoMapper) Get(node Node) string {
	m.Locker.RLock()
	defer m.Locker.RUnlock()
	return m.Inner[Node(node)]
}

func (m *NodeInfoMapper) Clone() *NodeInfoMapper {
	m.Locker.RLock()
	defer m.Locker.RUnlock()
	return deepcopy.Copy(m).(*NodeInfoMapper)
}
