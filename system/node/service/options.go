package service

import (
	"time"

	"github.com/PunyGoood/DCS/system/interfaces"
	"github.com/crawlab-team/crawlab/core/interfaces"
)

type Option func(svc interfaces.NodeService)

func WithConfigPath(path string) Option {
	return func(svc interfaces.NodeService) {
		svc.SetConfigPath(path)
	}
}

func WithAddress(address interfaces.Address) Option {
	return func(svc interfaces.NodeService) {
		svc.SetAddress(address)
	}
}
func WithSetStop() option {
	return func(svc interfaces.NodeService) {
		svc2, ok := svc.(interfaces.NodeMasterService)
		if ok {
			svc2.SetStop()
		}
	}
}

func WithHeartbeatInterval(duration time.Duration) Option {
	return func(svc interfaces.NodeService) {
		svc2, ok := svc.(interfaces.NodeWorkerService)
		if ok {
			svc2.SetHeartbeatInterval(duration)
		}
	}
}
