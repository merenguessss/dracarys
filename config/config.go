package config

import (
	"errors"
	"io/ioutil"
	"os"
	"sync/atomic"
	"time"

	"github.com/merenguessss/dracarys/client"
	"github.com/merenguessss/dracarys/log"
	"github.com/merenguessss/dracarys/plugin"
	"github.com/merenguessss/dracarys/server"

	"gopkg.in/yaml.v3"
)

var configStream []byte
var path string
var config = newDefault()
var ErrorConfigFileNotExist = errors.New("config file not exist")

// Setting 整体设置结构体.
// Client 客户端默认设置.
// Server 服务端默认设置.
type Setting struct {
	Client *client.Options `yaml:"client"`
	Server *server.Options `yaml:"server"`
	Plugin *plugin.Options `yaml:"plugin"`
	Logger *log.Options    `yaml:"log"`
	isLoad int32
}

var newDefault = func() *Setting {
	return &Setting{
		Client: &client.Options{
			SerializerType: "json",
			NetWork:        "tcp",
			CodecType:      "proto",
		},
		Server: &server.Options{
			Network:         "tcp",
			KeepAlivePeriod: 200 * time.Second,
			SerializerType:  "json",
			CodecType:       "proto",
		},
		Plugin: plugin.DefaultConfig(),
		Logger: log.DefaultOptions,
	}
}

// readConfigBytes 读取文件中的默认配置字节流.
// 默认为当前目录下的dracarys.yml文件.
// 可以在系统环境变量中设置DRACARYS_CONFIG,设置为要使用的路径.
// 例如: classpath:/rpc/dracarys.yml
func readConfigBytes() ([]byte, error) {
	if path == "" {
		path = "dracarys.yml"
		if c := os.Getenv("DRACARYS_CONFIG"); c != "" {
			path = c
		}
	}

	if !pathExist(path) {
		return nil, ErrorConfigFileNotExist
	}
	return ioutil.ReadFile(path)
}

// pathExist 判断是否存在config文件.
// 如果存在返回true, 不存在返回false.
func pathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// SetPath 设置自定义的path路径.
// 例如: classpath:/rpc/dracarys.yml
func SetPath(p string) error {
	var err error
	path = p
	configStream, err = readConfigBytes()
	if err != nil {
		return err
	}
	return nil
}

// loadConfig 加载Config函数,通过yaml.v3包从读取的字节流中获取config内容.
func loadConfig() error {
	if config == nil {
		config = newDefault()
	}
	// 已经加载过配置,直接返回.
	if config.isLoad == 1 {
		return nil
	}

	pluginFactory := plugin.New(config.Plugin)
	config.Client.PluginFactory = pluginFactory
	config.Server.PluginFactory = pluginFactory

	var err error
	if configStream == nil {
		configStream, err = readConfigBytes()
		if err != nil {
			// 如果config文件不存在,直接返回默认配置.
			if err == ErrorConfigFileNotExist {
				return nil
			}
			return err
		}
	}

	err = yaml.Unmarshal(configStream, config)
	if err != nil {
		return err
	}
	log.DefaultLogger = log.New(config.Logger)
	atomic.AddInt32(&config.isLoad, 1)
	return nil
}

// GetClient 获取Client端默认配置.
func GetClient() (*client.Options, error) {
	err := loadConfig()
	if err != nil {
		return nil, err
	}
	return config.Client, nil
}

// GetServer 获取Server端默认配置.
func GetServer() (*server.Options, error) {
	err := loadConfig()
	if err != nil {
		return nil, err
	}
	return config.Server, nil
}
