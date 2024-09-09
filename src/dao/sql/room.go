package sql

import (
	"fmt"
	"laneIM/src/config"
	"laneIM/src/model"
	"log"

	mysql2 "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func DB(conf config.Mysql) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", conf.Username, conf.Password, conf.Addr, conf.DataBase)
	log.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableAutomaticPing:   true,
		SkipDefaultTransaction: true, // 关闭默认事务
	})
	if err != nil {
		log.Fatalln("faild to get db")
		return nil
	}
	sqldb, err := db.DB()
	sqldb.SetMaxOpenConns(100)
	sqldb.SetMaxIdleConns(50)

	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
		return nil
	}
	return db
}

func RoomNew(db *gorm.DB, roomid int64, userid int64, serverAddr string) error {
	roommgr := model.RoomMgr{
		RoomID: roomid,
	}
	roomOnline := model.RoomOnline{
		RoomID:      roomid,
		OnlineCount: 1,
	}
	roomComet := model.RoomComet{
		RoomID:    roomid,
		CometAddr: serverAddr,
	}
	roomUserid := model.RoomUserid{
		RoomID: roomid,
		UserID: userid,
	}
	rt := db.Create(&roommgr)
	if rt.Error != nil {
		log.Fatalf("could not set room info: %v", rt.Error)
		return rt.Error
	}
	rt = db.Create(&roomOnline)
	if rt.Error != nil {
		log.Fatalf("could not set room info: %v", rt.Error)
		return rt.Error
	}
	rt = db.Create(&roomComet)
	if rt.Error != nil {
		log.Fatalf("could not set room info: %v", rt.Error)
		return rt.Error
	}
	rt = db.Create(&roomUserid)
	if rt.Error != nil {
		log.Fatalf("could not set room info: %v", rt.Error)
		return rt.Error
	}
	return nil

}

// func SqlRoomNew(db *gorm.DB, roomid int64, userid int64, serverAddr string) error {
// 	roommgr := model.RoomMgr{
// 		RoomID: roomid,
// 	}
// 	roomOnline := model.RoomOnline{
// 		RoomID:      roomid,
// 		OnlineCount: 1,
// 	}
// 	roomComet := model.RoomComet{
// 		RoomID:    roomid,
// 		CometAddr: serverAddr,
// 	}
// 	roomUserid := model.RoomUserid{
// 		RoomID: roomid,
// 		UserID: userid,
// 	}
// 	err := db.Transaction(func(tx *gorm.DB) error {
// 		rt := tx.Create(&roommgr)
// 		if rt.Error != nil {
// 			log.Fatalf("could not set room info: %v", rt.Error)
// 			return rt.Error
// 		}
// 		rt = tx.Create(&roomOnline)
// 		if rt.Error != nil {
// 			log.Fatalf("could not set room info: %v", rt.Error)
// 			return rt.Error
// 		}
// 		rt = tx.Create(&roomComet)
// 		if rt.Error != nil {
// 			log.Fatalf("could not set room info: %v", rt.Error)
// 			return rt.Error
// 		}
// 		rt = tx.Create(&roomUserid)
// 		if rt.Error != nil {
// 			log.Fatalf("could not set room info: %v", rt.Error)
// 			return rt.Error
// 		}
// 		return nil
// 	})

// 	return err
// }

