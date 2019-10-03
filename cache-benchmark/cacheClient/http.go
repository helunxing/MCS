package cacheClient

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type httpClient struct {
	*http.Client
	server string
}

// 实现get操作
func (c *httpClient) get(key string) string {
	// 根据缓存服务的地址和key拼成get请求并发送
	resp, e := c.Get(c.server + key)
	if e != nil {
		log.Println(key)
		panic(e)
	}
	if resp.StatusCode == http.StatusNotFound {
		return ""
	}
	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}
	b, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		panic(e)
	}
	// 返回结果
	return string(b)
}

// 实现set操作
func (c *httpClient) set(key, value string) {
	// put方法
	req, e := http.NewRequest(http.MethodPut,
		c.server+key, strings.NewReader(value))
	if e != nil {
		log.Println(key)
		panic(e)
	}
	resp, e := c.Do(req)
	if e != nil {
		log.Println(key)
		panic(e)
	}
	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}
}

// 根据name成员决定调用get还是set
func (c *httpClient) Run(cmd *Cmd) {
	if cmd.Name == "get" {
		cmd.Value = c.get(cmd.Key)
		return
	}
	if cmd.Name == "set" {
		c.set(cmd.Key, cmd.Value)
		return
	}
	panic("unknown cmd name " + cmd.Name)
}

// 创建httpClient结构体的指针
func newHTTPClient(server string) *httpClient {
	client := &http.Client{Transport: &http.Transport{MaxIdleConnsPerHost: 1}}
	return &httpClient{client, "http://" + server + ":12345/cache/"}
}

// 未实现pipeline
func (c *httpClient) PipelinedRun([]*Cmd) {
	panic("httpClient pipelined run not implement")
}
