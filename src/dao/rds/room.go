package rds

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"

	"laneIM/src/model"
	"laneIM/src/pkg/laneLog.go"
)

//--------------Room------------

//room:mgr

func SetNERoomMgr(rdb *redis.ClusterClient, room *model.RoomMgr) error {

	strRoomid := strconv.FormatInt(room.RoomID, 36)
	key := "room:" + strRoomid

	err := rdb.Watch(func(tx *redis.Tx) error {
		// 检查键是否存在
		exists, err := tx.Exists(key).Result()
		if err != nil {
			return err
		}

		// 如果键不存在，则执行写入操作
		if exists == 0 {
			pipe := rdb.Pipeline()
			pipe.HMSet(key, map[string]interface{}{"OC": room.OnlineCount}).Err()
			for i := range room.Users {
				pipe.SAdd(fmt.Sprintf("%s:userS", key), room.Users[i].UserID).Err()
			}
			for i := range room.Comets {
				pipe.SAdd(fmt.Sprintf("%s:cometS", key), room.Comets[i].CometAddr).Err()
			}
			pipe.Expire(key, time.Second*60)
			pipe.Expire(fmt.Sprintf("%s:userS", key), time.Second*60)
			pipe.Expire(fmt.Sprintf("%s:cometS", key), time.Second*60)

			_, err := pipe.Exec()
			if err != nil {
				laneLog.Logger.Infoln("faild to save roommgr")
				return err
			}
			return nil
		}
		return nil
	}, key)

	return err

}

func SetEXRoomMgr(rdb *redis.ClusterClient, room *model.RoomMgr) error {

	strRoomid := strconv.FormatInt(room.RoomID, 36)
	key := "room:" + strRoomid

	err := rdb.Watch(func(tx *redis.Tx) error {
		// 检查键是否存在
		exists, err := tx.Exists(key).Result()
		if err != nil {
			return err
		}

		// 如果键存在，则执行写入操作
		if exists != 0 {
			pipe := rdb.Pipeline()
			pipe.Del(fmt.Sprintf("%s:userS", key))
			pipe.Del(fmt.Sprintf("%s:cometS", key))
			pipe.HMSet(key, map[string]interface{}{"OC": room.OnlineCount}).Err()
			for i := range room.Users {
				pipe.SAdd(fmt.Sprintf("%s:userS", key), room.Users[i].UserID).Err()
			}
			for i := range room.Comets {
				pipe.SAdd(fmt.Sprintf("%s:cometS", key), room.Comets[i].CometAddr).Err()
			}
			pipe.Expire(key, time.Second*60)
			pipe.Expire(fmt.Sprintf("%s:userS", key), time.Second*60)
			pipe.Expire(fmt.Sprintf("%s:cometS", key), time.Second*60)

			_, err := pipe.Exec()
			if err != nil {
				laneLog.Logger.Infoln("faild to save roommgr")
				return err
			}
			return nil
		}
		return nil
	}, key)

	return err

}

func RoomMgr(rdb *redis.ClusterClient, roomid int64) (*model.RoomMgr, error) {
	strRoomid := strconv.FormatInt(roomid, 36)
	hashKey := fmt.Sprintf("room:%s", strRoomid)

	// Retrieve OnlineCount
	onlineCountStr, err := rdb.HGet(hashKey, "OC").Result()
	if err != nil {
		return nil, err
	}
	onlineCount, err := strconv.Atoi(onlineCountStr)
	if err != nil {
		return nil, err
	}

	// Retrieve Users
	userSetKey := fmt.Sprintf("room:%s:userS", strRoomid)
	userIDs, err := rdb.SMembers(userSetKey).Result()

	if err != nil {
		return nil, err
	}

	users := make([]model.UserMgr, len(userIDs))
	for i, idStr := range userIDs {
		userID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, err
		}
		users[i] = model.UserMgr{UserID: userID}
	}

	// Retrieve Comets
	cometSetKey := fmt.Sprintf("room:%s:cometS", strRoomid)
	cometAddrs, err := rdb.SMembers(cometSetKey).Result()
	if err != nil {
		return nil, err
	}

	comets := make([]model.CometMgr, len(cometAddrs))
	for i, addr := range cometAddrs {
		comets[i] = model.CometMgr{CometAddr: addr}
	}

	return &model.RoomMgr{
		RoomID:      roomid,
		Users:       users,
		Comets:      comets,
		OnlineCount: onlineCount,
	}, nil
}

