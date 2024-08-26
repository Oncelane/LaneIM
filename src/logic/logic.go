package logic

import (
	"context"
	"fmt"
	pb "laneIM/proto/logic"
	"laneIM/src/config"
	"laneIM/src/model"
	"laneIM/src/pkg"
	"log"
	"net"

	"google.golang.org/grpc"
)

type ServiceLogic struct {
	etcd pkg.EtcdClient
	conf config.Logic
}

func NewServiceLogic(conf config.Logic) *ServiceLogic {
	s := &ServiceLogic{
		etcd: *pkg.NewEtcd(),
		conf: conf,
	}
	s.etcd.SetAddr("grpc:logic/"+conf.Name, conf.Addr)
	lis, err := net.Listen("tcp", conf.Addr)
	if err != nil {
		log.Fatalf("error: logic start faild")
	}
	gServer := grpc.NewServer()
	pb.RegisterLogicServer(gServer, s)
	fmt.Println("Logic serivce is running on port")
	if err := gServer.Serve(lis); err != nil {
		log.Fatalln("failed to serve : ", err.Error())
	}
	return s
}

var _ pb.LogicServer = new(ServiceLogic)

func (s *ServiceLogic) SendMsg(context.Context, *pb.SendMsgReq) (*pb.NoResp, error) {

	return nil, nil
}
func (s *ServiceLogic) NewUser(_ context.Context, in *pb.NewUserReq) (*pb.NoResp, error) {

	err := s.etcd.NewUser(model.UserStage{
		Userid:  in.Userinfo.Userid,
		Machine: in.Userinfo.Machine,
	})
	return nil, err
}
func (s *ServiceLogic) DelUser(context.Context, *pb.DelUserReq) (*pb.NoResp, error) {
	return nil, nil
}
func (s *ServiceLogic) SetOnline(_ context.Context, in *pb.SetOnlineReq) (*pb.NoResp, error) {
	err := s.etcd.SetUserOnline(in.Userinfo.Userid)
	return nil, err
}
func (s *ServiceLogic) SetOffline(_ context.Context, in *pb.SetOfflineReq) (*pb.NoResp, error) {
	err := s.etcd.SetUserOffline(in.Userinfo.Userid)
	return nil, err
}
func (s *ServiceLogic) JoinRoom(_ context.Context, in *pb.JoinRoomReq) (*pb.NoResp, error) {
	return nil, nil
}
func (s *ServiceLogic) QuitRoom(context.Context, *pb.QuitRoomReq) (*pb.NoResp, error) {
	return nil, nil
}
func (s *ServiceLogic) QueryRoom(context.Context, *pb.QueryRoomReq) (*pb.QueryRoomResp, error) {
	return nil, nil
}
func (s *ServiceLogic) QueryServer(context.Context, *pb.QueryServerReq) (*pb.QueryServerResp, error) {
	return nil, nil
}
