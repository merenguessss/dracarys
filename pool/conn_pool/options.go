package conn_pool

import "time"

type Options struct {
	// 最大连接数
	MaxIdle int
	// 核心连接数
	CoreIdle int
	// 最大连接数
	MaxActive int
	// 超过最大连接数是否等待
	Wait bool
	// 空闲连接超时时间
	IdleTimeout time.Duration
	// 建立连接超时时间
	DialTimeout time.Duration
}

type Option func(*Options)
