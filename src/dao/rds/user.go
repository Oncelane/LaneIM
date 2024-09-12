package rds

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"

	"laneIM/src/model"
	"laneIM/src/pkg/laneLog.go"
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
				laneLog.Logger.Infoln("faild to save usermgr")
				return err
			}
			return nil
		}
		return nil
	}, key)
	if err != nil {
		laneLog.Logger.Infoln("faild to set user mgr")
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
				laneLog.Logger.Infoln("faild to save usermgr")
				return err
			}
			return nil
		}
		return nil
	}, key)
	if err != nil {
		laneLog.Logger.Infoln("faild to set user mgr")
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
	laneLog.Logger.Infoln("TODO注意一下Online是个什么东西:", Online)

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

func UserMgrRoomBatch(rdb *redis.ClusterClient, userids []int64) ([][]int64, bool) {
	// Retrieve Users
	rt := make([][]int64, len(userids))
	full := true
	// pipe := rdb.Pipeline()
	for i := range userids {
		strUserid := strconv.FormatInt(userids[i], 36)
		roomSetKey := fmt.Sprintf("user:%s:roomS", strUserid)
		roomIDs, err := rdb.SMembers(roomSetKey).Result()
		// laneLog.Logger.Debugf("redis smembers key[%s] value[%s]", roomSetKey, roomIDs)
		if len(roomIDs) == 0 {
			full = false
			continue
		}
		if err != nil && err != redis.Nil {
			laneLog.Logger.Errorln("faild redis query", err)
			full = false
			continue
		}
		rooms := make([]int64, len(roomIDs))
		for i, idStr := range roomIDs {
			roomID, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				laneLog.Logger.Errorln("faild read roomid", err)
				full = false
				continue
			}
			rooms[i] = roomID
		}
		rt[i] = rooms
	}
	return rt, full
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
		laneLog.Logger.Infoln("faild to set user room", err)
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
		laneLog.Logger.Infoln("faild to set user room", err)
		return err
	}
	return nil
}

func SetNEUSerMgrRoomBatch(rdb *redis.ClusterClient, userids []int64, roomss [][]int64) error {
	pipe := rdb.Pipeline()
	exists := make([]bool, len(userids))
	for i := range userids {
		strRoomid := strconv.FormatInt(userids[i], 36)
		key := "user:" + strRoomid
		exist, err := pipe.Exists(key).Result()
		if err != nil {
			continue
		}
		exists[i] = exist != 0
	}
	start := time.Now()
	_, err := pipe.Exec()
	laneLog.Logger.Debugln("pipe exist spand", time.Since(start))
	if err != nil {
		laneLog.Logger.Infoln("faild to set user room", err)
		return err
	}
	pipe = rdb.Pipeline()

	for i := range userids {
		strRoomid := strconv.FormatInt(userids[i], 36)
		key := "user:" + strRoomid
		// 如果键不存在，则执行写入操作
		if !exists[i] {
			for j := range roomss[i] {
				// laneLog.Logger.Debugf("redis set key[%s] value[%s]", fmt.Sprintf("%s:roomS", key), roomss[i][j])
				pipe.SAdd(fmt.Sprintf("%s:roomS", key), roomss[i][j]).Err()
			}
			pipe.Expire(fmt.Sprintf("%s:roomS", key), time.Second*60)
		}

	}
	start = time.Now()
	_, err = pipe.Exec()
	laneLog.Logger.Debugln("pipe set user room spand", time.Since(start))
	if err != nil {
		laneLog.Logger.Infoln("faild to set user room", err)
		return err
	}
	return nil
}

func SetEXUSerMgrRoomBatch(rdb *redis.ClusterClient, userids []int64, roomss [][]int64) error {
	pipe := rdb.Pipeline()
	exists := make([]bool, len(userids))
	for i := range userids {
		strRoomid := strconv.FormatInt(userids[i], 36)
		key := "user:" + strRoomid
		exist, err := pipe.Exists(key).Result()
		if err != nil {
			continue
		}
		exists[i] = exist != 0
	}
	_, err := pipe.Exec()
	if err != nil {
		laneLog.Logger.Infoln("faild to set user room", err)
		return err
	}
	pipe = rdb.Pipeline()

	for i := range userids {
		strRoomid := strconv.FormatInt(userids[i], 36)
		key := "user:" + strRoomid
		// 如果键存在，则执行写入操作
		if exists[i] {
			for j := range roomss[i] {
				pipe.SAdd(fmt.Sprintf("%s:roomS", key), roomss[i][j]).Err()
			}
			pipe.Expire(fmt.Sprintf("%s:roomS", key), time.Second*60)
		}

	}
	_, err = pipe.Exec()
	if err != nil {
		laneLog.Logger.Infoln("faild to set user room", err)
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
			pipe := tx.Pipeline()
			pipe.HSet(key, "CA", comet).Result()
			pipe.Expire(key, time.Second*60)
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)
	if err != nil {
		laneLog.Logger.Infoln("faild to set user comet", err.Error())
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
			pipe := tx.Pipeline()
			pipe.HSet(key, "CA", comet).Result()
			pipe.Expire(key, time.Second*60)
			_, err := pipe.Exec()
			return err
		}
		return nil
	}, key)
	if err != nil {
		laneLog.Logger.Infoln("faild to set user comet", err.Error())
		return err
	}
	return nil
}
