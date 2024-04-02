package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	MonDb    MonConfig
	EventRpc zrpc.RpcClientConf
}

type MonConfig struct {
	Url    string
	DbName string
}
