package crawler

import (
	"sync"

	log2 "github.com/apex/log"
	config2 "github.com/PunyGoood/DCS/system/config"
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/errors"
	"github.com/crawlab-team/crawlab/core/interfaces"
	model "github.com/crawlab-team/crawlab/core/models/models/v2"

	"github.com/PunyGoood/DCS/system/grpc/server"
	"github.com/PunyGoood/DCS/system/interfaces"
	"github.com/PunyGoood/DCS/system/node/config"
	"github.com/PunyGoood/DCS/system/task/manager"
	"github.com/PunyGoood/DCS/system/task/scheduler"
	"github.com/apex/log"

type Service struct {
	// dependencies
	NodeCfg      interfaces.NodeCfgfig
	schedulerSvc *scheduler.Service
	cron         *cron.Cron
	syncLock     bool

	// basic configs
	cfgPath string
}

//Schedule
func (svc *Service) Init(id int64, opts *interfaces.SpiderRunOptions) (taskIds []int64, err error) {

	s, err := service.NewModelService[model.Crawler]().GetById(id)
	if err != nil {
		return nil, err
	}

	return svc.schedule(s, opts)
}

//schedule Tasks
func (svc *Service) schedule(s *model.Crawler, opts *interfaces.CrawlerOptions) (taskIds []int64, err error) {

	t := &model.Task{


	}
	t.SetId()



	return taskIds, nil
}

func (svc *Service) getNodeIds(opts *interfaces.CrawlerOptions) (nodeIds []int64, err error) {

	return nodeIds, nil
}


func newCrawlerService() (svc *Service, err error) {
	svc := &Service{
		NodeCfg: config.GetNodeCfgfig(),
		cfgPath: config2.GetConfigPath(),
	}
	svc.schedulerSvc, err = scheduler.GetTaskSchedulerService()
	if err != nil {
		return nil, err
	}


	// validate 
	if !svc.NodeCfg.IsMaster() {
		//待定
		return nil, trace.TraceError()

	}

	return svc, nil
}

var svcV2 *ServiceV2

func GetCrawlerService() (svc *Service, err error) {
	sync.Once.Do(func() {
		svc, err = newCrawlerService()
		if err != nil {
			log2.Errorf("[GetSpiderAdminServiceV2] error: %v", err)
		}
	})
	if err != nil {
		return nil, err
	}
	return svc, nil
}
