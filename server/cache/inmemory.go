package cache

import (
	"sync"
	"time"
)

type value struct {
	v       []byte
	created time.Time
}

type inMemoryCache struct {
	// 更改为value结构体
	c     map[string]value
	mutex sync.RWMutex
	Stat
	// 未提供成员名只提供结构体类型，内嵌 语法，用来实现继承。
	// 调用时使用类型指代。
	ttl time.Duration
}

// 设置缓存值
func (c *inMemoryCache) Set(k string, v []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.c[k] = value{v, time.Now()}
	c.add(k, v)
	return nil
}

// 获取值
func (c *inMemoryCache) Get(k string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.c[k].v, nil
}

// 删除键值
func (c *inMemoryCache) Del(k string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	v, exist := c.c[k]
	if exist {
		delete(c.c, k)
		c.del(k, v.v)
	}
	return nil
}

// 获取状态
func (c *inMemoryCache) GetStat() Stat {
	return c.Stat
}

// 创建结构体实例
func newInMemoryCache(ttl int) *inMemoryCache {
	c := &inMemoryCache{make(map[string]value), sync.RWMutex{}, Stat{}, time.Duration(ttl) * time.Second}
	if ttl > 0 {
		// 新协程定期清理过期元素
		go c.expirer()
	}
	return c
}

func (c *inMemoryCache) expirer() {
	for {
		time.Sleep(c.ttl)
		c.mutex.RLock()
		for k, v := range c.c {
			c.mutex.RUnlock()
			if v.created.Add(c.ttl).Before(time.Now()) {
				c.Del(k)
			}
			c.mutex.RLock()
		}
		c.mutex.RUnlock()
	}
}
