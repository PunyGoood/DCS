package client

import (

	"sync"
	"time"
    "context"
    "crypto/tls"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "github.com/cenkalti/backoff"
	grpc2 "github.com/PunyGoood/DCS/grpc/proto"


)

type GrpcClient struct {

	nodeCfg *NodeConfig
	conn *grpc.ClientConn
	once sync.Once

	address string
 
	timeout time.Duration
	err error

	substream grpc2.NodeService_SubscribeClient
	msgCh  chan *grpc2.StreamMessage

	NodeClient               grpc2.NodeServiceClient
	TaskClient               grpc2.TaskServiceClient
	ModelBaseServiceV2Client grpc2.ModelBaseServiceV2Client


}

func (cc *GrpcClient) Context() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), cc.timeout)
}

func (cc *GrpcClientV2) Connect() (err error) {
	op := func() error {
		address := cc.address.String()

		var opts []grpc.DialOption
		
		//opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		creds := credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: false, 
		})
		opts = append(opts, grpc.WithTransportCredentials(creds))
		opts = append(opts, grpc.WithBlock())
		opts = append(opts, grpc.WithTimeout(cc.timeout))
		opts = append(opts, grpc.WithChainUnaryInterceptor(auth.GetAuthTokenUnaryChainInterceptor(cc.nodeCfg)))
		opts = append(opts, grpc.WithChainStreamInterceptor(auth.GetAuthTokenStreamChainInterceptor(cc.nodeCfg)))

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
		if err != nil {
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
		NodeKey: cc.nodeCfgSvc.GetNodeKey(),
		Data:    cc.GetRequestData(d),
	}
}

func (cc *GrpcClient) GetRequestData(d interface{}) (data []byte) {
	if d == nil {
		return data
	}
	switch d.(type) {
	case []byte:
		data = d.([]byte)
	default:
		var err error
		data, err = json.Marshal(d)
		if err != nil {
			panic(err)
		}
	}
	return data
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

func (cc *GrpcClient) IsClosed() (res bool) {
	if cc.conn != nil {
		return cc.conn.GetState() == connectivity.Shutdown
	}
	return false
}

func (cc *GrpcClient) manageStreamMessage() {
	log.Infof("[GrpcClient] start managing stream message...")
	for {
		// resubscribe if stream is set to nil
		if cc.stream == nil {
			if err := backoff.RetryNotify(cc.subscribe, backoff.NewExponentialBackOff(), utils.BackoffErrorNotify("grpc client subscribe")); err != nil {
				log.Errorf("subscribe")
				return
			}
		}

		//node_service内容
		msg, err := c.stream.Recv()

		// send stream message to channel
		cc.msgCh <- msg

		// receive stream message
		log.Debugf("[GrpcClient] received message: %v", msg)
		if err != nil {

			cc.err = err

			// end
			if err == io.EOF {
				log.Infof("[GrpcClient] received EOF signal, can not be connected")
				return
			}

			if cc.IsClosed() {
				return
			}

			// error
			trace.PrintError(err)
			c.stream = nil
			time.Sleep(1 * time.Second)
			continue
		}

		// reset error
		cc.err = nil
	}
}

func newGrpcClient() (cc *GrpcClient){

	client := &GrpcClient{
		nodecfg: nodeconfig.GetANode(),
		timeout: 10 * time.Second,
		address: entity.NewAddress(&entity.AddressOptions{
			Host : =cons.DefaultClientRemoteHost
			Port : =cons.DefaultClientRemotePort
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

