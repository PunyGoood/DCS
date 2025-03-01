package manager

import (
	"github.com/PunyGoood/DCS/system/interfaces"
	"time"
)

type Option func(s interfaces.TaskManagerService)

func WithConfigPath(path string) Option {
	return func(s interfaces.TaskManagerService) {
		s.SetConfigPath(path)
	}
}

func WithExitWatchDuration(duration time.Duration) Option {
	return func(s interfaces.TaskManagerService) {
		s.SetExitWatchDuration(duration)
	}
}

func WithReportInterval(interval time.Duration) Option {
	return func(s interfaces.TaskManagerService) {
		s.SetReportInterval(interval)
	}
}

func WithCancelTimeout(timeout time.Duration) Option {
	return func(s interfaces.TaskManagerService) {
		s.SetCancelTimeout(timeout)
	}
}

type RunnerOption func(r interfaces.TaskRunner)

func WithSubscribeTimeout(timeout time.Duration) RunnerOption {
	return func(r interfaces.TaskRunner) {
		r.SetSubscribeTimeout(timeout)
	}
}
