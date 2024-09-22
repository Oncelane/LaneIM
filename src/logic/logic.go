package logic

import (
	"context"
	"fmt"
	pb "laneIM/proto/logic"
	"laneIM/proto/msg"
	"laneIM/src/config"
	"laneIM/src/dao"
	"laneIM/src/dao/cql"
	"laneIM/src/dao/localCache"
	"laneIM/src/dao/sql"
	"laneIM/src/model"
	"laneIM/src/pkg"
	"laneIM/src/pkg/laneLog"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/allegro/bigcache"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type Logic struct {
	conf     config.Logic
	etcd     *pkg.EtcdClient
	cache    *bigcache.BigCache
	redis    *pkg.RedisClient
	db       *sql.SqlDB
	scyllaDB *cql.ScyllaDB
	kafka    *pkg.KafkaProducer
	grpc     *grpc.Server
	uuid     *pkg.UuidGenerator
	daoo     *dao.Dao

	//only to sync comet infomation to mysql
	cometMu sync.RWMutex
	comets  map[string]struct{}
}

// new and register
func NewLogic(conf config.Logic) *Logic {
	s := &Logic{
		etcd:   pkg.NewEtcd(conf.Etcd),
		conf:   conf,
		daoo:   dao.NewDao(conf.Mysql.BatchWriter),
		comets: make(map[string]struct{}),
		cache:  localCache.Cache(time.Minute),
	}
	s.uuid = pkg.NewUuidGenerator(int64(conf.Id))

	s.db = sql.NewDB(conf.Mysql)
	s.scyllaDB = cql.NewCqlDB(conf.ScyllaDB)
	dao.Init(s.db.DB, s.scyllaDB)

	// init redis
	redisAddrs := s.etcd.GetAddr("redis")
	// laneLog.Logger.Infoln("获取到的redis地址：", redisAddrs)
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
	go s.WatchComet()
	s.etcd.SetAddr("grpc:logic:"+s.conf.Name, s.conf.Addr)
	return s
}

func (l *Logic) Close() {
	laneLog.Logger.Infoln("logic exit:", l.conf.Addr)
	l.etcd.DelAddr("grpc:logic:"+l.conf.Name, l.conf.Addr)
	l.grpc.Stop()
}

func (l *Logic) WatchComet() {
	for {
		addrs := l.etcd.GetAddr("grpc:comet")
		remoteAddrs := make(map[string]struct{})
		for _, addr := range addrs {
			remoteAddrs[addr] = struct{}{}
			if _, exist := l.comets[addr]; exist {
				continue
			}

			l.comets[addr] = struct{}{}
			// discovery comet
			err := l.db.AddComet(addr)
			if err != nil {
				continue
			}
			laneLog.Logger.Infoln("[laneEtcd] discovery comet:", addr)
		}
		for addr := range l.comets {
			if _, exist := remoteAddrs[addr]; !exist {
				l.cometMu.Lock()
				delete(l.comets, addr)
				l.cometMu.Unlock()
				//discovery missing comet
				err := l.db.DelComet(addr)
				if err != nil {
					continue
				}
				laneLog.Logger.Infoln("[server] remove comet:", addr)
			}
		}

		time.Sleep(time.Second)
	}
}

var _ pb.LogicServer = new(Logic)

func (s *Logic) SendMsg(_ context.Context, in *msg.SendMsgReq) (*pb.NoResp, error) {
	switch in.Path {
	case "sendRoom":

		//no op just send to kafka
		data, err := proto.Marshal(in)
		if err != nil {
			laneLog.Logger.Infoln("proto marshal error")
		}

		msg := &sarama.ProducerMessage{
			Topic: "laneIM",
			Value: sarama.ByteEncoder(data),
		}
		_, _, err = s.kafka.Client.SendMessage(msg)
		if err != nil {
			laneLog.Logger.Fatalln("[server] faild to send kafka:", err)
		}
		// laneLog.Logger.Infoln("success send message:", in.String())

	}
	return nil, nil
}
func (s *Logic) SendMsgBatch(_ context.Context, in *msg.SendMsgBatchReq) (*pb.NoResp, error) {
	// start := time.Now()
	data, err := proto.Marshal(in)
	if err != nil {
		laneLog.Logger.Infoln("proto marshal error")
	}
	msg := &sarama.ProducerMessage{
		Topic: "laneIM",
		Value: sarama.ByteEncoder(data),
	}
	s.scyllaDB.AddChatMessageBatch(in)
	_, _, err = s.kafka.Client.SendMessage(msg)
	// laneLog.Logger.Infof("[kafka] [sendMsg] %d bytes", len(data))
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to send kafka:", err)
	}
	// laneLog.Logger.Debugf("send message count %d spand time :%v", len(in.Msgs), time.Since(start))
	return nil, nil
}

