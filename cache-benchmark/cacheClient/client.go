package cacheClient

// 用于保存操作的类型、键、值、错误
type Cmd struct {
	Name  string
	Key   string
	Value string
	Error error
}

// client接口，run提供单个cmd运行，pipelinedrun提供多个cmd运行
type Client interface {
	Run(*Cmd)
	PipelinedRun([]*Cmd)
}

// 根据typ参数决定生成不同的符合client接口的结构体
func New(typ, server string) Client {
	if typ == "redis" {
		return newRedisClient(server)
	}
	if typ == "http" {
		return newHTTPClient(server)
	}
	if typ == "tcp" {
		return newTCPClient(server)
	}
	panic("unknown client type " + typ)
}
