package logic

import (
	"context"
	"fmt"
	pb "laneIM/proto/logic"
	"laneIM/src/config"
	"laneIM/src/pkg"
	"log"
	"net"

	"github.com/IBM/sarama"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type Logic struct {
	conf  config.Logic
	etcd  *pkg.EtcdClient
	redis *pkg.RedisClient
	kafka *pkg.KafkaProducer
	grpc  *grpc.Server
}

// new and register
func NewLogic(conf config.Logic) *Logic {
	s := &Logic{
		etcd: pkg.NewEtcd(conf.Etcd),
		conf: conf,
	}

	// init redis
	redisAddrs := s.etcd.GetAddr("redis")
	log.Println("获取到的redis地址：", redisAddrs)
	redis := pkg.NewRedisClient(redisAddrs)
	s.redis = redis

	// init kafka producer
	s.kafka = pkg.NewKafkaProducer(conf.KafkaProducer)

	// register etcd
	s.etcd.SetAddr("grpc:logic/"+s.conf.Name, s.conf.Addr)

	// server grpc
	lis, err := net.Listen("tcp", conf.Addr)
	if err != nil {
		log.Fatalf("error: logic start faild")
	}
	gServer := grpc.NewServer()
	pb.RegisterLogicServer(gServer, s)
	fmt.Println("Logic serivce is running on port")
	s.grpc = gServer
	go func() {
		if err := gServer.Serve(lis); err != nil {
			log.Fatalln("failed to serve : ", err.Error())
		}
	}()

	return s
}

func (l *Logic) Close() {
	log.Println("logic exit:", l.conf.Addr)
	l.etcd.DelAddr("grpc:logic/"+l.conf.Name, l.conf.Addr)
	l.grpc.Stop()
}

var _ pb.LogicServer = new(Logic)

func (s *Logic) SendMsg(_ context.Context, in *pb.SendMsgReq) (*pb.NoResp, error) {
	switch in.Path {
	case "brodcast/test":

		//no op just send to kafka
		data, err := proto.Marshal(in)
		if err != nil {
			log.Println("proto marshal error")
		}

		msg := &sarama.ProducerMessage{
			Topic: "laneIM",
			Value: sarama.ByteEncoder(data),
		}
		_, _, err = s.kafka.Client.SendMessage(msg)
		if err != nil {
			log.Println("faild to send kafka:", err)
		}
		log.Println("success send message:", in.String())

	}
	return nil, nil
}
func (s *Logic) NewUser(_ context.Context, in *pb.NewUserReq) (*pb.NoResp, error) {
	return nil, nil
}
func (s *Logic) DelUser(context.Context, *pb.DelUserReq) (*pb.NoResp, error) {
	return nil, nil
}
func (s *Logic) SetOnline(_ context.Context, in *pb.SetOnlineReq) (*pb.NoResp, error) {

	return nil, nil
}
func (s *Logic) SetOffline(_ context.Context, in *pb.SetOfflineReq) (*pb.NoResp, error) {

	return nil, nil
}
func (s *Logic) JoinRoom(_ context.Context, in *pb.JoinRoomReq) (*pb.NoResp, error) {
	return nil, nil
}
func (s *Logic) QuitRoom(context.Context, *pb.QuitRoomReq) (*pb.NoResp, error) {
	return nil, nil
}
func (s *Logic) QueryRoom(context.Context, *pb.QueryRoomReq) (*pb.QueryRoomResp, error) {
	return nil, nil
}
func (s *Logic) QueryServer(context.Context, *pb.QueryServerReq) (*pb.QueryServerResp, error) {
	return nil, nil
}
