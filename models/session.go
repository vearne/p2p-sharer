package models

import (
	"github.com/mohae/deepcopy"
	"sync"
	"time"
)

type Session struct {
	Map    map[Node]time.Time `json:"map"`
	Locker sync.Mutex         `json:"-"`
}

func NewSession() *Session {
	s := Session{}
	s.Map = make(map[Node]time.Time)
	return &s
}

func (s *Session) HeartBeat(node Node) {
	s.Locker.Lock()
	defer s.Locker.Unlock()
	s.Map[node] = time.Now()
}

func (s *Session) FindAndPurgeInvalid(expire time.Duration) []Node {
	s.Locker.Lock()
	defer s.Locker.Unlock()
	now := time.Now()
	result := make([]Node, 0)
	for node, last := range s.Map {
		if now.Sub(last) > expire {
			result = append(result, node)
			// delete
			delete(s.Map, node)
		}
	}
	return result
}

func (s *Session) Clone() *Session {
	s.Locker.Lock()
	defer s.Locker.Unlock()
	return deepcopy.Copy(s).(*Session)
}
