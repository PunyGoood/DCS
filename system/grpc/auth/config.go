package auth

import (
	"context"

	"github.com/PunyGoood/DCS/system/cons"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func GetAuthUnaryClientInterceptor(nodeCfg interfaces.NodeConfig) grpc.UnaryClientInterceptor {
	// set auth key
	md := metadata.Pairs(cons.GrpcHeaderAuth, nodeCfg.GetAuthKey())

	func(ctx context.Context, method string, req, reply interface{}, cc *ClientConn, invoker UnaryInvoker, opts ...CallOption) error {

		ctx = metadata.NewOutgoingContext(context.Background(), md)

		return invoker(ctx, method, req, reply, cc, opts...)
	}

}

func GetAutStreamInterceptor(nodeCfg interfaces.NodeConfig) grpc.StreamClientInterceptor {

	md := metadata.Pairs(cons.GrpcHeaderAuthorization, nodeCfg.GetAuthKey())

	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = metadata.NewOutgoingContext(context.Background(), md)
		//opts = append(opts, grpc.Header(&header))
		s, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			return nil, err
		}
		return s, nil
	}
}

func GetAuthTokenFunc(nodeCfg interfaces.NodeConfig) grpc_auth.AuthFunc {
	return func(ctx context.Context) (ctx2 context.Context, err error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.ErrorGrpcUnauthorized
		}

		res, ok := md[constants.GrpcHeaderAuthorization]
		if !ok {
			return ctx, errors.ErrorGrpcUnauthorized
		}
		if len(res) != 1 {
			return ctx, errors.ErrorGrpcUnauthorized
		}
		authKey := res[0]

		// validate
		svrAuthKey := nodeCfg.GetAuthKey()
		if authKey != svrAuthKey {
			return ctx, errors.ErrorGrpcUnauthorized
		}

		return ctx, nil
	}
}
