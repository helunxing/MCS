package http

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// 内嵌server结构体指针
type cacheHandler struct {
	*Server
}

// 实现该方法意味着实现了handler接口
func (h *cacheHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 获取key
	key := strings.Split(r.URL.EscapedPath(), "/")[2]
	if len(key) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// 根据请求方法调用方法
	m := r.Method
	if m == http.MethodPut {
		b, _ := ioutil.ReadAll(r.Body)
		if len(b) != 0 {
			e := h.Set(key, b)
			if e != nil {
				log.Println(e)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
		return
	}
	if m == http.MethodGet {
		b, e := h.Get(key)
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(b) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write(b)
		return
	}
	if m == http.MethodDelete {
		e := h.Del(key)
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

// 根据server对象，建立cachehandler
func (s *Server) cacheHandler() http.Handler {
	return &cacheHandler{s}
}
