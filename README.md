# laneIM golang 分布式 im

安装 protobuf grpc complier
sudo apt-get install protobuf-compiler

gRPC Go插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

protoc 命令
protoc --go-grpc_out=. --go_out=. service.proto