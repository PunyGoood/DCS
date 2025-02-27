package interfaces

import (
	"context"
	grpc "github.com/PunyGoood/DCS/grpc"
	"time"
)

type GrpcClient interface {

	WithConfigPath
	SetAddress(Address)
	SetTimeout(time.Duration)
	SetSubscribeType(string)
	Context() (context.Context, context.CancelFunc)
	NewRequest(interface{}) *grpc.Request
	Restart() error
	IsStarted() bool
	IsClosed() bool
	Err() error

}
