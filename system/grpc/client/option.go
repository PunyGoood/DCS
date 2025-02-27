package client

import (
	"github.com/PunyGoood/DCS/system/interfaces"
	"time"
)

type Option func(client interfaces.GrpcClient)

func WithConfigPath(path string) Option {
	return func(cc interfaces.GrpcClient) {
		cc.SetConfigPath(path)
	}
}

func WithAddress(address interfaces.Address) Option {
	return func(cc interfaces.GrpcClient) {
		cc.SetAddress(address)
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(cc interfaces.GrpcClient) {
		cc.SetTimeout(timeout)
	}
}

func WithSubscribeType(subscribeType string) Option {
	return func(cc interfaces.GrpcClient) {
		cc.SetSubscribeType(subscribeType)
	}
}

func WithManageMessage(manageMessage bool) Option {
	return func(cc interfaces.GrpcClient) {
		cc.SetManageMessage(manageMessage)
	}
}

type PoolOption func(p interfaces.GrpcClientPool)

func WithPoolConfigPath(path string) PoolOption {
	return func(cc interfaces.GrpcClientPool) {
		cc.SetConfigPath(path)
	}
}

func WithPoolSize(size int) PoolOption {
	return func(cc interfaces.GrpcClientPool) {
		cc.SetSize(size)
	}
}
