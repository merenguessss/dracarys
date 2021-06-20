package selector

import (
	"strconv"
	"sync"
	"time"
)

// ServiceNodes 单个服务.
type ServiceNodes struct {
	// Name 服务名.
	Name string
	// Nodes 服务结点.
	Nodes []*Node
	// 结点数.
	Length int
	// LastTime 服务最后更新时间.
	LastTime time.Duration
	mu       sync.RWMutex
}

// Node 服务结点.
type Node struct {
	// Key   键: 服务名.
	Key string
	// Value 值: 结点ip地址.
	Value string
	// Blacklist 黑名单.
	Blacklist string
	// Whitelist 白名单.
	Whitelist string
	// Weight 权重.
	Weight float64
	// LastTime 最新更新时间.
	LastTime time.Duration
}

// addNode 加入一个结点.
func (s *ServiceNodes) addNode(node *Node) {
	lastTime := node.LastTime
	if s.LastTime < lastTime {
		s.LastTime = lastTime
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.Length++
	s.Nodes = append(s.Nodes, node)
}

func newNode(key, value string, meta map[string]string) (*Node, error) {
	var err error
	var lastTime int64
	if t, ok := meta["last_time"]; ok {
		lastTime, err = strconv.ParseInt(t, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	var weight float64
	if w, ok := meta["weight"]; ok {
		weight, err = strconv.ParseFloat(w, 32)
		if err != nil {
			return nil, err
		}
	}

	return &Node{
		Key:       key,
		Value:     value,
		Weight:    weight,
		Blacklist: meta["blacklist"],
		Whitelist: meta["whitelist"],
		LastTime:  time.Duration(lastTime),
	}, nil
}
