package cluster

import (
	"io/ioutil"
	"time"

	"github.com/hashicorp/memberlist"
	"stathat.com/c/consistent"
)

type Node interface {
	// 应当处理的节点
	ShouldProcess(key string) (string, bool)
	// 整个集群的节点列表
	Members() []string
	// 本节点地址
	Addr() string
}

type node struct {
	*consistent.Consistent
	addr string
}

func (n *node) Addr() string {
	return n.addr
}

func New(addr, cluster string) (Node, error) {
	// 创建默认LAN设置结构体指针
	conf := memberlist.DefaultLANConfig()
	// 名字 监听地址 放弃输出
	conf.Name = addr
	conf.BindAddr = addr
	conf.LogOutput = ioutil.Discard
	// 创建memberlist.Memberlist结构体
	l, e := memberlist.Create(conf)
	if e != nil {
		return nil, e
	}
	if cluster == "" {
		cluster = addr
	}
	clu := []string{cluster}
	// 加入指定集群
	_, e = l.Join(clu)
	if e != nil {
		return nil, e
	}
	// 创建consisitent.Consistent结构体指针
	circle := consistent.New()
	// 将每个节点的虚拟节点数量置为256个
	circle.NumberOfReplicas = 256
	// 每秒将集群节点列表更新到circle中
	go func() {
		for {
			m := l.Members()
			nodes := make([]string, len(m))
			for i, n := range m {
				nodes[i] = n.Name
			}
			circle.Set(nodes)
			time.Sleep(time.Second)
		}
	}()
	return &node{circle, addr}, nil
}

func (n *node) ShouldProcess(key string) (string, bool) {
	addr, _ := n.Get(key)
	return addr, addr == n.addr
}
