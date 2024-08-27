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
	conf  config.Logic
	etcd  *pkg.EtcdClient
	redis *pkg.RedisClient

	comets map[string]*ClientComet
}

// new and register
func NewServiceLogic(conf config.Logic) *ServiceLogic {
	s := &ServiceLogic{
		etcd: pkg.NewEtcd(),
		conf: conf,
	}

	// register etcd
	s.etcd.SetAddr("grpc:logic/"+s.conf.Name, s.conf.Addr)

	// get redis address
	redisAddrs := s.etcd.GetAddr("redis")
	log.Println("获取到的redis地址：", redisAddrs)
	redis := pkg.NewRedisClient(redisAddrs)
	s.redis = redis

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

func (s *ServiceLogic) SendMsg(_ context.Context, in *pb.SendMsgReq) (*pb.NoResp, error) {
	switch in.Path {
	case "brodcast/test":
		// 从redis查询房间内其余成员
		room, err := model.RoomGet(s.redis.Client, in.Roomid)
		if err != nil {
			log.Panicln("get room member faild:", err)
		}
		log.Panicln("get roomInfo", room)
		for userid, online := range room.Users {
			if !online {
				log.Println("userid:", userid, "not online")
				continue
			}
			user, err := model.UserGet(s.redis.Client, userid)
			if err != nil {
				log.Panicln("ger user faild:", err)
			}
			//拿到user的comet地址,发起grpc调用

		}

		// 获取comet地址
	}
	return nil, nil
}
func (s *ServiceLogic) NewUser(_ context.Context, in *pb.NewUserReq) (*pb.NoResp, error) {
	return nil, nil
}
func (s *ServiceLogic) DelUser(context.Context, *pb.DelUserReq) (*pb.NoResp, error) {
	return nil, nil
}
func (s *ServiceLogic) SetOnline(_ context.Context, in *pb.SetOnlineReq) (*pb.NoResp, error) {

	return nil, nil
}
func (s *ServiceLogic) SetOffline(_ context.Context, in *pb.SetOfflineReq) (*pb.NoResp, error) {

	return nil, nil
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