func RoomMgrComet(rdb *redis.ClusterClient, roomid int64) ([]string, error) {
	// Retrieve Comets
	strRoomid := strconv.FormatInt(roomid, 36)
	cometSetKey := fmt.Sprintf("room:%s:cometS", strRoomid)
	cometAddrs, err := rdb.SMembers(cometSetKey).Result()
	if len(cometAddrs) == 0 {
		return nil, redis.Nil
	}
	if err != nil {
		return nil, err
	}
	return cometAddrs, nil
}

func SetEXRoomMgrComet(rdb *redis.ClusterClient, roomid int64, cometAddrs []string) error {
	strRoomid := strconv.FormatInt(roomid, 36)
	key := "room:" + strRoomid

	err := rdb.Watch(func(tx *redis.Tx) error {
		// 检查键是否存在
		exists, err := tx.Exists(key).Result()
		if err != nil {
			return err
		}

		// 如果键存在，则执行写入操作
		if exists != 0 {
			pipe := rdb.Pipeline()
			pipe.Del(fmt.Sprintf("%s:cometS", key))
			for i := range cometAddrs {
				pipe.SAdd(fmt.Sprintf("%s:cometS", key), cometAddrs[i]).Err()
			}
			pipe.Expire(fmt.Sprintf("%s:cometS", key), time.Second*60)
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)
	if err != nil {
		laneLog.Logger.Infoln("faild to set room comet")
		return err
	}
	return nil
}

func SetNERoomMgrComet(rdb *redis.ClusterClient, roomid int64, cometAddrs []string) error {
	strRoomid := strconv.FormatInt(roomid, 36)
	key := "room:" + strRoomid

	err := rdb.Watch(func(tx *redis.Tx) error {
		// 检查键是否存在
		exists, err := tx.Exists(key).Result()
		if err != nil {
			return err
		}

		// 如果键不存在，则执行写入操作
		if exists == 0 {
			pipe := rdb.Pipeline()
			for i := range cometAddrs {
				pipe.SAdd(fmt.Sprintf("%s:cometS", key), cometAddrs[i]).Err()
			}
			pipe.Expire(fmt.Sprintf("%s:cometS", key), time.Second*60)
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)
	if err != nil {
		laneLog.Logger.Infoln("faild to set room comet")
		return err
	}
	return nil
}

// room:user
func RoomMgrUserid(rdb *redis.ClusterClient, roomid int64) ([]int64, error) {
	// Retrieve Comets
	// Retrieve Users
	strRoomid := strconv.FormatInt(roomid, 36)
	userSetKey := fmt.Sprintf("room:%s:userS", strRoomid)
	userIDs, err := rdb.SMembers(userSetKey).Result()
	if len(userIDs) == 0 {
		return nil, redis.Nil
	}
	if err != nil {
		return nil, err
	}

	users := make([]int64, len(userIDs))
	for i, idStr := range userIDs {
		userID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return nil, err
		}
		users[i] = userID
	}
	return users, nil
}

func SetEXRoomMgrUsers(rdb *redis.ClusterClient, roomid int64, users []int64) error {
	strRoomid := strconv.FormatInt(roomid, 36)
	key := "room:" + strRoomid

	err := rdb.Watch(func(tx *redis.Tx) error {
		// 检查键是否存在
		exists, err := tx.Exists(key).Result()
		if err != nil {
			return err
		}

		// 如果键存在，则执行写入操作
		if exists != 0 {
			pipe := rdb.Pipeline()
			pipe.Del(fmt.Sprintf("%s:userS", key))
			for i := range users {
				pipe.SAdd(fmt.Sprintf("%s:userS", key), users[i]).Err()
			}
			pipe.Expire(fmt.Sprintf("%s:userS", key), time.Second*60)
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)
	if err != nil {
		laneLog.Logger.Infoln("faild to set room user", err)
		return err
	}
	return nil
}

func SetNERoomMgrUsers(rdb *redis.ClusterClient, roomid int64, users []int64) error {
	strRoomid := strconv.FormatInt(roomid, 36)
	key := "room:" + strRoomid

	err := rdb.Watch(func(tx *redis.Tx) error {
		// 检查键是否存在
		exists, err := tx.Exists(key).Result()
		if err != nil {
			return err
		}

		// 如果键不存在，则执行写入操作
		if exists == 0 {
			pipe := rdb.Pipeline()
			for i := range users {
				pipe.SAdd(fmt.Sprintf("%s:userS", key), users[i]).Err()
			}
			pipe.Expire(fmt.Sprintf("%s:userS", key), time.Second*60)
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)
	if err != nil {
		laneLog.Logger.Infoln("faild to set room user", err)
		return err
	}
	return nil
}
