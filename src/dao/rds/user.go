package rds

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"

	"laneIM/src/model"
)

//room:mgr

func SetNEUserMgr(rdb *redis.ClusterClient, user *model.UserMgr) error {

	strUserid := strconv.FormatInt(user.UserID, 36)
	key := "user:" + strUserid

	err := rdb.Watch(func(tx *redis.Tx) error {
		// 检查键是否存在
		exists, err := tx.Exists(key).Result()
		if err != nil {
			return err
		}

		// 如果键不存在，则执行写入操作
		if exists == 0 {
			pipe := rdb.Pipeline()
			pipe.HMSet(key, map[string]interface{}{"CA": user.CometAddr, "OL": user.Online}).Err()
			for i := range user.Rooms {
				pipe.SAdd(fmt.Sprintf("%s:roomS", key), user.Rooms[i].RoomID).Err()
			}
			pipe.Expire(key, time.Second*60)
			pipe.Expire(fmt.Sprintf("%s:roomS", key), time.Second*60)

			_, err := pipe.Exec()
			if err != nil {
				log.Println("faild to save usermgr")
				return err
			}
			return nil
		}
		return nil
	}, key)
	if err != nil {
		log.Println("faild to set user mgr")
		return err
	}
	return err

}

func SetEXUserMgr(rdb *redis.ClusterClient, user *model.UserMgr) error {

	strUserid := strconv.FormatInt(user.UserID, 36)
	key := "user:" + strUserid

	err := rdb.Watch(func(tx *redis.Tx) error {
		// 检查键是否存在
		exists, err := tx.Exists(key).Result()
		if err != nil {
			return err
		}

		// 如果键存在，则执行写入操作
		if exists != 0 {
			pipe := rdb.Pipeline()
			pipe.Del(fmt.Sprintf("%s:roomS", key))
			pipe.HMSet(key, map[string]interface{}{"CA": user.CometAddr, "OL": user.Online}).Err()
			for i := range user.Rooms {
				pipe.SAdd(fmt.Sprintf("%s:roomS", key), user.Rooms[i].RoomID).Err()
			}
			pipe.Expire(key, time.Second*60)
			pipe.Expire(fmt.Sprintf("%s:roomS", key), time.Second*60)

			_, err := pipe.Exec()
			if err != nil {
				log.Println("faild to save usermgr")
				return err
			}
			return nil
		}
		return nil
	}, key)
	if err != nil {
		log.Println("faild to set user mgr")
		return err
	}
	return err

}

func UserMgr(rdb *redis.ClusterClient, userid int64) (*model.UserMgr, error) {
	strUserid := strconv.FormatInt(userid, 36)
	hashKey := fmt.Sprintf("user:%s", strUserid)

	// Retrieve OnlineCount
	CometAddr, err := rdb.HGet(hashKey, "CA").Result()
	if err != nil {
		return nil, err
	}
	Online, err := rdb.HGet(hashKey, "OL").Result()
	if err != nil {
		return nil, err
	}
	log.Println("TODO注意一下Online是个什么东西:", Online)

	// Retrieve Users
	roomSetKey := fmt.Sprintf("user:%s:roomS", strUserid)
	roomIDs, err := rdb.SMembers(roomSetKey).Result()
	if err != nil {
		return nil, err
	}

	rooms := make([]model.RoomMgr, len(roomIDs))
	for i, idStr := range roomIDs {
		roomID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, err
		}
		rooms[i] = model.RoomMgr{RoomID: roomID}
	}

	return &model.UserMgr{
		UserID:    userid,
		Rooms:     rooms,
		CometAddr: CometAddr,
	}, nil
}

