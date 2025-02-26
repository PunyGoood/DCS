package auth

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func GetAutoUnaryClientInterceptor(nodeCfg interfaces.NodeConfig) grpc.UnaryClientInterceptor {

	md := metadata.Pairs(cons.GrpcHeaderAuth, nodeCfg.GetAuthKey())

	func(ctx context.Context, method string, req, reply interface{}, cc *ClientConn, invoker UnaryInvoker, opts ...CallOption) error {

		ctx = metadata.NewOutgoingContext(context.Background(), md)

		return invoker(ctx, method, req, reply, cc, opts...)
	}

}
