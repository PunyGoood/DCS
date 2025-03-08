package service

import (
	"github.com/PunyGoood/DCS/system/grpc/server"
	"github.com/PunyGoood/DCS/system/interfaces"
	"github.com/PunyGoood/DCS/system/node/config"
	"github.com/PunyGoood/DCS/system/task/manager"
	"github.com/PunyGoood/DCS/system/task/scheduler"
	"github.com/apex/log"

	"sync"
	"time"
)

type MasterService struct {
	// dependencies
	nodeCfg      interfaces.NodeConfig
	server       *server.GrpcServer //	system/grpc
	schedulerSvc *scheduler.Service
	managerSvc   *manager.Service

	// basic configs
	cfgPath         string
	address         interfaces.Address
	monitorInterval time.Duration
	stopped         bool
}

func (svc *MasterService) Start() {

	// start grpc server
	if err := svc.server.Start(); err != nil {
		panic(err)
	}

	// register to db
	if err := svc.Register(); err != nil {
		panic(err)
	}

	go svc.Monitor()
	go svc.managerSvc.Start()
	go svc.schedulerSvc.Start()
	svc.Wait()
	svc.Stop()
}

func (svc *MasterService) Wait() {

}

func (svc *MasterService) Stop() {
	_ = svc.server.Stop()
	log.Infof("master[%s] service has stopped", svc.GetConfigService().GetNodeKey())
}

func (svc *MasterService) Monitor() {
	log.Infof("master[%s] monitoring started", svc.GetConfigService().GetNodeKey())

}

func (svc *MasterService) GetConfigService() (nodeCfg interfaces.NodeConfig) {
	return svc.nodeCfg
}

func (svc *MasterService) GetConfigPath() (path string) {
	return svc.cfgPath
}

func (svc *MasterService) SetConfigPath(path string) {
	svc.cfgPath = path
}

func (svc *MasterService) GetAddress() (address interfaces.Address) {
	return svc.address
}

func (svc *MasterService) SetAddress(address interfaces.Address) {
	svc.address = address
}

func (svc *MasterService) SetMonitorInterval(duration time.Duration) {
	svc.monitorInterval = duration
}

func (svc *MasterService) Register() (err error) {
	nodeKey := svc.GetConfigService().GetNodeKey()
	nodeName := svc.GetConfigService().GetNodeName()

}

func (svc *MasterService) SetStop() {
	svc.stopped = true
}

func (svc *MasterService) GetServer() (svr interfaces.GrpcServer) {
	return svc.server
}

func (svc *MasterService) monitor() (err error) {

	// update master node status in db

}

func newMasterService() (res *MasterService, err error) {

	// set basic configs
	svc := &MasterService{
		cfgPath:         config2.GetConfigPath(),
		monitorInterval: 15 * time.Second,
		stopped:         false,
	}

	svc.nodeCfg = config.GetNodeConfig()

	// grpc server
	svc.server, err = server.GetGrpcServer()
	if err != nil {
		return nil, err
	}

	svc.schedulerSvc, err = scheduler.GetTaskSchedulerService()
	if err != nil {
		return nil, err
	}

	svc.handlerSvc, err = manager.GetTaskManagerService()
	if err != nil {
		return nil, err
	}

	// no need to init

	return svc, nil
}

var MasterService *MasterService

func GetMasterService() (res *MasterService, err error) {
	sync.Once.Do(func() {
		MasterService, err = newMasterService()
		if err != nil {
			log.Errorf("failed to get master service: %v", err)
		}
	})
	return MasterService, err

}
