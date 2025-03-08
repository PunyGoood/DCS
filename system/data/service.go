package data

import (
	"fmt"

	"github.com/PunyGoood/DCS/system/errors"
	"github.com/PunyGoood/DCS/system/interfaces"
	"github.com/PunyGoood/DCS/system/models"
	"github.com/PunyGoood/DCS/system/trace"
)

func NewDataService(registryKey string, s *models.Spider) (svc interfaces.DataService, err error) {
	// Data service function
	var fn interfaces.DataServiceRegistryFn

	// set func default may be sql
	if registryKey == "" {
		fn = NewDataService_
	} else {
		// from registry
		reg := GetDataServiceRegistry()
		fn = reg.Get(registryKey)
		if fn == nil {
			return nil, errors.NewDataError(fmt.Sprintf("%s can not work", registryKey))
		}
	}

	// generate Data service
	svc, err := fn(s.ColId, s.DataSourceId)
	if err != nil {
		return nil, trace.TraceError(err)
	}

	return svc, nil
}


func GetResultService(spiderId int64) (svc2 interfaces.ResultService, err error) {  
	// model service
	modelSvc, err := service.GetService()
	if err != nil {
		return nil, trace.TraceError(err)
	}

	// spider
	s, err := modelSvc.GetSpiderById(spiderId)  // 
	if err != nil {
		return nil, trace.TraceError(err)
	}

	// store key
	storeKey := fmt.Sprintf("%d:%d", s.ColId, s.DataSourceId) 
