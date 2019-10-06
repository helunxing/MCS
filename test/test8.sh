# 启动节点
# ./server -node 1.1.1.1
# 插入键
./cache-benchmark -type tcp -n 10000 -d 1 -h 1.1.1.1
# 查看缓存状态
curl 1.1.1.1:12345/status
# 启动节点
# ./server -node 1.1.1.2 -cluster 1.1.1.1

curl 1.1.1.2:12345/status
# 再平衡
curl 1.1.1.1:12345/rebalance -XPOST
# 看状态
curl 1.1.1.1:12345/status

curl 1.1.1.2:12345/status
# 启动节点
# ./server -node 1.1.1.3 -cluster 1.1.1.2
# 再平衡
curl 1.1.1.1:12345/rebalance -XPOST

curl 1.1.1.2:12345/rebalance -XPOST

curl 1.1.1.1:12345/status

curl 1.1.1.2:12345/status

curl 1.1.1.3:12345/status