func RoomDel(db *gorm.DB, roomid int64) (int64, error) {

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

//--------------Room------------

// room:mgr
func AllRoomid(db *gorm.DB) ([]int64, error) {
	// 查询所有行的 ID
	var roomids []int64
	result := db.Model(&model.RoomMgr{}).Pluck("room_id", &roomids)
	if result.Error != nil {
		panic(result.Error)
	}
	return roomids, nil
}

func SetAllRoomid(db *gorm.DB, roomids []int64) error {
	// 删除所有现有的记录
	DelAllRoomid(db)

	// 插入新的记录
	for _, roomID := range roomids {
		roomMgr := model.RoomMgr{RoomID: roomID}
		if err := db.Create(&roomMgr).Error; err != nil {
			return err
		}
	}
	return nil
}

func AddRoomid(db *gorm.DB, roomid int64) error {
	roomMgr := model.RoomMgr{RoomID: roomid}
	if err := db.Create(&roomMgr).Error; err != nil {
		return err
	}
	return nil
}

func DelRoomid(db *gorm.DB, roomid int64) error {
	roomMgr := model.RoomMgr{RoomID: roomid}
	if err := db.Where("room_id = ?", roomid).Delete(&roomMgr).Error; err != nil {
		return err
	}
	return nil
}

func DelAllRoomid(db *gorm.DB) error {
	if err := db.Where("1 = 1").Delete(&model.RoomMgr{}).Error; err != nil {
		return err
	}
	return nil
}

// room:online
func SetRoomOnlie(db *gorm.DB, roomid int64, onlineCount int) error {
	err := db.Save(&model.RoomOnline{RoomID: roomid, OnlineCount: onlineCount}).Error
	if err != nil {
		log.Fatalf("could not set room online: %v", err)
		return err
	}
	return nil
}

func DelRoomOnline(db *gorm.DB, roomid int64) error {
	err := db.Where("room_id = ?", roomid).Delete(&model.RoomOnline{}).Error
	if err != nil {
		log.Fatalf("could not del room online: %v", err)
		return err
	}
	return nil
}

// room:user
func RoomUserid(db *gorm.DB, roomid int64) ([]int64, error) {
	// 查询所有行的 ID
	var userids []int64
	err := db.Model(&model.RoomUserid{}).Where("room_id = ?", roomid).Pluck("user_id", &userids).Error
	if err != nil {
		log.Println("faild to query room userid", err)
	}
	return userids, nil
}

func RoomCount(db *gorm.DB, roomid int64) (int, error) {
	// 查询所有行的 ID
	var userids []int64
	err := db.Model(&model.RoomUserid{}).Where("room_id = ?", roomid).Pluck("user_id", &userids).Error
	if err != nil {
		log.Println("faild to query room userid", err)
	}
	return len(userids), nil
}

func AddRoomUser(db *gorm.DB, roomid int64, userid int64) error {
	err := db.Create(&model.RoomUserid{RoomID: roomid, UserID: userid}).Error
	if err != nil {
		if sqlerr, ok := err.(*mysql2.MySQLError); ok {
			if sqlerr.Number == 1062 {
				return nil
			}
		}
		log.Println("faild to add room user", err)
		return err
	}
	return nil
}

func DelRoomUser(db *gorm.DB, roomid int64, userid int64) error {
	err := db.Where("room_id = ? and user_id = ?", roomid, userid).Delete(&model.RoomUserid{}).Error
	if err != nil {
		log.Println("faild to del room user", err)
		return err
	}
	return nil
}

func DelRoomAllUser(db *gorm.DB, roomid int64) error {
	err := db.Where("room_id = ?", roomid).Delete(&model.RoomUserid{}).Error
	if err != nil {
		log.Println("faild to del all room user", err)
		return err
	}
	return nil
}

// room:comet
func RoomComet(db *gorm.DB, roomid int64) ([]string, error) {
	// 查询所有行的 ID
	var cometAddrs []string
	err := db.Model(&model.RoomComet{}).Pluck("comet_addr", &cometAddrs).Error
	if err != nil {
		log.Println("faild to query room comet", err)
	}
	return cometAddrs, nil
}

func AddRoomComet(db *gorm.DB, roomid int64, comet string) error {
	err := db.Create(&model.RoomComet{RoomID: roomid, CometAddr: comet}).Error
	if err != nil {
		if sqlerr, ok := err.(*mysql2.MySQLError); ok {
			if sqlerr.Number == 1062 {
				return nil
			}
		}
		log.Println("faild to add room comet", err)
		return err
	}
	return err

}

func AddRoomCometWithUserid(db *gorm.DB, userid int64, cometAddr string) error {

	tx := db.Begin()

	// 插入网关地址到所有包含用户的房间
	err := tx.Exec(`
INSERT INTO room_comets (room_id, comet_addr)
SELECT r.room_id, ?
FROM room_userids r
LEFT JOIN room_comets c
ON r.room_id = c.room_id AND c.comet_addr = ?
WHERE r.user_id = ?
AND c.comet_addr IS NULL
`, cometAddr, cometAddr, userid).Error

	if err != nil {
		tx.Rollback()
		// 处理错误
		log.Println("faild to add comet to all room")
		return err
	} else {
		tx.Commit()
	}
	return nil
}

func DelRoomComet(db *gorm.DB, roomid int64, comet string) error {
	err := db.Where("room_id = ? AND comet = ?", roomid, comet).Delete(&model.RoomComet{}).Error
	if err != nil {
		log.Println("faild to del room comet", err)
		return err
	}
	return nil
}

func DelRoomAllComet(db *gorm.DB, roomid int64) error {
	err := db.Where("room_id = ?", roomid).Delete(&model.RoomComet{}).Error
	if err != nil {
		log.Println("faild to del room comet", err)
		return err
	}
	return nil
}

func NewRoom(db *gorm.DB, roomid int64, userid int64, serverAddr string) error {
	AddRoomid(db, roomid)
	SetRoomOnlie(db, roomid, 1)
	AddRoomUser(db, roomid, userid)
	return nil
}

func DelRoom(db *gorm.DB, roomid int64) error {

	DelRoomid(db, roomid)
	DelRoomOnline(db, roomid)
	DelRoomAllComet(db, roomid)
	DelRoomAllUser(db, roomid)
	return nil
}
