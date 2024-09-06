package dao

import (
	"fmt"
	"log"

	"github.com/go-redis/redis"

	lane "laneIM/src/common"
)

func UserNew(rdb *redis.ClusterClient, userid lane.Int64) error {
	err := rdb.SAdd("user:mgr", userid.String()).Err()
	if err != nil {
		log.Fatalf("could not set room info: %v", err)
		return err
	}

	return nil
}

func UserDel(rdb *redis.ClusterClient, userid lane.Int64) (int64, error) {
	err := rdb.SRem("user:mgr", userid.String()).Err()
	if err != nil {
		log.Fatalf("could not set room info: %v", err)
		return 0, err
	}
	rt, err := rdb.Del(fmt.Sprintf("user:online:%s", userid.String())).Result()
	if err != nil {
		return rt, err
	}
	rt, err = rdb.Del(fmt.Sprintf("user:comet:%s", userid.String())).Result()
	if err != nil {
		return rt, err
	}
	rt, err = rdb.Del(fmt.Sprintf("user:room:%s", userid.String())).Result()
	if err != nil {
		return rt, err
	}
	return rt, nil
}

func UserOnline(rdb *redis.ClusterClient, userid lane.Int64, comet string) error {
	err := rdb.Set(fmt.Sprintf("user:online:%s", userid.String()), 1, 0).Err()
	if err != nil {
		return err
	}
	err = rdb.Set(fmt.Sprintf("user:comet:%s", userid.String()), comet, 0).Err()
	return err
}
func UserOffline(rdb *redis.ClusterClient, userid lane.Int64, comet string) error {
	err := rdb.Set(fmt.Sprintf("user:online:%s", userid.String()), 0, 0).Err()
	return err
}

func UserQueryOnline(rdb *redis.ClusterClient, userid lane.Int64) (bool, string, error) {
	rt, err := rdb.Get(fmt.Sprintf("user:online:%s", userid.String())).Int64()
	if err != nil {
		return false, "", err
	}
	if rt != 1 {
		return false, "", nil
	}
	cometAddr, err := rdb.Get(fmt.Sprintf("user:comet:%s", userid.String())).Result()
	return true, cometAddr, err
}

func UserJoinRoomid(rdb *redis.ClusterClient, userid lane.Int64, roomid lane.Int64) error {
	_, err := rdb.SAdd(fmt.Sprintf("user:room:%s", userid.String()), roomid.String()).Result()
	return err
}

func UserQuitRoomid(rdb *redis.ClusterClient, userid lane.Int64, roomid lane.Int64) error {
	_, err := rdb.SRem(fmt.Sprintf("user:room:%s", userid.String()), roomid.String()).Result()
	return err
}

func UserQueryRoomid(rdb *redis.ClusterClient, userid lane.Int64) ([]int64, error) {
	roomStr, err := rdb.SMembers(fmt.Sprintf("user:room:%s", userid.String())).Result()
	if err != nil {
		return nil, err
	}
	return RedisStrsToInt64(roomStr)
}
