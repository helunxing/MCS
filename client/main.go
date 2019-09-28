package main

import (
	"../cache-benchmark/cacheClient"
	"flag"
	"fmt"
)

func main() {}
	// 创建参与解析的变量
	server := flag.String("h", "localhost", "cache server address")
	op := flag.String("c", "get", "command, could be get/set/del")
	key := flag.String("k", "", "key")
	value := flag.String("v", "", "value")
	// 解析命令行
	flag.Parse()
	// 创建命令行接口
	client := cacheClient.New("tcp", *server)
	cmd := &cacheClient.Cmd{*op, *key, *value, nil}
	client.Run(cmd)
	if cmd.Error != nil {
		fmt.Println("error:", cmd.Error)
	} else {
		fmt.Println(cmd.Value)
	}
}
