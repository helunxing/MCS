package main

import (
	"./cache"
	"./http"
	"./tcp"
)

func main() {
	ca := cache.New("inmemory")
	// 新协程建立tcp.Server对象
	go tcp.New(ca).Listen()
	http.New(ca).Listen()
}
