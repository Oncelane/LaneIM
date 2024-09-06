package dao

import (
	lane "laneIM/src/common"
	"laneIM/src/model"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SqlDB() *gorm.DB {
	dsn := "debian-sys-maint:FJho5xokpFqZygL5@tcp(127.0.0.1:3306)/laneIM?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	return db
}

func SqlAllRoomid(db *gorm.DB) ([]int64, error) {
	// 查询所有行的 ID
	var roomids []int64
	result := db.Model(&model.RoomMgr{}).Pluck("room_id", &roomids)
	if result.Error != nil {
		panic(result.Error)
	}
	return roomids, nil
}

func SqlRoomNew(db *gorm.DB, roomid lane.Int64, userid lane.Int64, serverAddr string) error {
	roommgr := model.RoomMgr{
		RoomID: int64(roomid),
	}
	roomOnline := model.RoomOnline{
		RoomID:      int64(roomid),
		OnlineCount: 1,
	}
	roomComet := model.RoomComet{
		RoomID:    int64(roomid),
		CometAddr: serverAddr,
	}
	roomUserid := model.RoomUserid{
		RoomID: int64(roomid),
		UserID: int64(userid),
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		rt := tx.Create(&roommgr)
		if rt.Error != nil {
			log.Fatalf("could not set room info: %v", rt.Error)
			return rt.Error
		}
		rt = tx.Create(&roomOnline)
		if rt.Error != nil {
			log.Fatalf("could not set room info: %v", rt.Error)
			return rt.Error
		}
		rt = tx.Create(&roomComet)
		if rt.Error != nil {
			log.Fatalf("could not set room info: %v", rt.Error)
			return rt.Error
		}
		rt = tx.Create(&roomUserid)
		if rt.Error != nil {
			log.Fatalf("could not set room info: %v", rt.Error)
			return rt.Error
		}
		return nil
	})

	return err
}

func SqlRoomDel(db *gorm.DB, roomid lane.Int64) (lane.Int64, error) {

	// 使用事务处理删除操作
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("room_id = ?", roomid).Delete(&model.RoomMgr{}).Error; err != nil {
			return err
		}
		if err := tx.Where("room_id = ?", roomid).Delete(&model.RoomOnline{}).Error; err != nil {
			return err
		}
		if err := tx.Where("room_id = ?", roomid).Delete(&model.RoomComet{}).Error; err != nil {
			return err
		}
		if err := tx.Where("room_id = ?", roomid).Delete(&model.RoomUserid{}).Error; err != nil {
			return err
		}
		return nil
	})
	return 0, err
}

// func RoomJoinUser(db *gorm.DB, roomid lane.Int64, userid lane.Int64) error {
// 	_, err := db.SAdd(fmt.Sprintf("room:userid:%s", roomid.String()), userid.String()).Result()
// 	return err
// }

// func RoomQuitUser(db *gorm.DB, roomid lane.Int64, userid lane.Int64) error {
// 	_, err := db.SRem(fmt.Sprintf("room:userid:%s", roomid.String()), userid.String()).Result()
// 	return err
// }

// func RoomQueryUserid(db *gorm.DB, roomid lane.Int64) ([]int64, error) {
// 	userStr, err := db.SMembers(fmt.Sprintf("room:userid:%s", roomid.String())).Result()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return RedisStrsToInt64(userStr)
// }

// func RoomPutComet(db *gorm.DB, roomid lane.Int64, comet string) error {
// 	_, err := db.SAdd(fmt.Sprintf("room:comet:%s", roomid.String()), comet).Result()
// 	return err
// }

// func RoomDelComet(db *gorm.DB, roomid lane.Int64, comet string) error {
// 	_, err := db.SRem(fmt.Sprintf("room:comet:%s", roomid.String()), comet).Result()
// 	return err
// }
// func RoomQueryComet(db *gorm.DB, roomid lane.Int64) ([]string, error) {
// 	cometStr, err := db.SMembers(fmt.Sprintf("room:comet:%s", roomid.String())).Result()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return cometStr, err
// }
// func RoomCount(db *gorm.DB, roomid lane.Int64) (int64, error) {
// 	count, err := db.SCard(fmt.Sprintf("room:userid:%s", roomid.String())).Result()
// 	if err != nil {
// 		return -1, err
// 	}
// 	return count, nil
// }

