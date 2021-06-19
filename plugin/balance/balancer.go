package balance

import (
	"math/rand"
	"time"

	"github.com/merenguessss/dracarys-go/plugin/selector"
)

type Balancer interface {
	Get(*selector.ServiceNodes) *selector.Node
}

type Picker interface {
	pick(*selector.ServiceNodes) *selector.Node
}

type Options struct {
	BalancerName string `yaml:"name"`
}

const (
	Random     = "random"
	Poll       = "poll"
	WeightPoll = "weight_poll"
)

var balancerMap = make(map[string]Balancer)

func init() {
	Register(Random, randomBalancer)
	Register(Poll, pollBalancer)
	Register(WeightPoll, weightPollBalancer)
}

// Register 注册负载均衡器.
func Register(name string, b Balancer) {
	if balancerMap == nil {
		balancerMap = make(map[string]Balancer)
	}
	balancerMap[name] = b
}

func Get(name string) Balancer {
	if v, ok := balancerMap[name]; ok {
		return v
	}
	return balancerMap[Random]
}

var randomBalancer = &RandomBalancer{}

// RandomBalancer 获取随机结点的负载均衡.
type RandomBalancer struct {
}

// Get 通过时间戳随机生成一个随机数,
func (b *RandomBalancer) Get(s *selector.ServiceNodes) *selector.Node {
	rand.Seed(time.Now().Unix())
	randomNum := rand.Intn(s.Length)
	return s.Nodes[randomNum]
}
