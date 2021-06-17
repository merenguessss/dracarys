package selector

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/merenguessss/dracarys-go/log"
)

type Consul struct {
	config       *api.Config
	client       *api.Client
	opts         *Options
	registration *api.AgentServiceRegistration
	once         sync.Once
}

func init() {
	Register("consul", consul)
}

var consul = &Consul{}

// LoadConfig 加载consul配置.
func (c *Consul) LoadConfig(opts *Options) error {
	c.opts = opts
	var err error
	config := api.DefaultConfig()
	if c.opts.Address != "" {
		config.Address = c.opts.Address
	}
	config.Scheme = c.opts.Scheme
	config.HttpClient = http.DefaultClient

	c.config = config
	c.client, err = api.NewClient(config)
	return err
}

// RegisterClient 将client注册到服务中心.
func (c *Consul) RegisterClient(name, address string) error {
	// 使用consul ,将不对client进行注册.
	return nil
}

// RegisterService 将service注册到服务中心.
// 此处与RegisterClient代码略微重复,但需要适应接口.
func (c *Consul) RegisterService(name, address string) error {
	if name == "" || address == "" {
		return ErrorEmptyNameAddress
	}

	host := getHost(address)
	if host == "" {
		return ErrorEmptyNameAddress
	}
	c.opts.Host = host

	registration := new(api.AgentServiceRegistration)
	c.registration = registration
	if c.opts.EnableHeartbeat {
		c.RegisterHeartbeat()
	}

	// 直接使用address作为ID,避免重复注册.
	registration.ID = name + "/" + address
	registration.Address = address
	registration.Name = name
	// 初始化权重,黑白名单.
	registration.Meta = map[string]string{
		"weights":   "100",
		"blacklist": "",
		"whitelist": "",
		"last_time": fmt.Sprintf("%d", time.Now().Unix()),
	}

	if err := c.client.Agent().ServiceRegister(registration); err != nil {
		return err
	}
	return nil
}

// Select 通过 Service.Name 查询具体结点.
func (c *Consul) Select(name string) (*ServiceNodes, error) {
	filter := "Service == \"" + name + "\""
	services, err := c.client.Agent().ServicesWithFilter(filter)
	if err != nil {
		return nil, err
	}

	nodes := &ServiceNodes{
		Name:  name,
		Nodes: make([]*Node, 0),
	}
	for serviceName, v := range services {
		node, err := newNode(serviceName, v.Address, v.Meta)
		if err != nil {
			log.Error(err)
			continue
		}

		nodes.addNode(node)
	}
	return nodes, nil
}

// RegisterHeartbeat 注册并监听心跳.
func (c *Consul) RegisterHeartbeat() {
	check := new(api.AgentServiceCheck)
	host := c.opts.Host
	port := c.opts.Port
	addr := host + ":" + port
	path := "/actuator/health"

	check.HTTP = "http://" + addr + path
	check.Timeout = c.opts.Timeout
	check.Interval = c.opts.Interval
	check.DeregisterCriticalServiceAfter = c.opts.DeregisterCriticalServiceAfter
	c.registration.Check = check

	// 监听心跳.
	c.once.Do(func() {
		go func() {
			http.HandleFunc(path, func(writer http.ResponseWriter, request *http.Request) {
				if request.Method == "GET" {
					writer.WriteHeader(http.StatusOK)
				}
			})
			_ = http.ListenAndServe(addr, nil)
		}()
	})
}
