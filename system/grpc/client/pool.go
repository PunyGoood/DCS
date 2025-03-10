package client

import (
	"github.com/PunyGoood/DCS/system/interfaces"
	"github.com/emirpasic/gods/lists/arraylist"
	"math/rand"
)

type Pool struct {

	size    int
	cfgPath string

	// internals
	clients *arraylist.List
}

func (p *Pool) GetConfigPath() (path string) {
	return p.cfgPath
}

func (p *Pool) SetConfigPath(path string) {
	p.cfgPath = path
}

func (p *Pool) Init() (err error) {
	for i := 0; i < p.size; i++ {
		err := p.NewClient()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Pool) NewClient() (err error) {
	cc, err := NewClient()
	if err != nil {
		return trace.TraceError(err)
	}
	if err := cc.Start(); err != nil {
		return err
	}
	p.clients.Add(cc)
	return nil
}

//Round-Robin
func (p *Pool) GetClient_RR() (c interfaces.GrpcClient, err error) {
    var currentIdx int
    if p.currentIdx >= p.clients.Size() {
        p.currentIdx = 0
    }
    res, ok := p.clients.Get(p.currentIdx)
    if !ok {
        return nil, trace.TraceError(errors.ErrorGrpcClientNotExists)
    }
    p.currentIdx++
    c, ok = res.(interfaces.GrpcClient)
    if !ok {
        return nil, trace.TraceError(errors.ErrorGrpcInvalidType)
    }
    return c, nil
}

//
func (p *Pool) GetClient() (c interfaces.GrpcClient, err error) {
	idx := p.getRandomIndex()
	res, ok := p.clients.Get(idx)
	if !ok {
		return nil, trace.TraceError(errors.ErrorGrpcClientNotExists)
	}
	c, ok = res.(interfaces.GrpcClient)
	if !ok {
		return nil, trace.TraceError(errors.ErrorGrpcInvalidType)
	}
	return c, nil
}

func (p *Pool) SetSize(size int) {
	p.size = size
}

func (p *Pool) getRandomIndex() (idx int) {
	return rand.Intn(p.clients.Size())
}

func NewPool(opts ...PoolOption) (p interfaces.GrpcClientPool, err error) {

	p = &Pool{
		size:    1,
		clients: arraylist.New(),
	}

	// apply options
	for _, opt := range opts {
		opt(p)
	}

	// initialize
	if err := p.Init(); err != nil {
		return nil, err
	}

	return p, nil
}
