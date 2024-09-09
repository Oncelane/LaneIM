package logic

import (
	"context"
	"fmt"
	pb "laneIM/proto/logic"
	"laneIM/src/config"
	"laneIM/src/dao"
	"laneIM/src/dao/sql"
	"laneIM/src/pkg"
	"log"
	"net"

	"github.com/IBM/sarama"
	"github.com/bwmarrin/snowflake"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

type UuidGenerator struct {
	node *snowflake.Node
}

func NewUuidGenerator(name int64) *UuidGenerator {
	node, err := snowflake.NewNode(name)
	if err != nil {
		log.Panicln("faild to create snowflake node")
	}
	return &UuidGenerator{
		node: node,
	}
}

func (u *UuidGenerator) Generator() (rt int64) {
	return u.node.Generate().Int64()
}

type Logic struct {
	conf  config.Logic
	etcd  *pkg.EtcdClient
	redis *pkg.RedisClient
	db    *gorm.DB
	kafka *pkg.KafkaProducer
	grpc  *grpc.Server
	uuid  *UuidGenerator
}

// new and register
func NewLogic(conf config.Logic) *Logic {
	s := &Logic{
		etcd: pkg.NewEtcd(conf.Etcd),
		conf: conf,
	}
	s.uuid = NewUuidGenerator(int64(conf.Id))

	s.db = sql.DB(conf.Mysql)

	// init redis
	redisAddrs := s.etcd.GetAddr("redis")
	log.Println("获取到的redis地址：", redisAddrs)
	redis := pkg.NewRedisClient(config.Redis{Addr: redisAddrs})
	s.redis = redis

	// init kafka producer
	s.kafka = pkg.NewKafkaProducer(conf.KafkaProducer)

	// server grpc
	lis, err := net.Listen("tcp", conf.Addr)
	if err != nil {
		log.Fatalln("error: logic start faild", err)
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
	// register etcd
	s.etcd.SetAddr("grpc:logic:"+s.conf.Name, s.conf.Addr)
	return s
}

func (l *Logic) Close() {
	log.Println("logic exit:", l.conf.Addr)
	l.etcd.DelAddr("grpc:logic:"+l.conf.Name, l.conf.Addr)
	l.grpc.Stop()
}

var _ pb.LogicServer = new(Logic)

func (s *Logic) SendMsg(_ context.Context, in *pb.SendMsgReq) (*pb.NoResp, error) {
	switch in.Path {
	case "sendRoom":

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

func (s *Logic) NewUser(_ context.Context, in *pb.NewUserReq) (*pb.NewUserResp, error) {
	uuid := s.uuid.Generator()
	log.Println("new user id:", uuid)
	err := sql.NewUser(s.db, uuid)
	if err != nil {
		log.Println("faild to new user", err)
		return nil, err
	}
	resp := &pb.NewUserResp{
		Userid: uuid,
	}
	return resp, nil
}

func (s *Logic) DelUser(_ context.Context, in *pb.DelUserReq) (*pb.NoResp, error) {
	err := sql.DelUser(s.db, in.Userid)
	if err != nil {
		log.Println("faild to del user:", s.uuid, err)
	}
	return nil, nil
}

func (s *Logic) SetOnline(_ context.Context, in *pb.SetOnlineReq) (*pb.NoResp, error) {
	err := sql.SetUserOnline(s.db, in.Userid, in.Server)
	if err != nil {
		log.Println("faild to new user", err)
		return nil, err
	}

	err = sql.AddRoomCometWithUserid(s.db, in.Userid, in.Server)
	if err != nil {
		log.Println("faild to new user", err)
		return nil, err
	}
	return nil, nil
}

func (s *Logic) SetOffline(_ context.Context, in *pb.SetOfflineReq) (*pb.NoResp, error) {
	err := sql.SetUseroffline(s.db, in.Userid)
	if err != nil {
		log.Println("faild to new user", err)
	}
	return nil, nil
}

func (s *Logic) JoinRoom(_ context.Context, in *pb.JoinRoomReq) (*pb.NoResp, error) {
	err := sql.AddUserRoom(s.db, in.Userid, in.Roomid)
	if err != nil {
		log.Printf("faild user %d to join room:%d\n", in.Userid, in.Roomid)
		return nil, err
	}
	err = sql.AddRoomUser(s.db, in.Roomid, in.Roomid)
	if err != nil {
		log.Println("faild room to add user:", err)
		return nil, err
	}
	return nil, err
}

func (s *Logic) QuitRoom(_ context.Context, in *pb.QuitRoomReq) (*pb.NoResp, error) {
	err := sql.DelUserRoom(s.db, in.Userid, in.Roomid)
	if err != nil {
		log.Printf("faild user %d to quit room:%d\n", in.Roomid, in.Userid)
	}
	return nil, err
}

func (s *Logic) QueryRoom(_ context.Context, in *pb.QueryRoomReq) (*pb.QueryRoomResp, error) {
	out := pb.QueryRoomResp{
		Roomids: make([]*pb.QueryRoomResp_RoomSlice, 0),
	}
	for _, userid := range in.Userid {
		roomids, err := dao.UserRoom(s.redis.Client, s.db, userid)
		log.Printf("logic dao query user%d have room%v\n", userid, roomids)
		if err != nil {
			return nil, err
		}
		out.Roomids = append(out.Roomids, &pb.QueryRoomResp_RoomSlice{
			Roomid: roomids,
		})

	}

	return &out, nil
}

func (s *Logic) QueryServer(context.Context, *pb.QueryServerReq) (*pb.QueryServerResp, error) {
	return nil, nil
}

func (s *Logic) Auth(_ context.Context, in *pb.AuthReq) (*pb.AuthResp, error) {
	log.Println("auth pass:", in.Userid)
	out := &pb.AuthResp{Pass: true}
	return out, nil
}
