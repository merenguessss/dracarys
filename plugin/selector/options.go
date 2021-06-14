package selector

type Options struct {
	SelectorName      string `yaml:"name"`
	Address           string `yaml:"address"`
	Scheme            string `yaml:"scheme"`
	EnableHeartbeat   bool   `yaml:"heartbeat"`
	*HeartbeatOptions `yaml:"heartbeat"`
}

type HeartbeatOptions struct {
	Host                           string `yaml:"host"`
	Port                           string `yaml:"port"`
	Timeout                        string `yaml:"timeout"`
	Interval                       string `yaml:"interval"`
	DeregisterCriticalServiceAfter string `yaml:"deregister-critical-service-after"`
}
