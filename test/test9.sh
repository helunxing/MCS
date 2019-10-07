# 启动节点
# ./server
# 设置并查看状态
curl 127.0.0.1:12345/cache/a -XPUT -daa

curl 127.0.0.1:12345/cache/a

curl 127.0.0.1:12345/status

# wait 30 seconds
# 再次尝试
curl 127.0.0.1:12345/cache/a

curl 127.0.0.1:12345/status
