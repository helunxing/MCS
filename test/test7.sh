# 建立集群
# ./server -node 1.1.1.1
# 添加节点
# ./server -node 1.1.1.2 -cluster 1.1.1.1
# 查看值归属的节点
./client -h 1.1.1.1 -c set -k keya -v a

./client -h 1.1.1.1 -c set -k keyb -v b

./client -h 1.1.1.1 -c set -k keyc -v c

./client -h 1.1.1.1 -c set -k keyd -v d

./client -h 1.1.1.1 -c set -k keye -v e
# 添加节点
# ./server -node 1.1.1.3 -cluster 1.1.1.2
# 查看值归属的节点
./client -h 1.1.1.1 -c set -k keya -v a

./client -h 1.1.1.1 -c set -k keyb -v b

./client -h 1.1.1.1 -c set -k keyc -v c

./client -h 1.1.1.1 -c set -k keyd -v d

./client -h 1.1.1.1 -c set -k keye -v e
# 停止
# stop 1.1.1.1
# 查看值归属的节点，其会自动调整负载
./client -h 1.1.1.2 -c set -k keya -v a

./client -h 1.1.1.2 -c set -k keyb -v b

./client -h 1.1.1.2 -c set -k keyc -v c

./client -h 1.1.1.2 -c set -k keyd -v d

./client -h 1.1.1.2 -c set -k keye -v e
