package balance

import (
	"sync"
	"time"

	"github.com/merenguessss/dracarys/plugin/selector"
)

var pollBalancer = &PollBalancer{
	pickers: new(sync.Map),
	timeout: 2 * time.Second,
}

// PollBalancer 轮询负载均衡器.
type PollBalancer struct {
	pickers *sync.Map
	timeout time.Duration
}

type pollPicker struct {
	lastIndex    int
	createTime   time.Time
	timeout      time.Duration
	serviceNodes *selector.ServiceNodes
}

// Get 选择合适的选择器来进行选择.
func (pb *PollBalancer) Get(s *selector.ServiceNodes) *selector.Node {
	var picker Picker
	if p, ok := pb.pickers.Load(s.Name); !ok {
		picker = &pollPicker{
			lastIndex:    -1,
			createTime:   time.Now(),
			timeout:      pb.timeout,
			serviceNodes: s,
		}
		pb.pickers.Store(s.Name, picker)
	} else {
		picker = p.(Picker)
	}
	return picker.pick(s)
}

// pick 从选择器中选择一个结果.
func (p *pollPicker) pick(s *selector.ServiceNodes) *selector.Node {
	// 判断选择器是否过期.
	if time.Since(p.createTime) > p.timeout ||
		s.LastTime != p.serviceNodes.LastTime {
		p.serviceNodes = s
		p.lastIndex = -1
		p.createTime = time.Now()
	}

	index := p.lastIndex + 1
	if index >= s.Length {
		index = 0
	}

	p.lastIndex = index
	return p.serviceNodes.Nodes[index]
}