// user:room
func UserMgrRoom(rdb *redis.ClusterClient, userid int64) ([]int64, error) {
	// Retrieve Users
	strUserid := strconv.FormatInt(userid, 36)
	roomSetKey := fmt.Sprintf("user:%s:roomS", strUserid)
	roomIDs, err := rdb.SMembers(roomSetKey).Result()
	if len(roomIDs) == 0 {
		return nil, redis.Nil
	}
	if err != nil {
		return nil, err
	}
	rooms := make([]int64, len(roomIDs))
	for i, idStr := range roomIDs {
		roomID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, err
		}
		rooms[i] = roomID
	}
	return rooms, nil
}

func SetNEUSerMgrRoom(rdb *redis.ClusterClient, userid int64, rooms []int64) error {
	strRoomid := strconv.FormatInt(userid, 36)
	key := "user:" + strRoomid

	err := rdb.Watch(func(tx *redis.Tx) error {
		// 检查键是否存在
		exists, err := tx.Exists(key).Result()
		if err != nil {
			return err
		}

		// 如果键不存在，则执行写入操作
		if exists == 0 {
			pipe := rdb.Pipeline()
			for i := range rooms {
				pipe.SAdd(fmt.Sprintf("%s:roomS", key), rooms[i]).Err()
			}
			pipe.Expire(fmt.Sprintf("%s:roomS", key), time.Second*60)
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)
	if err != nil {
		log.Println("faild to set user room", err)
		return err
	}
	return nil
}

func SetEXUSerMgrRoom(rdb *redis.ClusterClient, userid int64, rooms []int64) error {
	strUserid := strconv.FormatInt(userid, 36)
	key := "user:" + strUserid

	err := rdb.Watch(func(tx *redis.Tx) error {
		// 检查键是否存在
		exists, err := tx.Exists(key).Result()
		if err != nil {
			return err
		}

		// 如果键存在，则执行写入操作
		if exists != 0 {
			pipe := rdb.Pipeline()
			pipe.Del(fmt.Sprintf("%s:roomS", key))
			for i := range rooms {
				pipe.SAdd(fmt.Sprintf("%s:roomS", key), rooms[i]).Err()
			}
			pipe.Expire(fmt.Sprintf("%s:roomS", key), time.Second*60)
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)
	if err != nil {
		log.Println("faild to set user room", err)
		return err
	}
	return nil
}

// user:comet
func UserMgrComet(rdb *redis.ClusterClient, userid int64) (string, error) {
	strUserid := strconv.FormatInt(userid, 36)
	hashKey := fmt.Sprintf("user:%s", strUserid)

	// Retrieve OnlineCount
	CometAddr, err := rdb.HGet(hashKey, "CA").Result()
	if err != nil {
		return "", err
	}
	return CometAddr, err
}

func SetNEUSerMgrComet(rdb *redis.ClusterClient, userid int64, comet string) error {
	strUserid := strconv.FormatInt(userid, 36)
	key := fmt.Sprintf("user:%s", strUserid)
	err := rdb.Watch(func(tx *redis.Tx) error {
		// 检查键是否存在
		exists, err := tx.Exists(key).Result()
		if err != nil {
			return err
		}

		// 如果键不存在，则执行写入操作
		if exists == 0 {
			pipe := rdb.Pipeline()
			pipe.HSet(key, "CA", comet).Result()
			pipe.Expire(key, time.Second*60)
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)
	if err != nil {
		log.Println("faild to set user comet", err)
		return err
	}
	return nil
}

func SetEXUSerMgrComet(rdb *redis.ClusterClient, userid int64, comet string) error {
	strUserid := strconv.FormatInt(userid, 36)
	key := fmt.Sprintf("user:%s", strUserid)

	err := rdb.Watch(func(tx *redis.Tx) error {
		// 检查键是否存在
		exists, err := tx.Exists(key).Result()
		if err != nil {
			return err
		}

		// 如果键存在，则执行写入操作
		if exists != 0 {
			pipe := rdb.Pipeline()
			pipe.HSet(key, "CA", comet).Result()
			pipe.Expire(key, time.Second*60)
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)
	if err != nil {
		log.Println("faild to set user comet", err)
		return err
	}
	return nil
}
