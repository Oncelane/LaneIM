package main

import (
	"context"
	"fmt"
	"laneIM/proto/pb"
	"laneIM/src/pkg"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedExampleServiceServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloReq) (*pb.HelloResp, error) {
	return &pb.HelloResp{
		Message: "Hello " + req.Name,
	}, nil
}

var Addr = "127.0.0.1:50051"

type Logic struct {
	etcd pkg.EtcdClient
}

func main() {
	logic := Logic{
		etcd: *pkg.NewEtcd(),
	}
	logic.etcd.SetAddr("logic1", Addr)
	fmt.Println("etcd set logic1:", Addr)
	lis, err := net.Listen("tcp", Addr)
	if err != nil {
		log.Fatalln("failed to listen: ", err.Error())
	}
	s := grpc.NewServer()
	pb.RegisterExampleServiceServer(s, &server{})
	fmt.Println("Logic serivce is running on port")
	if err := s.Serve(lis); err != nil {
		log.Fatalln("failed to serve : ", err.Error())
	}
}
