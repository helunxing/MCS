package tcp

import (
	"bufio"
	"io"
	"log"
	"net"
)

type result struct {
	v []byte
	e error
}

// 此类函数先创建一个channel，再将其放入resultCh
func (s *Server) get(ch chan chan *result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c
	k, e := s.readKey(r)
	if e != nil {
		c <- &result{nil, e}
		return
	}
	go func() {
		v, e := s.Get(k)
		c <- &result{v, e}
	}()
}

func (s *Server) set(ch chan chan *result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c
	k, v, e := s.readKeyAndValue(r)
	if e != nil {
		c <- &result{nil, e}
		return
	}
	go func() {
		c <- &result{nil, s.Set(k, v)}
	}()
}

func (s *Server) del(ch chan chan *result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c
	k, e := s.readKey(r)
	if e != nil {
		c <- &result{nil, e}
		return
	}
	go func() {
		c <- &result{nil, s.Del(k)}
	}()
}

// 接受resultCh并发送响应，resultCh用于传递结果channel
func reply(conn net.Conn, resultCh chan chan *result) {
	defer conn.Close()
	for {
		// 接收resultCh
		c, open := <-resultCh
		// 已关闭则退出方法
		if !open {
			return
		}
		// 从chan中接收result
		r := <-c
		// 向连接返回响应
		e := sendResponse(r.v, r.e, conn)
		if e != nil {
			log.Println("close connection due to error:", e)
			return
		}
	}
}

// 用于处理来自tcp连接的客户请求
func (s *Server) process(conn net.Conn) {
	r := bufio.NewReader(conn)
	resultCh := make(chan chan *result, 5000)
	defer close(resultCh)
	// 新协程处理请求并响应
	go reply(conn, resultCh)
	for {
		op, e := r.ReadByte()
		if e != nil {
			if e != io.EOF { // 写入结束
				log.Println("close connection due to error:", e)
			}
			return
		}
		if op == 'S' {
			s.set(resultCh, r)
		} else if op == 'G' {
			s.get(resultCh, r)
		} else if op == 'D' {
			s.del(resultCh, r)
		} else {
			log.Println("close connection due to invalid operation:", op)
			return
		}
	}
}
