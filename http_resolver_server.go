package tdc

import (
	"encoding/json"
	"net/http"
)

type ResolverQuery interface {
	Query(name, env string) ([]byte, uint64, bool, error)
}

type Server struct {
	ResolverQuery ResolverQuery
}

func (s *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	name, env := req.FormValue("name"), req.FormValue("env")
	data, version, exist, err := s.ResolverQuery.Query(name, env)
	solverResp := httpResourceSolverResp{
		Exist:   exist,
		Data:    data,
		Version: version,
	}
	if err != nil {
		solverResp.Error = err.Error()
	}
	jsonResp, _ := json.Marshal(solverResp)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(jsonResp)
}
