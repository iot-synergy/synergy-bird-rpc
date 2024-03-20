package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	MonDb MonConfig
}

type MonConfig struct {
	Url    string
	DbName string
}
