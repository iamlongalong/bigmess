package rmap

import (
	"strings"
	"sync"
)

func NewSubTreeRoot() *SubTreeRoot {
	return &SubTreeRoot{
		treeNode: NewSubTreeNode(true),
		splitKey: ".", // 先默认所有都一样，以后可提供选择
	}
}

type SubTreeRoot struct {
	treeNode *SubTreeNode
	splitKey string
}

func (r *SubTreeRoot) Add(key string, v interface{}) {

	r.treeNode.Add(r.subs(key), v)
}

func (r *SubTreeRoot) subs(key string) []string {
	subkeys := strings.Split(key, r.splitKey)

	return append([]string{""}, subkeys...)
}

func (r *SubTreeRoot) Remove(key string, v interface{}) {
	r.treeNode.Remove(r.subs(key), v)
}

func (r *SubTreeRoot) Get(key string) []interface{} {
	return r.treeNode.GetVals(r.subs(key))
}

type SubTreeNode struct {
	mu   sync.Mutex
	root bool

	Vals *sliceVals

	Children map[string]*SubTreeNode
}

func NewSubTreeNode(root bool) *SubTreeNode {
	return &SubTreeNode{root: root, Vals: newSliceVals(), Children: make(map[string]*SubTreeNode)}
}

func (n *SubTreeNode) Remove(subs []string, v interface{}) {
	if len(subs) == 0 {
		n.Vals.Remove(v)
		return
	}

	n.mu.Lock()
	c := n.Children[subs[0]]
	n.mu.Unlock()

	if c != nil {
		c.Remove(subs[1:], v)
	}
}

func (n *SubTreeNode) Add(subs []string, v interface{}) {

	if len(subs) == 0 {
		n.Vals.Add(v)
		return
	}

	n.mu.Lock()
	cn := n.Children[subs[0]]
	if cn == nil {
		cn = NewSubTreeNode(false)
		n.Children[subs[0]] = cn
	}
	n.mu.Unlock()

	cn.Add(subs[1:], v)
}

// GetVals 获取包含 prefix 的值
func (n *SubTreeNode) GetVals(keys []string) []interface{} {
	v := make([]interface{}, 0)
	v = append(v, n.Vals.Get()...)

	if len(keys) == 0 {
		return v
	}

	n.mu.Lock()
	cn := n.Children[keys[0]]
	n.mu.Unlock()

	if cn == nil {
		return v
	}

	return append(v, cn.GetVals(keys[1:])...)
}

func newSliceVals() *sliceVals {
	return &sliceVals{
		val:     make([]interface{}, 0),
		realLen: 0,
	}
}

type sliceVals struct {
	mu sync.Mutex

	val     []interface{}
	realLen int
}

func (s *sliceVals) Get() []interface{} {
	res := make([]interface{}, s.realLen)
	s.mu.Lock()
	copy(res, s.val[0:s.realLen])
	s.mu.Unlock()

	return res
}

func (s *sliceVals) Add(v interface{}) {
	s.mu.Lock()
	if len(s.val) == s.realLen {
		s.val = append(s.val, v)
	} else {
		s.val[s.realLen] = v
	}

	s.realLen++

	defer s.mu.Unlock()
}

func (s *sliceVals) Remove(v interface{}) {
	s.mu.Lock()

	for i := 0; i < s.realLen; i++ {
		if s.val[i] == v {
			s.val[i], s.val[s.realLen-1] = s.val[s.realLen-1], nil

			s.realLen--
			i--
		}
	}

	s.mu.Unlock()
}
