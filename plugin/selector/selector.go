package selector

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Selector interface {
	RegisterClient(name, address string) error
	RegisterService(name, address string) error
	RegisterHeartbeat()
	Select(string) ([]*Node, error)
	LoadConfig(*Options) error
}

var ErrorEmptyNameAddress = errors.New("empty name or address")

var selectorMap = make(map[string]Selector)

func Register(name string, s Selector) {
	if selectorMap == nil {
		selectorMap = make(map[string]Selector)
	}
	selectorMap[name] = s
}

func Get(name string) Selector {
	if v, ok := selectorMap[name]; ok {
		return v
	}
	return selectorMap["default"]
}

// getHost 通过address解析出host.
func getHost(addr string) string {
	index := strings.LastIndex(addr, ":")
	if index == 0 || index == -1 || !verifyIPv4Host(addr[:index]) {
		return ""
	}
	return addr[:index]
}

// verifyIPv4Host 验证host是否符合IPv4.
func verifyIPv4Host(host string) bool {
	strArr := strings.Split(host, ".")
	if len(strArr) != 4 {
		return false
	}
	for _, str := range strArr {
		num, err := strconv.Atoi(str)
		if err != nil || num > 255 || num < 0 {
			return false
		}
		newStr := fmt.Sprint(num)
		if str != newStr {
			return false
		}
	}
	return true
}
