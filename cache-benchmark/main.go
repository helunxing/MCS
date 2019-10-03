package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"./cacheClient"
)

type statistic struct {
	count int
	time  time.Duration
}

type result struct {
	getCount    int
	missCount   int
	setCount    int
	statBuckets []statistic
}

func (r *result) addStatistic(bucket int, stat statistic) {
	if bucket > len(r.statBuckets)-1 {
		newStatBuckets := make([]statistic, bucket+1)
		copy(newStatBuckets, r.statBuckets)
		r.statBuckets = newStatBuckets
	}
	s := r.statBuckets[bucket]
	s.count += stat.count
	s.time += stat.time
	r.statBuckets[bucket] = s
}

func (r *result) addDuration(d time.Duration, typ string) {
	bucket := int(d / time.Millisecond)
	r.addStatistic(bucket, statistic{1, d})
	if typ == "get" {
		r.getCount++
	} else if typ == "set" {
		r.setCount++
	} else {
		r.missCount++
	}
}

func (r *result) addResult(src *result) {
	for b, s := range src.statBuckets {
		r.addStatistic(b, s)
	}
	r.getCount += src.getCount
	r.missCount += src.missCount
	r.setCount += src.setCount
}

//
func run(client cacheClient.Client, c *cacheClient.Cmd, r *result) {
	expect := c.Value
	start := time.Now()
	// 使用cacheClient.Client.Run
	client.Run(c)
	d := time.Now().Sub(start)
	resultType := c.Name
	if resultType == "get" {
		if c.Value == "" {
			resultType = "miss"
		} else if c.Value != expect {
			panic(c)
		}
	}
	// 将结果记录
	r.addDuration(d, resultType)
}

func pipeline(client cacheClient.Client, cmds []*cacheClient.Cmd, r *result) {
	expect := make([]string, len(cmds))
	for i, c := range cmds {
		if c.Name == "get" {
			expect[i] = c.Value
		}
	}
	start := time.Now()
	client.PipelinedRun(cmds)
	d := time.Now().Sub(start)
	for i, c := range cmds {
		resultType := c.Name
		if resultType == "get" {
			if c.Value == "" {
				resultType = "miss"
			} else if c.Value != expect[i] {
				fmt.Println(expect[i])
				panic(c.Value)
			}
		}
		r.addDuration(d, resultType)
	}
}

// count是请求数量
func operate(id, count int, ch chan *result) {
	// 创建cacheClient.Client接口
	client := cacheClient.New(typ, server)
	// 创建命令数组
	cmds := make([]*cacheClient.Cmd, 0)
	valuePrefix := strings.Repeat("a", valueSize)
	// 返回结果
	r := &result{0, 0, 0, make([]statistic, 0)}
	//使用循环记录到result结构体数组指针r中
	for i := 0; i < count; i++ {
		var tmp int
		if keyspacelen > 0 {
			tmp = rand.Intn(keyspacelen)
		} else {
			tmp = id*count + i
		}
		key := fmt.Sprintf("%d", tmp)
		value := fmt.Sprintf("%s%d", valuePrefix, tmp)
		name := operation
		if operation == "mixed" {
			if rand.Intn(2) == 1 {
				name = "set"
			} else {
				name = "get"
			}
		}
		c := &cacheClient.Cmd{name, key, value, nil}
		// 使用pipelen
		if pipelen > 1 {
			cmds = append(cmds, c)
			if len(cmds) == pipelen {
				pipeline(client, cmds, r)
				cmds = make([]*cacheClient.Cmd, 0)
			}
		} else {
			// 不使用pipelen，仅把1个command传给服务端，并保存结果
			run(client, c, r)
		}
	}
	// 将剩余命令执行
	if len(cmds) != 0 {
		pipeline(client, cmds, r)
	}
	//将结果发到channel
	ch <- r
}

var typ, server, operation string
var total, valueSize, threads, keyspacelen, pipelen int

// 解析命令行参数，对rand种子初始化
func init() {
	flag.StringVar(&typ, "type", "redis", "cache server type")
	flag.StringVar(&server, "h", "localhost", "cache server address")
	flag.IntVar(&total, "n", 1000, "total number of requests")
	flag.IntVar(&valueSize, "d", 1000, "data size of SET/GET value in bytes")
	flag.IntVar(&threads, "c", 1, "number of parallel connections")
	flag.StringVar(&operation, "t", "set", "test set, could be get/set/mixed")
	flag.IntVar(&keyspacelen, "r", 0, "keyspacelen, use random keys from 0 to keyspacelen-1")
	flag.IntVar(&pipelen, "P", 1, "pipeline length")
	flag.Parse()
	fmt.Println("type is", typ)
	fmt.Println("server is", server)
	fmt.Println("total", total, "requests")
	fmt.Println("data size is", valueSize)
	fmt.Println("we have", threads, "connections")
	fmt.Println("operation is", operation)
	fmt.Println("keyspacelen is", keyspacelen)
	fmt.Println("pipeline length is", pipelen)

	rand.Seed(time.Now().UnixNano())
}

func main() {
	// 创建channel以传递结果
	ch := make(chan *result, threads)
	res := &result{0, 0, 0, make([]statistic, 0)}
	start := time.Now()
	// 开启指定个数个协程执行操作
	for i := 0; i < threads; i++ {
		go operate(i, total/threads, ch)
	}
	// 接收结果
	for i := 0; i < threads; i++ {
		res.addResult(<-ch)
	}
	d := time.Now().Sub(start)
	totalCount := res.getCount + res.missCount + res.setCount
	// 输出结果
	fmt.Printf("%d records get\n", res.getCount)
	fmt.Printf("%d records miss\n", res.missCount)
	fmt.Printf("%d records set\n", res.setCount)
	fmt.Printf("%f seconds total\n", d.Seconds())
	statCountSum := 0
	statTimeSum := time.Duration(0)
	// 分
	for b, s := range res.statBuckets {
		if s.count == 0 {
			continue
		}
		statCountSum += s.count
		statTimeSum += s.time
		fmt.Printf("%d%% requests < %d ms\n", statCountSum*100/totalCount, b+1)
	}
	fmt.Printf("%d usec average for each request\n", int64(statTimeSum/time.Microsecond)/int64(statCountSum))
	fmt.Printf("throughput is %f MB/s\n", float64((res.getCount+res.setCount)*valueSize)/1e6/d.Seconds())
	fmt.Printf("rps is %f\n", float64(totalCount)/float64(d.Seconds()))
}
