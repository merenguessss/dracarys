package config

import (
	"testing"
)

func TestConfig(t *testing.T) {
	err := SetPath("test_dracarys.yml")
	if err != nil {
		t.Error(err)
		return
	}

	client, err := GetClient()
	if err != nil {
		t.Errorf("get client error %v", err)
		return
	}
	if client.SerializerType != "proto" {
		t.Error("config fail")
	}

	server, err := GetServer()
	if err != nil {
		t.Errorf("get server error %v", err)
		return
	}
	if server.ServerName != "dracarys" {
		t.Error("config fail")
	}

	if config.Plugin.Selector.SelectorName != "default" ||
		config.Plugin.Selector.Timeout != "2s" {
		t.Error("config fail")
	}
}

func TestPathExist(t *testing.T) {
	if !pathExist("test_dracarys.yml") {
		t.Error("eeee")
	}
	if pathExist("sssss") {
		t.Error("eee")
	}
}
