package client

import (
	"context"
	"sync"
	"time"

	"github.com/apex/log"
	"google.golang.org/grpc"
)

type GrpcClient struct {
	conn *grpc.ClientConn
	once sync.Once
	err  error

	//
	address interfaces.Address
	timeout time.Duration

	//
}

func Start(sou *GrpcClient) (err error) {

	sou.once.Do(func() {

		err = sou.Connect()
		if err != nil {
			return
		}

		err = sou.Register()
		if err != nil {
			return
		}

	})

	return err
}

func Connect(sou *GrpcClient) (err error) {
	ops := func() error {

		Address := sou.address.String()

		ctx, cancel := context.WithTimeout(context.Background(), sou.timeout)
		defer cancel()
	}
}

func Register(sou *GrpcClient) (err error) {

}

func (sou *GrpcClient) Context() (ctx context.Context, cancel context.CancelFunc) {

	return context.WithTimeout(context.Background(), sou.timeout)
}

func (sou *GrpcClient) Stop() (err error) {

	if sou.conn == nil {
		return nil
	}

	// grpc server address
	Address := sou.address.String()

	// shut connection
	if err := sou.conn.Close(); err != nil {
		return err
	}
	log.Infof("grpc client disconnected from %s", Address)

	return nil
}

func newGrpcClient() (sou *GrpcClient) {

	return client
}

var client *GrpcClient
var clientOnce sync.Once

func GetGrpcClient() *GrpcClient {
	clientOnce.Do(func() {
		client = newGrpcClient()
	})
	return client
}
