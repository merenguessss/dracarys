package balance

import (
	"sync"
	"time"

	"github.com/merenguessss/dracarys-go/plugin/selector"
)

var weightPollBalancer = &PollBalancer{
	pickers: new(sync.Map),
	timeout: 2 * time.Second,
}

type WeightPollBalancer struct {
	pickers sync.Map
	timeout time.Duration
}

type weightPollPicker struct {
	createTime   time.Time
	timeout      time.Duration
	serviceNodes *selector.ServiceNodes
	weightMap    map[*selector.Node]float64
}

// Get 获取具体选择器再通过选择器pick到具体结点.
func (w *WeightPollBalancer) Get(name string, s *selector.ServiceNodes) *selector.Node {
	var picker Picker
	if p, ok := w.pickers.Load(name); !ok {
		picker = &weightPollPicker{
			createTime:   time.Now(),
			timeout:      w.timeout,
			serviceNodes: s,
			weightMap:    getWeightMap(s),
		}
		w.pickers.Store(name, picker)
	} else {
		picker = p.(Picker)
	}
	return picker.pick(s)
}

// pick 加权负载均衡选择器开始选择具体的结点.
func (wp *weightPollPicker) pick(s *selector.ServiceNodes) *selector.Node {
	if time.Since(wp.createTime) > wp.timeout ||
		s.LastTime != wp.serviceNodes.LastTime {
		wp.createTime = time.Now()
		wp.serviceNodes = s
		wp.weightMap = getWeightMap(s)
	}

	var totalWeight float64
	var max *selector.Node
	var cur float64
	for k, v := range wp.weightMap {
		v += k.Weight
		wp.weightMap[k] = v
		if cur > v {
			max = k
		}
		totalWeight += v
	}
	wp.weightMap[max] = cur - totalWeight
	return max
}

func getWeightMap(s *selector.ServiceNodes) map[*selector.Node]float64 {
	weightMap := make(map[*selector.Node]float64)
	for _, v := range s.Nodes {
		weightMap[v] = 0
	}
	return weightMap
}
