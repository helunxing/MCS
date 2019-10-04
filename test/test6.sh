# pipeline为10
./cache-benchmark -type tcp -n 100000 -r 100000 -t get -P 10
# pipeline为100
./cache-benchmark -type tcp -n 100000 -r 100000 -t get -P 100
# 50个客户端，效果显著
./cache-benchmark -type tcp -n 100000 -r 100000 -t get -c 50
# redis50个客户端
redis-benchmark -c 50 -n 100000 -d 1000 -t get -r 100000 -P 1
