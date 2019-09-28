package http

import (
	"net/http"

	"../cache"
)

type Server struct {
	cache.Cache
}

// 开始监听
func (s *Server) Listen() {
	http.Handle("/cache/", s.cacheHandler())
	http.Handle("/status", s.statusHandler())
	http.ListenAndServe(":12345", nil)
}

// 创建并返回结构体
func New(c cache.Cache) *Server {
	return &Server{c}
}
