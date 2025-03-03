package server

import (
	"errors"
	"io"
	"sync"
	"time"
	pb "your-module-name/proto/proto"
)

type Status int
type MessageType int32

const (
	Prepare Status = iota
	Doing
	Finished
)
const (
	HaveTasks       int32 = iota //只有心跳，有任务正在进行
	NoTasks                      //心跳同时没有任务进行
	FinishTasks                  //任务完成主动发送的
	DistributeTasks              //server->client，分发任务
)

type server struct {
	pb.UnimplementedCrawlerServiceServer
	pb.UnimplementedSubscribeServer

	mutex      sync.Mutex
	taskIndex  int
	taskQueue  []string
	taskStatus map[string]Status
	streams    map[pb.Subscribe_SubscribeServer]bool

	//长时间没收到消息发送heartbeat
	timeLimit time.Duration

	//存储正在处理的url，如果client出故障将该client处理的url交给其他client
	DoingTasks map[pb.Subscribe_SubscribeServer][]string
}

func (s *server) Subscribe(in *pb.NodeInfo, stream pb.Subscribe_SubscribeServer) error {
	if s.streams[stream] {
		return errors.New("already subscribed")
	}
	s.streams[stream] = true
	go s.StartClient(stream)
	return nil
}

func (s *server) StartClient(stream pb.Subscribe_SubscribeServer) {
	//从用户接受数据，每一段时间接受心跳消息保证client正常运行
	go func() {
		for {
			var ResMsg pb.ResponseMessage
			err := stream.RecvMsg(ResMsg)
			//收到报错信息，停止连接并将对应的任务重新入队
			if err != nil {
				if err == io.EOF {
					//for key,str :=range s.streams{
					//	if str == stream{
					//		delete(s.streams, key)
					//		break
					//	}
					//}
				}
				s.streams[stream] = false
				break
			}
			if ResMsg.Success && ResMsg.MessageType == HaveTasks { //收到心跳并且client正在处理任务

			} else if ResMsg.Success && ResMsg.MessageType == NoTasks { //收到心跳并且client此时没有任务
				//空闲，发任务
				s.mutex.Lock()
				if len(s.taskQueue) == 0 {
					continue
				}
				s.DistributeTasks(stream)
				s.mutex.Unlock()
			} else if ResMsg.Success && ResMsg.MessageType == FinishTasks { //收到client完成任务的消息，调用数据合并模块

			}
		}
	}()

}

func (s *server) DistributeTasks(stream pb.Subscribe_SubscribeServer) {
	//为client分配任务
	msg := pb.ResponseMessage{
		MessageType: DistributeTasks,
		Message:     time.Now().Format("2006-01-02 15:04:05") + " distribute tasks",
		Url:         s.taskQueue[0],
		Success:     true,
	}
	err := stream.Send(&msg)
	if err != nil {
		return
	}
	s.taskStatus[s.taskQueue[0]] = Doing
	s.DoingTasks[stream] = append(s.DoingTasks[stream], s.taskQueue[0])
	s.taskQueue = s.taskQueue[1:]
}
