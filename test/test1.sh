# 查看状态
curl 127.0.0.1:12345/status
# 设置
curl -v 127.0.0.1:12345/cache/testkey -XPUT -dtestvalue
# 读取
curl 127.0.0.1:12345/cache/testkey
# 查看状态
curl 127.0.0.1:12345/status
# 删除
curl 127.0.0.1:12345/cache/testkey -XDELETE
# 查看状态  
curl 127.0.0.1:12345/status
# 测试性能
./cache-benchmark -type http -n 100000 -r 100000 -t set

./cache-benchmark -type http -n 100000 -r 100000 -t get

redis-benchmark -c 1 -n 100000 -d 1000 -t set,get -r 100000
