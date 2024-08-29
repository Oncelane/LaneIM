package model

import (
	"fmt"
	"log"

	pb "laneIM/proto/msg"

	"github.com/go-redis/redis"
	"google.golang.org/protobuf/proto"
)

func RoomNew(rdb *redis.ClusterClient, roomid int64, userid int64, serverAddr string) error {
	err := rdb.SAdd("roomMgr", roomid, 0).Err()
	if err != nil {
		log.Fatalf("could not set room info: %v", err)
		return err
	}
	err = rdb.Set(fmt.Sprintf("room:online:%d", roomid), 1, 0).Err()
	if err != nil {
		log.Fatalf("could not set room info: %v", err)
		return err
	}
	err = rdb.SAdd(fmt.Sprintf("room:comet:%d", roomid), serverAddr, 0).Err()
	if err != nil {
		log.Fatalf("could not set room info: %v", err)
		return err
	}
	err = rdb.SAdd(fmt.Sprintf("room:userid:%d", roomid), userid, 0).Err()
	if err != nil {
		log.Fatalf("could not set room info: %v", err)
		return err
	}
	return nil
}

func RoomDel(rdb *redis.ClusterClient, roomid int64) (int64, error) {

	err := rdb.SRem("roomMgr", roomid).Err()
	if err != nil {
		log.Fatalf("could not set room info: %v", err)
		return 0, err
	}
	rt, err := rdb.Del(fmt.Sprintf("room:online:%d", roomid), fmt.Sprintf("room:comet:%d", roomid), fmt.Sprintf("room:userid:%d", roomid)).Result()
	return rt, err
}

func RoomJoinUser(rdb *redis.ClusterClient, roomid int64, userid int64, comet string) error {
	_, err := rdb.SAdd(fmt.Sprintf("room:comet:%d", roomid), comet).Result()
	if err != nil {
		return err
	}
	_, err = rdb.SAdd(fmt.Sprintf("room:userid:%d", roomid), userid).Result()
	if err != nil {
		return err
	}

	return err
}

func RoomQuit(rdb *redis.ClusterClient, roomid int64, userid int64, comet string) error {
	_, err := rdb.SRem(fmt.Sprintf("room:comet:%d", roomid), comet).Result()
	if err != nil {
		return err
	}
	_, err = rdb.SRem(fmt.Sprintf("room:userid:%d", roomid), userid).Result()
	if err != nil {
		return err
	}
	return nil
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

func UserDel(rdb *redis.ClusterClient, r *pb.UserInfo) (int64, error) {
	// 获取房间信息
	num, err := rdb.Del(fmt.Sprintf("user:%d", r.Userid)).Result()
	if err != nil {
		if err == redis.Nil {
			err = nil
			return 0, err
		}
		return 0, err
	}
	return num, nil
}

func AtomicRedis(rdb *redis.ClusterClient, key string, value string) error {
	// Lua脚本，用于实现原子性的读取-修改-写入
	// KEYS[1] 是键名，ARGV[1] 是新值
	script := `  
    local currentValue = redis.call("GET", KEYS[1])  
    if currentValue == false then  
        -- 如果键不存在，则直接设置新值  
        redis.call("SET", KEYS[1], ARGV[1])  
    else  
        -- 如果键存在，则进行某种修改（这里假设只是打印出来，实际可以修改后再设置）  
        -- 注意：这里为了示例简单，并没有真正修改值，只是打印  
        print("Current value:", currentValue)  
        -- 真实场景中，你可能需要基于currentValue计算新值，然后设置  
        -- redis.call("SET", KEYS[1], newValue)  
    end  
    return currentValue  
    `

	// 执行Lua脚本
	result, err := rdb.Eval(script, []string{key}, new).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("Result:", result)
	return nil
}
