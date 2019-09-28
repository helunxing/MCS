# 设置值
./client -c set -k testkey -v testvalue
# 获取值
./client -c get -k testkey
# 读取状态
curl 127.0.0.1:12345/status
# 删除
./client -c del -k testkey
# 查看状态
curl 127.0.0.1:12345/status
# 测试性能
./cache-benchmark -type tcp -n 100000 -r 100000 -t set

./cache-benchmark -type tcp -n 100000 -r 100000 -t get
