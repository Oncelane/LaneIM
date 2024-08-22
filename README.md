# laneIM golang 分布式 im

安装 protobuf grpc complier
sudo apt-get install protobuf-compiler

gRPC Go 插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

protoc 命令
protoc --go_out=. --go-grpc_out=. -I. -Iproto proto/msg/msg.proto proto/comet/comet.proto proto/logic/logic.proto