func (s *Logic) QueryStoreMsgBatch(_ context.Context, in *msg.QueryMultiRoomPagesReq) (out *msg.QueryMultiRoomPagesReply, err error) {
	out = new(msg.QueryMultiRoomPagesReply)
	out.RoomMultiPageMsgs = make([]*msg.QueryMultiRoomPagesReply_RoomMultiPageMsg, len(in.RoomMultiPageInfos))
	for i, q := range in.RoomMultiPageInfos { //each room
		out.RoomMultiPageMsgs[i] = &msg.QueryMultiRoomPagesReply_RoomMultiPageMsg{
			PagesMsgs: make([]*msg.QueryMultiRoomPagesReply_RoomMultiPageMsg_PageMsgs, len(q.PageInfos)),
		}
		for pageIndex, page := range q.PageInfos {
			modelData, valid, err := s.scyllaDB.PageChatMessageByMessageid(q.Roomid, page.MessageId, int(page.Size))
			if err != nil {
				laneLog.Logger.Errorln("faild to cql get q.Roomid ", q.Roomid, "page.MessageId", page.MessageId, "page.Size", page.Size, "err =", err)
			}
			if valid != 0 {
				out.RoomMultiPageMsgs[i].PagesMsgs[pageIndex] = &msg.QueryMultiRoomPagesReply_RoomMultiPageMsg_PageMsgs{
					Msgs: TrueMsgFromModelToProto(modelData, valid),
				}
			}
		}
	}

	return out, nil
}

func TrueMsgFromModelToProto(in []model.ChatMessage, valid int) []*msg.RoomStoreChatMessage {
	out := make([]*msg.RoomStoreChatMessage, valid)
	for i := 0; i < valid; i++ {
		out[i] = &msg.RoomStoreChatMessage{
			Messageid: in[i].MessageID,
			Userid:    in[i].UserID,
			Userseq:   in[i].UserSeq,
			Timeunix:  in[i].Timestamp.Unix(),
			Data:      in[i].Content,
		}
	}
	return out
}

func (s *Logic) NewUser(_ context.Context, in *pb.NewUserReq) (*pb.NewUserResp, error) {
	uuid := s.uuid.Generator()
	// laneLog.Logger.Infoln("new user id:", uuid)
	err := s.db.NewUser(uuid)
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to new user", err)
		return nil, err
	}
	resp := &pb.NewUserResp{
		Userid: uuid,
	}
	return resp, nil
}

func (s *Logic) NewUserBatch(_ context.Context, in *pb.NewUserBatchReq) (*pb.NewUserBatchResp, error) {
	if in.Count <= 0 {
		return nil, fmt.Errorf("wrong count")
	}
	uuids := make([]int64, in.Count)
	for i := range in.Count {
		uuids[i] = s.uuid.Generator()
	}

	// laneLog.Logger.Infoln("new user id:", uuid)
	err := s.db.NewUserBatch(nil, uuids)
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to new user", err)
		return nil, err
	}
	resp := &pb.NewUserBatchResp{
		Userid: uuids,
	}
	return resp, nil
}

func (s *Logic) NewRoom(_ context.Context, in *pb.NewRoomReq) (*pb.NewRoomResp, error) {
	uuid := s.uuid.Generator()
	// laneLog.Logger.Infoln("new user id:", uuid)
	err := s.db.NewRoom(uuid, in.Userid, in.CometAddr)
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to new user", err)
		return nil, err
	}
	resp := &pb.NewRoomResp{
		Roomid: uuid,
	}
	return resp, nil
}

func (s *Logic) DelUser(_ context.Context, in *pb.DelUserReq) (*pb.NoResp, error) {
	err := s.db.DelUser(in.Userid)
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to del user:", s.uuid, err)
	}
	return nil, nil
}

func (s *Logic) SetOnline(_ context.Context, in *pb.SetOnlineReq) (*pb.NoResp, error) {
	err := s.db.SetUserOnline(in.Userid, in.Server)
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to new user", err)
		return nil, err
	}
	err = s.db.AddRoomCometWithUseridBatch(s.redis.Client, []int64{in.Userid}, in.Server)
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to new user", err)
		return nil, err
	}
	return nil, nil
}

func (s *Logic) SetOnlineBatch(_ context.Context, in *pb.SetOnlineBatchReq) (*pb.NoResp, error) {
	// start := time.Now()
	tx := s.db.DB.Begin()
	err := s.db.SetUserOnlineBatch(tx, in.Userid, in.Server)
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to new user batch", err)
		return nil, err
	}
	err = s.db.AddRoomCometWithUseridBatch(s.redis.Client, in.Userid, in.Server)
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to new user batch", err)
		return nil, err
	}
	err = tx.Commit().Error
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to commit set online batch", err)
		return nil, err
	}
	// laneLog.Logger.Debugln("SetOnlineBatch spand", time.Since(start))
	return nil, nil
}

