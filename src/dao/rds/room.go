package rds

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"

	lane "laneIM/src/common"
)

//--------------Room------------

// room:mgr
func AllRoomid(rdb *redis.ClusterClient) ([]int64, error) {
	roomStr, err := rdb.SMembers("room:mgr").Result()
	if len(roomStr) == 0 {
		return nil, redis.Nil
	}
	if err != nil {
		return nil, err
	}
	return lane.RedisStrsToInt64(roomStr)
}

func SetNEAllRoomid(rdb *redis.ClusterClient, roomids []int64) error {
	key := "room:mgr"
	// 监视键的变化
	err := rdb.Watch(func(tx *redis.Tx) error {
		// 检查键是否存在
		exists, err := tx.Exists(key).Result()
		if err != nil {
			return err
		}

		// 如果键不存在，则执行写入操作
		if exists == 0 {
			pipe := tx.Pipeline()
			pipe.Del(key)
			for _, member := range roomids {
				pipe.SAdd(key, lane.Int64ToString(member))
			}
			pipe.Expire(key, time.Second*30)
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)

	return err
}

func EXAllRoomid(rdb *redis.ClusterClient) (bool, error) {
	key := "room:mgr"
	return EXKey(rdb, key)
}

func SetEXAllRoomid(rdb *redis.ClusterClient, roomids []int64) error {
	key := "room:mgr"
	// 监视键的变化
	err := rdb.Watch(func(tx *redis.Tx) error {
		// 检查键是否存在
		exists, err := tx.Exists(key).Result()
		if err != nil {
			return err
		}

		// 如果键存在，则执行写入操作
		if exists != 0 {
			pipe := tx.Pipeline()
			pipe.Del(key)
			for _, member := range roomids {
				pipe.SAdd(key, lane.Int64ToString(member))
			}
			pipe.Expire(key, time.Second*30)
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)

	return err
}

// func AddRoomid(rdb *redis.ClusterClient, roomid lane.Int64) error {
// 	err := rdb.SAdd("room:mgr", roomid.String()).Err()
// 	if err != nil {
// 		log.Fatalf("could not set room info: %v", err)
// 		return err
// 	}
// 	return err
// }

// func DelRoomid(rdb *redis.ClusterClient, roomid lane.Int64) error {
// 	err := rdb.SRem("room:mgr", roomid.String()).Err()
// 	if err != nil {
// 		log.Fatalf("could not del room info: %v", err)
// 		return err
// 	}
// 	return err
// }

func DelAllRoomid(rdb *redis.ClusterClient) error {
	err := rdb.Del("room:mgr").Err()
	if err != nil {
		log.Fatalf("could not del room info: %v", err)
		return err
	}
	return err
}

// room:online
func EXRoomOnline(rdb *redis.ClusterClient, roomid int64) (bool, error) {
	key := fmt.Sprintf("room:online:%s", lane.Int64ToString(roomid))
	return EXKey(rdb, key)
}

func SetEXRoomOnlie(rdb *redis.ClusterClient, roomid int64, onlineCount int) error {
	key := fmt.Sprintf("room:online:%s", lane.Int64ToString(roomid))
	// 监视键的变化
	err := rdb.Watch(func(tx *redis.Tx) error {
		// 检查键是否存在
		exists, err := tx.Exists(key).Result()
		if err != nil {
			return err
		}

		// 如果键存在，则执行写入操作
		if exists != 0 {
			pipe := tx.Pipeline()
			pipe.Set(key, onlineCount, 0).Err()
			pipe.Expire(key, time.Second*30)
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)
	return err
}

func SetNERoomOnlie(rdb *redis.ClusterClient, roomid int64, onlineCount int) error {
	key := fmt.Sprintf("room:online:%s", lane.Int64ToString(roomid))
	// 监视键的变化
	err := rdb.Watch(func(tx *redis.Tx) error {
		// 检查键是否存在
		exists, err := tx.Exists(key).Result()
		if err != nil {
			return err
		}

		// 如果键不存在，则执行写入操作
		if exists == 0 {
			pipe := tx.Pipeline()
			pipe.Set(key, onlineCount, 0).Err()
			pipe.Expire(key, time.Second*30)
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)
	return err
}

func DelRoomOnline(rdb *redis.ClusterClient, roomid int64) error {
	err := rdb.Del(fmt.Sprintf("room:online:%s", lane.Int64ToString(roomid))).Err()
	if err != nil {
		log.Fatalf("could not del room online: %v", err)
		return err
	}
	return err
}

// room:user
func RoomUserid(rdb *redis.ClusterClient, roomid int64) ([]int64, error) {
	userStr, err := rdb.SMembers(fmt.Sprintf("room:userid:%s", lane.Int64ToString(roomid))).Result()
	if len(userStr) == 0 {
		return nil, redis.Nil
	}
	if err != nil {
		return nil, err
	}
	return lane.RedisStrsToInt64(userStr)
}

func RoomCount(rdb *redis.ClusterClient, roomid int64) (int64, error) {
	count, err := rdb.SCard(fmt.Sprintf("room:userid:%s", lane.Int64ToString(roomid))).Result()
	if err != nil {
		return -1, err
	}
	return count, nil
}

func EXRoomUser(rdb *redis.ClusterClient, roomid int64) (bool, error) {
	key := fmt.Sprintf("room:userid:%s", lane.Int64ToString(roomid))
	return EXKey(rdb, key)
}

func SetEXRoomUser(rdb *redis.ClusterClient, roomid int64, userids []int64) error {
	key := fmt.Sprintf("room:userid:%s", lane.Int64ToString(roomid))
	// 监视键的变化
	err := rdb.Watch(func(tx *redis.Tx) error {
		// 检查键是否存在
		exists, err := tx.Exists(key).Result()
		if err != nil {
			return err
		}

		// 如果键存在，则执行写入操作
		if exists != 0 {
			pipe := tx.Pipeline()
			pipe.Del(key)
			for _, member := range userids {
				pipe.SAdd(key, lane.Int64ToString(member))
			}
			pipe.Expire(key, time.Second*30)
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)

	return err
}

func SetNERoomUser(rdb *redis.ClusterClient, roomid int64, userids []int64) error {
	key := fmt.Sprintf("room:userid:%s", lane.Int64ToString(roomid))
	// 监视键的变化
	err := rdb.Watch(func(tx *redis.Tx) error {
		// 检查键是否存在
		exists, err := tx.Exists(key).Result()
		if err != nil {
			return err
		}

		// 如果键不存在，则执行写入操作
		if exists == 0 {
			pipe := tx.Pipeline()
			pipe.Del(key)
			for _, member := range userids {
				pipe.SAdd(key, lane.Int64ToString(member))
			}
			pipe.Expire(key, time.Second*30)
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)

	return err
}

func SetNERoomUserNotExist(rdb *redis.ClusterClient, roomid int64, userids []int64) error {
	key := fmt.Sprintf("room:userid:%s", lane.Int64ToString(roomid))
	// 监视键的变化
	err := rdb.Watch(func(tx *redis.Tx) error {
		// 检查键是否存在
		exists, err := tx.Exists(key).Result()
		if err != nil {
			return err
		}

		// 如果键不存在，则执行写入操作
		if exists == 0 {
			pipe := tx.Pipeline()
			pipe.Del(key)
			for _, member := range userids {
				pipe.SAdd(key, lane.Int64ToString(member))
			}
			pipe.Expire(key, time.Second*30)
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)

	return err
}

// func DelRoomUser(rdb *redis.ClusterClient, roomid lane.Int64, userid lane.Int64) error {
// 	err := rdb.SRem(fmt.Sprintf("room:userid:%s", roomid.String()), userid.String()).Err()
// 	if err != nil {
// 		log.Fatalf("could not del room user: %v", err)
// 		return err
// 	}
// 	return err
// }

func DelRoomAllUser(rdb *redis.ClusterClient, roomid int64) error {
	err := rdb.Del(fmt.Sprintf("room:userid:%s", lane.Int64ToString(roomid))).Err()
	if err != nil {
		log.Fatalf("could not del room all user: %v", err)
		return err
	}
	return err
}

// room:comet
func RoomComet(rdb *redis.ClusterClient, roomid int64) ([]string, error) {
	cometStr, err := rdb.SMembers(fmt.Sprintf("room:comet:%s", lane.Int64ToString(roomid))).Result()
	if len(cometStr) == 0 {
		return nil, redis.Nil
	}
	if err != nil {
		return nil, err
	}
	return cometStr, err
}

func EXRoomComet(rdb *redis.ClusterClient, roomid int64) (bool, error) {
	key := fmt.Sprintf("room:comet:%s", lane.Int64ToString(roomid))
	return EXKey(rdb, key)
}

func SetEXRoomComet(rdb *redis.ClusterClient, roomid int64, comets []string) error {
	key := fmt.Sprintf("room:comet:%s", lane.Int64ToString(roomid))
	// 监视键的变化
	err := rdb.Watch(func(tx *redis.Tx) error {
		// 检查键是否存在
		exists, err := tx.Exists(key).Result()
		if err != nil {
			return err
		}

		// 如果键存在，则执行写入操作
		if exists != 0 {
			pipe := tx.Pipeline()
			pipe.Del(key)
			for _, member := range comets {
				pipe.SAdd(key, member)
			}
			pipe.Expire(key, time.Second*30)
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)

	return err
}

func SetNERoomComet(rdb *redis.ClusterClient, roomid int64, comets []string) error {
	key := fmt.Sprintf("room:comet:%s", lane.Int64ToString(roomid))
	// 监视键的变化
	err := rdb.Watch(func(tx *redis.Tx) error {
		// 检查键是否存在
		exists, err := tx.Exists(key).Result()
		if err != nil {
			return err
		}

		// 如果键不存在，则执行写入操作
		if exists == 0 {
			pipe := tx.Pipeline()
			pipe.Del(key)
			for _, member := range comets {
				pipe.SAdd(key, member)
			}
			pipe.Expire(key, time.Second*30)
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)

	return err
}

func DelRoomAllComet(rdb *redis.ClusterClient, roomid int64) error {
	_, err := rdb.Del(fmt.Sprintf("room:comet:%s", lane.Int64ToString(roomid))).Result()
	if err != nil {
		log.Println("faild to del room all comet", err)
		return err
	}
	return nil
}
