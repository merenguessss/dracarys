package balance

import (
	"sync"
	"testing"
	"time"

	"github.com/merenguessss/dracarys/log"
	"github.com/merenguessss/dracarys/plugin/selector"
)

var t = time.Now().Unix()

var serviceNodes = &selector.ServiceNodes{
	Name:   "name",
	Length: 4,
	Nodes: []*selector.Node{
		{
			Key:    "name",
			Value:  "1",
			Weight: 1,
		},
		{
			Key:    "name",
			Value:  "2",
			Weight: 2,
		},
		{
			Key:    "name",
			Value:  "3",
			Weight: 3,
		},
		{
			Key:    "name",
			Value:  "4",
			Weight: 4,
		},
	},
	LastTime: time.Duration(t),
}

func TestRandomBalancer(T *testing.T) {
	testBalancer(Random, 100)
}

func TestPollBalancer(T *testing.T) {
	testBalancer(Poll, 100)
}

func TestWeightPollBalancer(T *testing.T) {
	testBalancer(WeightPoll, 100)
}

func testBalancer(name string, n int) {
	balancer := Get(name)
	m := map[string]int{
		"1": 0,
		"2": 0,
		"3": 0,
		"4": 0,
	}
	var lock sync.Mutex

	for i := 0; i < n; i++ {
		node := balancer.Get(serviceNodes)
		lock.Lock()
		m[node.Value] = m[node.Value] + 1
		lock.Unlock()
	}

	log.Info("balancer name: ", name)
	for k, v := range m {
		log.Infof("%s - %d", k, v)
	}
}
