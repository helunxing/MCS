package http

import (
	"encoding/json"
	"log"
	"net/http"
)

type clusterHandler struct {
	*Server
}

// 实现server接口
func (h *clusterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// 返回节点列表
	m := h.Members()
	b, e := json.Marshal(m)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func (s *Server) clusterHandler() http.Handler {
	return &clusterHandler{s}
}
