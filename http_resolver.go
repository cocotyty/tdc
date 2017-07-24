package tdc

import (
	"errors"
	"github.com/cocotyty/httpclient"
)

type HTTPResourceSolver struct {
	Address string
	Env     string
}
type httpResourceSolverResp struct {
	Exist bool   `json:"exist"`
	Data  []byte `json:"data"`
}

func (s *HTTPResourceSolver) ConfigurationRefByName(name string) (data []byte, err error) {
	resp := &httpResourceSolverResp{}
	err = httpclient.Get(s.Address).Query("name", name).Query("env", s.Env).Send().JSON(resp)
	if err != nil {
		return nil, err
	}
	if !resp.Exist {
		return nil, errors.New(name + " not exist")
	}
	return resp.Data, nil
}
