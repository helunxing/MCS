package tcp

import (
	"net"

	"../cache"
	"../cluster"
)

type Server struct {
	cache.Cache
	cluster.Node //额外内嵌了该接口
}

func (s *Server) Listen() {
	l, e := net.Listen("tcp", s.Addr()+":12346")
	if e != nil {
		panic(e)
	}
	// 接受连接并处理
	for {
		c, e := l.Accept()
		if e != nil {
			panic(e)
		}
		go s.process(c)
	}
}

func New(c cache.Cache, n cluster.Node) *Server {
	return &Server{c, n}
}
