## 分布式缓存

#### c1 基于http/rest的内存缓存

由于http解析，性能是redis的约四分之一

#### c2 基于tcp的内存缓存

性能有明显提升

#### c3 数据持久化

rocksdb：字节流形式键值对持久化。压缩系数调优。

cgo调用capi

rocksdb可以虚拟内存，还可以重启恢复。但语言转化有一定开销


#### c4 pipeline

不等待回复，连续发多个请求

响应的接受通常由其他协程负责

#### c5 rocksdb批量写入
多个set指令合并成一个，连续写入磁盘。无法得知写入失败，没有实时一致性。

channel类似polling。允许等待多事件同时发生

timer.Timer含一个成员chan Time C，触发后会发送time.Time结构体。重置时要先用stop确认，触发要先取出C中第一个

#### c6 异步操作提升读性能
服务端使用channel保证异步返回顺序

异步操作支出：channel和协程。rocksdb操作时间较慢，客户端请求密度较高。

本章功能使用channel实现，好处在于易读易实现地完成异步取结果的功能
