package main

import (
	"./cache"
	"./http"
)

// 创建cache对象，根据其创建server对象
func main() {
	c := cache.New("inmemory")
	http.New(c).Listen()
}
