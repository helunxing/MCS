package main

import (
	"flag"
	"log"

	"./cache"
	"./cluster"
	"./http"
	"./tcp"
)

func main() {
	// 接受命令行参数的值
	typ := flag.String("type", "inmemory", "cache type")
	ttl := flag.Int("ttl", 30, "cache time to live")
	node := flag.String("node", "127.0.0.1", "node address")
	clus := flag.String("cluster", "", "cluster address")
	flag.Parse()
	log.Println("type is", *typ)
	log.Println("ttl is", *ttl)
	log.Println("node is", *node)
	log.Println("cluster is", *clus)
	// 新建cache
	c := cache.New(*typ, *ttl)
	// 新建cluster.Node接口
	n, e := cluster.New(*node, *clus)
	if e != nil {
		panic(e)
	}
	// 根据cache在新协程建立tcp.Server对象
	go tcp.New(c, n).Listen()
	// 根据cache创建http.Server对象
	http.New(c, n).Listen()
}
