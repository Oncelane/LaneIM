# laneIM golang 分布式 im

安装 protobuf grpc complier
sudo apt-get install protobuf-compiler

gRPC Go 插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

protoc 命令
protoc --go_out=.. --go-grpc_out=.. --go-grpc_opt=require_unimplemented_servers=false -I. -Iproto proto/msg/msg.proto proto/comet/comet.proto proto/logic/logic.proto

# 启动 redis 集群

redis-server ./redis.conf

redis-cli --cluster create 127.0.0.1:7001 127.0.0.1:7002 127.0.0.1:7003 127.0.0.1:7004 127.0.0.1:7005 127.0.0.1:7006 --cluster-replicas 1

​ 执行 redis-cli -p 7001 进入客户端并通过 `info replication` 查看集群信息

通过`cluster nodes`查看集群关系

通过`cluster info`查看集群信息
