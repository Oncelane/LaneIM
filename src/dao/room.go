package dao

import (
	"fmt"
	"log"

	"laneIM/proto/msg"
	lane "laneIM/src/common"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

func AllRoomid(rdb *redis.ClusterClient, db *gorm.DB) ([]int64, error) {
	rt, err := RdsAllRoomid(rdb)
	if err != nil {
		return rt, err
	}
	if len(rt) != 0 {
		return rt, nil
	}
	log.Println("触发sql查询")
	rt, err = SqlAllRoomid(db)
	if err != nil {
		return rt, err
	}
	log.Println("同步到redis")
	rtt := make([]lane.Int64, len(rt))
	for i := range rt {
		rtt[i] = lane.Int64(rt[i])
	}
	RdsSetAllRoomid(rdb, rtt)
	return rt, nil
}

func RdsAllRoomid(rdb *redis.ClusterClient) ([]int64, error) {
	roomStr, err := rdb.SMembers("room:mgr").Result()
	if err != nil {
		return nil, err
	}
	return RedisStrsToInt64(roomStr)
}

func RdsSetAllRoomid(rdb *redis.ClusterClient, roomids []lane.Int64) error {
	pipe := rdb.Pipeline()
	for _, member := range roomids {
		pipe.SAdd("room:mgr", member.String())
	}
	_, err := pipe.Exec()
	return err
}

func RoomNew(rdb *redis.ClusterClient, roomid lane.Int64, userid lane.Int64, serverAddr string) error {
	err := rdb.SAdd("room:mgr", roomid.String()).Err()
	if err != nil {
		log.Fatalf("could not set room info: %v", err)
		return err
	}
	err = rdb.Set(fmt.Sprintf("room:online:%s", roomid.String()), 1, 0).Err()
	if err != nil {
		log.Fatalf("could not set room info: %v", err)
		return err
	}
	err = rdb.SAdd(fmt.Sprintf("room:comet:%s", roomid.String()), serverAddr).Err()
	if err != nil {
		log.Fatalf("could not set room info: %v", err)
		return err
	}
	err = rdb.SAdd(fmt.Sprintf("room:userid:%s", roomid.String()), userid.String()).Err()
	if err != nil {
		log.Fatalf("could not set room info: %v", err)
		return err
	}
	return nil
}

func RoomDel(rdb *redis.ClusterClient, roomid lane.Int64) (lane.Int64, error) {

	err := rdb.SRem("room:mgr", roomid.String()).Err()
	if err != nil {
		log.Fatalf("could not set room info: %v", err)
		return 0, err
	}
	rt, err := rdb.Del(fmt.Sprintf("room:online:%s", roomid.String())).Result()
	if err != nil {
		return lane.Int64(rt), err
	}
	rt, err = rdb.Del(fmt.Sprintf("room:comet:%s", roomid.String())).Result()
	if err != nil {
		return lane.Int64(rt), err
	}
	rt, err = rdb.Del(fmt.Sprintf("room:userid:%s", roomid.String())).Result()
	if err != nil {
		return lane.Int64(rt), err
	}
	return lane.Int64(rt), err
}

func RoomJoinUser(rdb *redis.ClusterClient, roomid lane.Int64, userid lane.Int64) error {
	_, err := rdb.SAdd(fmt.Sprintf("room:userid:%s", roomid.String()), userid.String()).Result()
	return err
}

func RoomQuitUser(rdb *redis.ClusterClient, roomid lane.Int64, userid lane.Int64) error {
	_, err := rdb.SRem(fmt.Sprintf("room:userid:%s", roomid.String()), userid.String()).Result()
	return err
}

func RoomQueryUserid(rdb *redis.ClusterClient, roomid lane.Int64) ([]int64, error) {
	userStr, err := rdb.SMembers(fmt.Sprintf("room:userid:%s", roomid.String())).Result()
	if err != nil {
		return nil, err
	}
	return RedisStrsToInt64(userStr)
}

func RoomPutComet(rdb *redis.ClusterClient, roomid lane.Int64, comet string) error {
	_, err := rdb.SAdd(fmt.Sprintf("room:comet:%s", roomid.String()), comet).Result()
	return err
}

func RoomDelComet(rdb *redis.ClusterClient, roomid lane.Int64, comet string) error {
	_, err := rdb.SRem(fmt.Sprintf("room:comet:%s", roomid.String()), comet).Result()
	return err
}
func RoomQueryComet(rdb *redis.ClusterClient, roomid lane.Int64) ([]string, error) {
	cometStr, err := rdb.SMembers(fmt.Sprintf("room:comet:%s", roomid.String())).Result()
	if err != nil {
		return nil, err
	}
	return cometStr, err
}
func RoomCount(rdb *redis.ClusterClient, roomid lane.Int64) (int64, error) {
	count, err := rdb.SCard(fmt.Sprintf("room:userid:%s", roomid.String())).Result()
	if err != nil {
		return -1, err
	}
	return count, nil
}

func AllUserid(rdb *redis.ClusterClient) ([]int64, error) {
	userStr, err := rdb.SMembers("userMgr").Result()
	if err != nil {
		return nil, err
	}
	return RedisStrsToInt64(userStr)
}

func UpdateRoom(rdb *redis.ClusterClient, roomid int64) (*msg.RoomInfo, error) {
	roomInfo := &msg.RoomInfo{}
	comets, err := RoomQueryComet(rdb, lane.Int64(roomid))
	if err != nil {
		return nil, err
	}

	cometMap := make(map[string]bool)
	for _, addr := range comets {
		cometMap[addr] = true
	}
	userMap := make(map[int64]bool)
	userids, err := RoomQueryUserid(rdb, lane.Int64(roomid))
	if err != nil {
		return nil, err
	}
	for _, userid := range userids {
		userMap[userid] = true
	}
	roomInfo.OnlineNum = -1
	roomInfo.Server = cometMap
	roomInfo.Users = userMap
	return roomInfo, nil
}

// func AtomicRedis(rdb *redis.ClusterClient, key string, value string) error {
// 	// Lua脚本，用于实现原子性的读取-修改-写入
// 	// KEYS[1] 是键名，ARGV[1] 是新值
// 	script := `
//     local currentValue = redis.call("GET", KEYS[1])
//     if currentValue == false then
//         -- 如果键不存在，则直接设置新值
//         redis.call("SET", KEYS[1], ARGV[1])
//     else
//         -- 如果键存在，则进行某种修改（这里假设只是打印出来，实际可以修改后再设置）
//         -- 注意：这里为了示例简单，并没有真正修改值，只是打印
//         print("Current value:", currentValue)
//         -- 真实场景中，你可能需要基于currentValue计算新值，然后设置
//         -- redis.call("SET", KEYS[1], newValue)
//     end
//     return currentValue
//     `

// 	// 执行Lua脚本
// 	result, err := rdb.Eval(script, []string{key}, new).Result()
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("Result:", result)
// 	return nil
// }

func RedisStrsToInt64(strs []string) ([]int64, error) {
	ret := make([]int64, len(strs))
	var tmp lane.Int64
	for i, str := range strs {
		tmp.PasrseString(str)
		ret[i] = int64(tmp)
	}
	return ret, nil
}
