package tdc

import (
	"encoding/json"
	"net/http"
)

type ResolverQuery interface {
	Query(name, env string) ([]byte, error)
}

type Server struct {
	ResolverQuery ResolverQuery
}

func (s *Server) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	name, env := req.FormValue("name"), req.FormValue("env")
	data, err := s.ResolverQuery.Query(name, env)

	jsonResp, _ := json.Marshal(httpResourceSolverResp{
		Exist: err == nil,
		Data:  data,
	})
	resp.Header().Set("Content-Type", "application/json")
	resp.Write(jsonResp)
}