// func AllUserid(db *gorm.DB) ([]int64, error) {
// 	userStr, err := db.SMembers("userMgr").Result()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return RedisStrsToInt64(userStr)
// }

// func UpdateRoom(db *gorm.DB, roomid int64) (*msg.RoomInfo, error) {
// 	roomInfo := &msg.RoomInfo{}
// 	comets, err := RoomQueryComet(db, lane.Int64(roomid))
// 	if err != nil {
// 		return nil, err
// 	}

// 	cometMap := make(map[string]bool)
// 	for _, addr := range comets {
// 		cometMap[addr] = true
// 	}
// 	userMap := make(map[int64]bool)
// 	userids, err := RoomQueryUserid(db, lane.Int64(roomid))
// 	if err != nil {
// 		return nil, err
// 	}
// 	for _, userid := range userids {
// 		userMap[userid] = true
// 	}
// 	roomInfo.OnlineNum = -1
// 	roomInfo.Server = cometMap
// 	roomInfo.Users = userMap
// 	return roomInfo, nil
// }

// // func AtomicRedis(db *gorm.DB, key string, value string) error {
// // 	// Lua脚本，用于实现原子性的读取-修改-写入
// // 	// KEYS[1] 是键名，ARGV[1] 是新值
// // 	script := `
// //     local currentValue = redis.call("GET", KEYS[1])
// //     if currentValue == false then
// //         -- 如果键不存在，则直接设置新值
// //         redis.call("SET", KEYS[1], ARGV[1])
// //     else
// //         -- 如果键存在，则进行某种修改（这里假设只是打印出来，实际可以修改后再设置）
// //         -- 注意：这里为了示例简单，并没有真正修改值，只是打印
// //         print("Current value:", currentValue)
// //         -- 真实场景中，你可能需要基于currentValue计算新值，然后设置
// //         -- redis.call("SET", KEYS[1], newValue)
// //     end
// //     return currentValue
// //     `

// // 	// 执行Lua脚本
// // 	result, err := db.Eval(script, []string{key}, new).Result()
// // 	if err != nil {
// // 		panic(err)
// // 	}

// // 	fmt.Println("Result:", result)
// // 	return nil
// // }

// func RedisStrsToInt64(strs []string) ([]int64, error) {
// 	ret := make([]int64, len(strs))
// 	var tmp lane.Int64
// 	for i, str := range strs {
// 		tmp.PasrseString(str)
// 		ret[i] = int64(tmp)
// 	}
// 	return ret, nil
// }

// // Create RoomMgr
// func CreateRoomMgr() {
// 	roomMgr :=model.RoomMgr{
// 		ID: 0,
// 		RoomID: ,
// 	}

// 	db.Create(&roomMgr)
// 	c.JSON(http.StatusOK, roomMgr)
// }

// // Get RoomMgr by ID
// func GetRoomMgr() {
// 	id := c.Param("id")
// 	var roomMgr model.RoomMgr
// 	if err := db.First(&roomMgr, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "RoomMgr not found"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, roomMgr)
// }

// // Update RoomMgr by ID
// func UpdateRoomMgr() {
// 	id := c.Param("id")
// 	var roomMgr model.RoomMgr
// 	if err := db.First(&roomMgr, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "RoomMgr not found"})
// 		return
// 	}
// 	if err := c.BindJSON(&roomMgr); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
// 	db.Save(&roomMgr)
// 	c.JSON(http.StatusOK, roomMgr)
// }

// // Delete RoomMgr by ID
// func DeleteRoomMgr() {
// 	id := c.Param("id")
// 	if err := db.Delete(&model.RoomMgr{}, id).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{"error": "RoomMgr not found"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "RoomMgr deleted"})
// }
