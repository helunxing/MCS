package main

import (
	"flag"
	"log"

	"./cache"
	"./http"
	"./tcp"
)

func main() {
	// 接受命令行参数的值
	typ := flag.String("type", "inmemory", "cache type")
	flag.Parse()
	log.Println("type is", *typ)
	// 新建cache
	c := cache.New(*typ)
	// 根据cache在新协程建立tcp.Server对象
	go tcp.New(c).Listen()
	// 根据cache创建http.Server对象
	http.New(c).Listen()
}
