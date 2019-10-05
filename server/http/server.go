package http

import (
	"net/http"

	"../cache"
	"../cluster"
)

type Server struct {
	cache.Cache
	cluster.Node
}

// 开始监听
func (s *Server) Listen() {
	http.Handle("/cache/", s.cacheHandler())
	http.Handle("/status", s.statusHandler())
	http.Handle("/cluster", s.clusterHandler())
	http.ListenAndServe(s.Addr()+":12345", nil)
}

// 创建并返回结构体
func New(c cache.Cache, n cluster.Node) *Server {
	return &Server{c, n}
}
