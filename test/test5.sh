# 没有pipeline的性能
./cache-benchmark -type tcp -n 100000 -r 100000 -t set

# 使用pipeline的性能
./cache-benchmark -type tcp -n 100000 -r 100000 -t set -P 3

# 多客户端
./cache-benchmark -type tcp -n 100000 -r 100000 -t set -c 50
# 多客户端的+pipeline
redis-benchmark -c 50 -n 100000 -d 1000 -t set -r 100000 -P 1
