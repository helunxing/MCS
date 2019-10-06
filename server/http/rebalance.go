package http

import (
	"bytes"
	"net/http"
)

// 实现http.Handler接口
type rebalanceHandler struct {
	*Server
}

// POST方法调用
func (h *rebalanceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	// 新协程完成平衡
	go h.rebalance()
}

func (h *rebalanceHandler) rebalance() {
	// 获取缓存遍历器
	s := h.NewScanner()
	defer s.Close()
	c := &http.Client{}
	// 循环遍历键值，若不是本节点处理则插入删除
	for s.Scan() {
		k := s.Key()
		n, ok := h.ShouldProcess(k)
		if !ok {
			r, _ := http.NewRequest(http.MethodPut, "http://"+n+":12345/cache/"+k, bytes.NewReader(s.Value()))
			c.Do(r)
			h.Del(k)
		}
	}
}

func (s *Server) rebalanceHandler() http.Handler {
	return &rebalanceHandler{s}
}
