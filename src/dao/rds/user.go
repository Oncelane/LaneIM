package rds

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"

	lane "laneIM/src/common"
)

func EXKey(rdb *redis.ClusterClient, key string) (bool, error) {
	exists, err := rdb.Exists(key).Result()
	if err != nil {
		return false, err
	}
	return exists != 0, nil
}

// user:mgr

func AllUserid(rdb *redis.ClusterClient) ([]int64, error) {
	userStr, err := rdb.SMembers("user:mgr").Result()
	if len(userStr) == 0 {
		return nil, redis.Nil
	}
	if err != nil {
		return nil, err
	}
	return lane.RedisStrsToInt64(userStr)
}

func EXAllUserid(rdb *redis.ClusterClient) (bool, error) {
	return EXKey(rdb, "user:mgr")
}

func SetEXAllUserid(rdb *redis.ClusterClient, userids []int64) error {

	key := "user:mgr"
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
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)

	return err
}

func SetNEAllUserid(rdb *redis.ClusterClient, userids []int64) error {

	key := "user:mgr"
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
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)

	return err
}

// func UserNew(rdb *redis.ClusterClient, userid int64) error {
// 	err := rdb.SAdd("user:mgr", lane.Int64ToString(userid)).Err()
// 	if err != nil {
// 		log.Fatalf("could not set room info: %v", err)
// 		return err
// 	}

// 	return nil
// }

// func UserDel(rdb *redis.ClusterClient, userid int64) (int64, error) {
// 	err := rdb.SRem("user:mgr", lane.Int64ToString(userid)).Err()
// 	if err != nil {
// 		log.Fatalf("could not set room info: %v", err)
// 		return 0, err
// 	}
// 	rt, err := rdb.Del(fmt.Sprintf("user:online:%s", lane.Int64ToString(userid))).Result()
// 	if err != nil {
// 		return rt, err
// 	}
// 	rt, err = rdb.Del(fmt.Sprintf("user:comet:%s", lane.Int64ToString(userid))).Result()
// 	if err != nil {
// 		return rt, err
// 	}
// 	rt, err = rdb.Del(fmt.Sprintf("user:room:%s", lane.Int64ToString(userid))).Result()
// 	if err != nil {
// 		return rt, err
// 	}
// 	return rt, nil
// }

// user:online
func UserOnline(rdb *redis.ClusterClient, userid int64) (bool, string, error) {
	rt, err := rdb.Get(fmt.Sprintf("user:online:%s", lane.Int64ToString(userid))).Int64()
	if err != nil {
		return false, "", err
	}
	if rt != 1 {
		return false, "", nil
	}
	cometAddr, err := UserComet(rdb, userid)
	return true, cometAddr, err
}

func EXUserOnline(rdb *redis.ClusterClient, userid int64) (bool, error) {
	key := fmt.Sprintf("user:online:%s", lane.Int64ToString(userid))
	return EXKey(rdb, key)
}

func SetEXUserOnline(rdb *redis.ClusterClient, userid int64, comet string) error {

	key := fmt.Sprintf("user:online:%s", lane.Int64ToString(userid))
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
			pipe.Set(fmt.Sprintf("user:online:%s", lane.Int64ToString(userid)), 1, 0).Err()
			pipe.Expire(key, time.Second*3)
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)
	if err != nil {
		return err
	}
	err = SetEXUserComet(rdb, userid, comet)
	return err
}
func SetNEUserOnline(rdb *redis.ClusterClient, userid int64, comet string) error {
	key := fmt.Sprintf("user:online:%s", lane.Int64ToString(userid))
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
			pipe.Set(fmt.Sprintf("user:online:%s", lane.Int64ToString(userid)), 1, 0).Err()
			pipe.Expire(key, time.Second*3)
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)
	if err != nil {
		return err
	}
	err = SetNEUserComet(rdb, userid, comet)
	return err
}