func (s *Logic) SetOfflineBatch(_ context.Context, in *pb.SetOfflineBatchReq) (*pb.NoResp, error) {
	// start := time.Now()
	tx := s.db.DB.Begin()
	err := s.db.SetUserOfflineBatch(tx, in.Userid)
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to new user batch", err)
		return nil, err
	}
	err = tx.Commit().Error
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to commit set online batch", err)
		return nil, err
	}
	// laneLog.Logger.Debugln("SetOfflineBatch spand", time.Since(start))
	return nil, nil
}

func (s *Logic) SetOffline(_ context.Context, in *pb.SetOfflineReq) (*pb.NoResp, error) {
	err := s.db.SetUseroffline(in.Userid)
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to new user", err)
	}
	return nil, nil
}

func (s *Logic) JoinRoom(_ context.Context, in *pb.JoinRoomReq) (*pb.NoResp, error) {
	start := time.Now()
	err := s.db.AddRoomUser(in.Roomid, in.Userid)
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild room to add user:", err)
		return nil, err
	}

	laneLog.Logger.Debugln("join room normal spand", time.Since(start))
	return nil, err
}

func (s *Logic) JoinRoomBatch(_ context.Context, in *pb.JoinRoomBatchReq) (*pb.NoResp, error) {
	start := time.Now()
	// laneLog.Logger.Debugln("Userid ", in.Userid, "join Roomid", in.Roomid)
	// tx := s.db.DB.Begin()
	err := s.db.AddRoomUserBatch(nil, in.Roomid, in.Userid)
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild room to add user:", err)
		return nil, err
	}
	err = s.db.AddRoomCometWithUseridBatch(s.redis.Client, in.Userid, in.Comet)
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to new user batch", err)
		return nil, err
	}
	laneLog.Logger.Debugln("[dao] [JoinRoom] spand time:", time.Since(start))
	return nil, nil
}

func (s *Logic) QuitRoom(_ context.Context, in *pb.QuitRoomReq) (*pb.NoResp, error) {
	err := s.db.DelUserRoom(in.Userid, in.Roomid)
	if err != nil {
		laneLog.Logger.Fatalln("faild user %d to quit room:%d\n", in.Roomid, in.Userid)
	}
	return nil, err
}

func (s *Logic) QueryRoom(_ context.Context, in *pb.QueryRoomReq) (*pb.QueryRoomResp, error) {
	start := time.Now()
	out := pb.QueryRoomResp{
		Roomids: make([]*pb.QueryRoomResp_RoomSlice, len(in.Userid)),
	}
	roomidss, err := s.daoo.UserRoomBatch(s.cache, s.redis.Client, s.db, in.Userid)
	// laneLog.Logger.Infof("logic dao query user%d have room%v\n", userid, roomids)
	if err != nil {
		return nil, err
	}
	for i := range roomidss {
		out.Roomids[i] = &pb.QueryRoomResp_RoomSlice{
			Roomid: roomidss[i],
		}
	}
	laneLog.Logger.Debugln("[dao] [QueryRoom] spand time:", time.Since(start))
	return &out, nil
}

func (s *Logic) QueryServer(context.Context, *pb.QueryServerReq) (*pb.QueryServerResp, error) {
	return nil, nil
}

func (s *Logic) Auth(_ context.Context, in *pb.AuthReq) (*pb.AuthResp, error) {
	// laneLog.Logger.Infoln("auth pass:", in.Userid)
	out := &pb.AuthResp{Pass: true}
	return out, nil
}

func (s *Logic) QueryLast(_ context.Context, in *pb.QueryLastReq) (*pb.QueryLastRelpy, error) {
	start := time.Now()
	wait := sync.WaitGroup{}
	wait.Add(len(in.Roomid))
	out := new(pb.QueryLastRelpy)
	out.MessageId = make([]int64, len(in.Roomid))
	for i, roomid := range in.Roomid {
		go func(index int, roomid int64) {
			defer wait.Done()
			rt, err, _ := s.scyllaDB.SingleFlightGroup.Do(strconv.FormatInt(roomid, 36), func() (interface{}, error) {
				return s.scyllaDB.QueryLatestGroupMessageid(roomid)
			})
			if err != nil {
				laneLog.Logger.Fatalln("[dao] wrong querylastMessageid", err)
				return
			}
			msgId, ok := rt.(int64)
			if !ok {
				laneLog.Logger.Fatalln("[dao] wrong querylastMessageid", err)
				return
			}
			out.MessageId[i] = msgId

		}(i, roomid)

	}
	wait.Wait()
	laneLog.Logger.Debugln("[dao] [QueryLastMsgId] spand time:", time.Since(start))
	return out, nil
}
