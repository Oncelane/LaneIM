package model

import (
	"fmt"
	"log"

	pb "laneIM/proto/msg"

	"github.com/go-redis/redis"
	"google.golang.org/protobuf/proto"
)

func RoomSet(rdb *redis.ClusterClient, r *pb.RoomInfo) error {
	// 房间信息：包含房间名称和描述
	data, err := proto.Marshal(r)
	if err != nil {
		return err
	}
	err = rdb.Set(fmt.Sprintf("room:%d", r.Roomid), data, 0).Err()
	if err != nil {
		log.Fatalf("could not set room info: %v", err)
	}
	return nil
}

func RoomGet(rdb *redis.ClusterClient, roomid int64) (r *pb.RoomInfo, err error) {
	r = &pb.RoomInfo{}
	// 获取房间信息
	data, err := rdb.Get(fmt.Sprintf("room:%d", roomid)).Result()
	if err != nil {
		if err == redis.Nil {
			err = nil
			return
		}
		return
	}
	err = proto.Unmarshal([]byte(data), r)
	if err != nil {
		return
	}
	return
}

func UserGet(rdb *redis.ClusterClient, userid int64) (r *pb.UserInfo, err error) {
	r = &pb.UserInfo{}
	// 获取房间信息
	data, err := rdb.Get(fmt.Sprintf("user:%d", userid)).Result()
	if err != nil {
		if err == redis.Nil {
			err = nil
			return
		}
		return
	}
	err = proto.Unmarshal([]byte(data), r)
	if err != nil {
		return
	}
	return
}

func UserSet(rdb *redis.ClusterClient, r *pb.UserInfo) error {
	// 房间信息：包含房间名称和描述
	data, err := proto.Marshal(r)
	if err != nil {
		return err
	}
	err = rdb.Set(fmt.Sprintf("user:%d", r.Userid), data, 0).Err()
	if err != nil {
		log.Fatalf("could not set room info: %v", err)
	}
	return nil
}
