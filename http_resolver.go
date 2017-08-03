package tdc

import (
	"errors"
	"github.com/cocotyty/httpclient"
	"github.com/golang/glog"
	"sync"
	"time"
)

type nodeInfo struct {
	Name     string
	Version  uint64
	Listener Listener
}
type httpResourceSolver struct {
	Address   string
	Env       string
	Tick      time.Duration
	WatchNode []nodeInfo
	lock      sync.RWMutex
}

func NewHTTPResourceSolver(Address string, Env string, Tick time.Duration) *httpResourceSolver {
	solver := &httpResourceSolver{
		Address:   Address,
		Env:       Env,
		Tick:      Tick,
		WatchNode: []nodeInfo{},
	}
	go solver.StartWatch()
	return solver
}

type httpResourceSolverResp struct {
	Error   string `json:"error"`
	Exist   bool   `json:"exist"`
	Data    []byte `json:"data"`
	Version uint64 `json:"version"`
}

func (s *httpResourceSolver) ConfigurationRefByName(name string, fn Listener) (data []byte, err error) {
	for i := 0; i < 3; i++ {
		data, err = s.configurationRefByName(name, fn)
		if err == nil {
			break
		}
		time.Sleep(time.Duration(i) * 500 * time.Millisecond)
	}
	return data, err
}

func (s *httpResourceSolver) configurationRefByName(name string, fn Listener) (data []byte, err error) {
	resp := &httpResourceSolverResp{}

	err = httpclient.Get(s.Address).Query("name", name).Query("env", s.Env).Send().JSON(resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, errors.New(resp.Error)
	}
	if !resp.Exist {
		return nil, errors.New(name + " not exist")
	}

	s.lock.Lock()
	s.WatchNode = append(s.WatchNode, nodeInfo{
		Name:     name,
		Version:  resp.Version,
		Listener: fn,
	})
	s.lock.Unlock()

	return resp.Data, nil
}
func (s *httpResourceSolver) nodes() []nodeInfo {
	s.lock.RLock()
	nodes := make([]nodeInfo, len(s.WatchNode))
	copy(nodes, s.WatchNode)
	s.lock.RUnlock()
	return nodes
}
func (s *httpResourceSolver) StartWatch() {
	if s.Tick == time.Duration(0) {
		s.Tick = 1 * time.Second
	}
	for {
		nodes := s.nodes()
		for _, node := range nodes {
			resp := &httpResourceSolverResp{}
			err := httpclient.Get(s.Address).Query("name", node.Name).Query("env", s.Env).Send().JSON(resp)
			if err != nil {
				glog.Error("error -> watch conf node", node.Name, err)
				continue
			}
			if resp.Error != "" {
				glog.Error("error -> watch conf node", node.Name, resp.Error)
				continue
			}
			if node.Version != resp.Version {
				if node.Listener != nil {
					glog.Info("event -> node changed", node.Name, resp.Version)
					node.Listener(node.Name, resp.Data, resp.Version, resp.Exist)
				}
			}
			node.Version = resp.Version
		}
		time.Sleep(s.Tick)
	}
}