func SetEXUserOffline(rdb *redis.ClusterClient, userid int64, comet string) error {

	key := fmt.Sprintf("user:online:%s", lane.Int64ToString(userid))
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
			pipe.Set(fmt.Sprintf("user:online:%s", lane.Int64ToString(userid)), 0, 0).Err()
			pipe.Expire(key, time.Second*3)
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)
	if err != nil {
		return err
	}
	err = SetEXUserComet(rdb, userid, comet)
	return err
}
func SetNEUserOffline(rdb *redis.ClusterClient, userid int64, comet string) error {
	key := fmt.Sprintf("user:online:%s", lane.Int64ToString(userid))
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
			pipe.Set(fmt.Sprintf("user:online:%s", lane.Int64ToString(userid)), 0, 0).Err()
			pipe.Expire(key, time.Second*3)
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)
	if err != nil {
		return err
	}
	err = SetNEUserComet(rdb, userid, comet)
	return err
}

func DelUserOnline(rdb *redis.ClusterClient, userid int64, comet string) error {
	err := rdb.Del(fmt.Sprintf("user:online:%s", lane.Int64ToString(userid))).Err()
	if err != nil {
		return err
	}
	return nil
}

// user:room
func UserRoom(rdb *redis.ClusterClient, userid int64) ([]int64, error) {
	roomStr, err := rdb.SMembers(fmt.Sprintf("user:room:%s", lane.Int64ToString(userid))).Result()
	if len(roomStr) == 0 {
		return nil, redis.Nil
	}
	if err != nil {
		return nil, err
	}
	return lane.RedisStrsToInt64(roomStr)
}

func EXUserRoom(rdb *redis.ClusterClient, userid int64) (bool, error) {
	key := fmt.Sprintf("user:room:%s", lane.Int64ToString(userid))
	return EXKey(rdb, key)
}

func SetEXUserRoom(rdb *redis.ClusterClient, userid int64, roomids []int64) error {
	key := fmt.Sprintf("user:room:%s", lane.Int64ToString(userid))
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
			pipe.Expire(key, time.Second*3)
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)

	return err
}

func SetNEUserRoom(rdb *redis.ClusterClient, userid int64, roomids []int64) error {
	key := fmt.Sprintf("user:room:%s", lane.Int64ToString(userid))
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
			pipe.Expire(key, time.Second*3)
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)

	return err
}

func DelUserAllRoom(rdb *redis.ClusterClient, userid int64) error {
	_, err := rdb.Del(fmt.Sprintf("user:room:%s", lane.Int64ToString(userid))).Result()
	if err != nil {
		log.Println("faild to del user all room", err)
		return err
	}
	return nil
}

// user:comet
func UserComet(rdb *redis.ClusterClient, userid int64) (string, error) {
	cometStr, err := rdb.Get(fmt.Sprintf("user:comet:%s", lane.Int64ToString(userid))).Result()
	if err != nil {
		return "", err
	}
	return cometStr, err
}

func EXUserComet(rdb *redis.ClusterClient, userid int64) (bool, error) {
	key := fmt.Sprintf("user:comet:%s", lane.Int64ToString(userid))
	return EXKey(rdb, key)
}

func SetEXUserComet(rdb *redis.ClusterClient, userid int64, comet string) error {

	key := fmt.Sprintf("user:comet:%s", lane.Int64ToString(userid))
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
			pipe.Set(key, comet, 0)
			pipe.Expire(key, time.Second*3)
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)

	return err
}

func SetNEUserComet(rdb *redis.ClusterClient, userid int64, comet string) error {

	key := fmt.Sprintf("user:comet:%s", lane.Int64ToString(userid))
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
			pipe.Set(key, comet, 0)
			pipe.Expire(key, time.Second*3)
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)

	return err
}

func DelUserComet(rdb *redis.ClusterClient, userid int64) error {
	_, err := rdb.Del(fmt.Sprintf("user:comet:%s", lane.Int64ToString(userid))).Result()
	if err != nil {
		log.Println("faild to del user comet", err)
		return err
	}
	return nil
}
