package manager

import (
	"context"
	"fmt"
	"time"
)

type TaskService struct {

	c *grpcclient.GrpcClient
	cfgSvc interfaces.NodeConfigService

}





func (s *TaskService) Fetch() {
	for {
		// 假设这里有获取任务的逻辑
		fmt.Println("Fetching tasks")
		time.Sleep(10 * time.Second)
	}
}


func start(s *TaskService) {

	if s.c.isClosed != nil {
		err := s.c.Start()
		if err != nil {
			return
		}
	}

	go svc.ReportStatus()
	go svc.Fetch()

}

func (s *TaskService) Run(taskId primitive.ObjectID) (err error) {
	return s.run(taskId)
}


func (s *TaskService) Cancel(taskId primitive.ObjectID) (err error) {
	r, err := s.getRunner(taskId)
	if err != nil {
		return err
	}
	if err := r.Cancel(); err != nil {
		return err
	}
	return nil
}

func (s *TaskService) Fetch() {
	ticker := time.NewTicker(s.fetchInterval)
	for {
		// wait
		<-ticker.C

		// current node
		n, err := svc.GetCurrentNode()
		if err != nil {
			continue
		}

		// skip if node is not active or enabled
		if !n.Active || !n.Enabled {
			continue
		}


		if s.getRunnerCount() >= n.MaxRunners {
			continue
		}

		if s.stopped {
			ticker.Stop()
			return
		}

		// fetch task
		tid, err := svc.fetch()
		if err != nil {
			trace.PrintError(err)
			continue
		}

		// skip if no task id
		if tid.IsZero() {
			continue
		}

		// run task
		if err := svc.run(tid); err != nil {
			trace.PrintError(err)
			t, err := svc.GetTaskById(tid)
			if err != nil && t.Status != constants.TaskStatusCancelled {
				t.Error = err.Error()
				t.Status = constants.TaskStatusError
				t.SetUpdated(t.CreatedBy)
				_ = client.NewModelServiceV2[models2.TaskV2]().ReplaceById(t.Id, *t)
				continue
			}
			continue
		}
	}
}