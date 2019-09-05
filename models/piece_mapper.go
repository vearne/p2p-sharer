package models

import (
	"github.com/mohae/deepcopy"
	gu "github.com/vearne/golib/utils"
	rc "github.com/vearne/randomchoice"
	"sync"
)

type PieceMapper struct {
	PN     map[Piece]NodeSet `json:"pn"`
	NP     map[Node]PieceSet `json:"np"`
	Locker sync.RWMutex      `json:"-"`
}

func NewPieceMapper() *PieceMapper {
	m := PieceMapper{}
	m.PN = make(map[Piece]NodeSet)
	m.NP = make(map[Node]PieceSet)
	return &m
}

func (m *PieceMapper) AddPiece(node Node, piece Piece) {
	m.Locker.Lock()
	defer m.Locker.Unlock()

	//1.  add node -> pieceSet
	var ok bool
	if _, ok = m.NP[node]; !ok {
		m.NP[node] = gu.NewStringSet()
	}

	(*gu.StringSet)(m.NP[node]).Add(string(piece))

	//2.  add piece -> nodeSet
	if _, ok = m.PN[piece]; !ok {
		m.PN[piece] = gu.NewStringSet()
	}
	(*gu.StringSet)(m.PN[piece]).Add(string(node))
}

func (m *PieceMapper) PurgeNode(node Node) {
	m.Locker.Lock()
	defer m.Locker.Unlock()

	_, ok := m.NP[node]
	if !ok {
		return
	}

	// delete piece-node
	pieceSet := (*gu.StringSet)(m.NP[node])
	for _, p := range pieceSet.ToArray() {
		nodeSet := (*gu.StringSet)(m.PN[Piece(p)])
		nodeSet.Remove(string(node))
		if nodeSet.Size() <= 0 {
			delete(m.PN, Piece(p))
		}
	}

	// delete node
	delete(m.NP, node)
}

func (m *PieceMapper) GetNodeList(piece Piece, n int) []Node {
	m.Locker.RLock()
	defer m.Locker.RUnlock()
	if _, ok := m.PN[piece]; !ok {
		return make([]Node, 0)
	}
	set := (*gu.StringSet)(m.PN[piece])
	total := set.Size()
	result := make([]Node, 0, n)
	all := set.ToArray()
	for _, idx := range rc.RandomChoice(total, n) {
		result = append(result, Node(all[idx]))
	}
	return result
}

func (m *PieceMapper) Clone() *PieceMapper {
	m.Locker.RLock()
	defer m.Locker.RUnlock()
	return deepcopy.Copy(m).(*PieceMapper)
}
