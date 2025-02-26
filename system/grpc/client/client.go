package client

import (

	"sync"
	"time"

	grpc2 "github.com/PunyGoood/DCS/grpc/proto"


)

type GrpcClient struct {

	nodecfg *NodeConfig
	conn *grpc.ClientConn
	once sync.Once

	address string

	timeout time.Duration
	err error

	substream grpc2.NodeService_SubscribeClient

	NodeClient               grpc2.NodeServiceClient
	TaskClient               grpc2.TaskServiceClient
	ModelBaseServiceV2Client grpc2.ModelBaseServiceV2Client


}

func (cc *GrpcClientV2) Context() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), cc.timeout)
}

func (cc *GrpcClientV2) Connect() (err error) {
	op := func() error {
		address := cC.address.String()

		var opts []grpc.DialOption
		
		//opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		creds := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: false, 
		})
		opts = append(opts, grpc.WithTransportCredentials(creds))
		opts = append(opts, grpc.WithBlock())
		opts = append(opts, grpc.WithTimeout(cc.timeout))
		// 添加拦截器
		opts = append(opts, grpc.WithChainUnaryInterceptor(auth.GetAuthTokenUnaryChainInterceptor(cc.nodeCfg)))
		opts = append(opts, grpc.WithChainStreamInterceptor(auth.GetAuthTokenStreamChainInterceptor(cc.nodeCfg)))

		// 创建gRPC连接
		cc.conn, err = grpc.NewClient(address, opts...)
		if err != nil {
			_ = trace.TraceError(err)
			return errors.ErrorGrpcClientFailedToStart
		}
		log.Infof("[GrpcClient] grpc client connected to %s", address)
		return nil
	}
	// 使用指数退避策略重试连接
	return backoff.RetryNotify(op, backoff.NewExponentialBackOff(), utils.BackoffErrorNotify("grpc client connect"))
}


func (cc *GrpcClientV2) Register() {
	cc.NodeClient = grpc2.NewNodeServiceClient(cc.conn)
	cc.TaskClient = grpc2.NewTaskServiceClient(cc.conn)


}

func Start(cc *GrpcClient) {

	cc.once.Do(func () {
		err=cc.connect
		if(err!=nil) {
			return
		}

		cc.connect()

		cc.Register()
	
		cc.Subscribe()

	})


	return err
}

func (cc *GrpcClient) Stop() (err error) {

	if cc.conn == nil {
		return nil
	}

	address := cc.address.String()

	if err := cc.unsubscribe(); err != nil {
		return err
	}
	log.Infof("grpc client unsubscribed from %s", address)

	if err := cc.conn.Close(); err != nil {
		return err
	}
	log.Infof("grpc client disconnected from %s", address)

	return nil
}

func (cc *GrpcClient) NewRequest(d interface{}) (req *grpc2.Request) {
	return &grpc2.Request{
		NodeKey: c.nodeCfgSvc.GetNodeKey(),
		Data:    c.getRequestData(d),
	}
}

func (cc *GrpcClient) subscribe() (err error) {
	op := func() error {
		req := cc.NewRequest(&entity.NodeInfo{
			Key:      c.nodeCfg.GetNodeKey(),
			IsMaster: false,
		})
		cc.substream, err = c.NodeClient.Subscribe(context.Background(), req)
		if err != nil {
			return trace.TraceError(err)
		}

		log.Infof("[GrpcClient] grpc client subscribed to remote server")

		return nil
	}
	return backoff.RetryNotify(op, backoff.NewExponentialBackOff(), utils.BackoffErrorNotify("grpc client subscribe"))
}


func (cc *GrpcClient) unsubscribe() (err error) {
	req := cc.NewRequest(&entity.NodeInfo{
		Key:      cc.nodeCfg.GetNodeKey(),
		IsMaster: false,
	})
	if _, err = cc.NodeClient.Unsubscribe(context.Background(), req); err != nil {
		return trace.TraceError(err)
	}
	return nil
}

func newGrpcClient() (cc *GrpcClient){

	client := &GrpcClient{
		nodecfg: nodeconfig.GetANode(),
		timeout: 10 * time.Second,
		address: entity.NewAddress(&entity.AddressOptions{
			Host=cons.DefaultClientRemoteHost
			Port=cons.DefaultClientRemotePort
		})

	}

	return client

	
	
}


var client *GrpcClient
var once sync.Once


func GetGrpcClient() (*GrpcClient, error) {
	var err error
	once.Do(func() {
		client, err = newGrpcClient(config.GetNodeConfig())
	})
	return client, err
}

