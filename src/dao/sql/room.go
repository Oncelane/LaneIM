package sql

import (
	"fmt"
	"laneIM/src/config"
	"laneIM/src/model"
	"laneIM/src/pkg/mergewrite"
	"log"

	mysql2 "github.com/go-sql-driver/mysql"
	"golang.org/x/sync/singleflight"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SqlDB struct {
	DB                *gorm.DB
	MergeWriter       *mergewrite.MergeWriter
	SingleFlightGroup singleflight.Group
}

func NewDB(conf config.Mysql) *SqlDB {
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
	return &SqlDB{
		DB:          db,
		MergeWriter: mergewrite.NewMergeWriter(conf.BatchWriter),
	}
}

func (d *SqlDB) RoomNew(roomid int64, userid int64, serverAddr string) error {
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
	rt := d.DB.Create(&roommgr)
	if rt.Error != nil {
		log.Fatalf("could not set room info: %v", rt.Error)
		return rt.Error
	}
	rt = d.DB.Create(&roomOnline)
	if rt.Error != nil {
		log.Fatalf("could not set room info: %v", rt.Error)
		return rt.Error
	}
	rt = d.DB.Create(&roomComet)
	if rt.Error != nil {
		log.Fatalf("could not set room info: %v", rt.Error)
		return rt.Error
	}
	rt = d.DB.Create(&roomUserid)
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
// 	err := d.DB.Transaction(func(tx *gorm.DB) error {
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

func (d *SqlDB) RoomDel(roomid int64) (int64, error) {

	// 使用事务处理删除操作
	err := d.DB.Transaction(func(tx *gorm.DB) error {
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
func (d *SqlDB) AllRoomid() ([]int64, error) {
	// 查询所有行的 ID
	var roomids []int64
	result := d.DB.Model(&model.RoomMgr{}).Pluck("room_id", &roomids)
	if result.Error != nil {
		log.Println("faild to query room:mgr:k", result.Error)
	}
	return roomids, nil
}

// single flight
func (d *SqlDB) AllRoomidSingleflight() ([]int64, error) {
	r, err, _ := d.SingleFlightGroup.Do("room_mgr", func() (any, error) {
		return d.AllRoomid()
	})
	var (
		rt []int64
		ok bool
	)
	if rt, ok = r.([]int64); !ok {
		return nil, fmt.Errorf("batch return type wrong")
	}
	return rt, err
}

// bath version
// func (d *SqlDB) BatchAllRoomid() ([]int64, error) {
// 	r, err := d.MergeWriter.Do("room_mgr", func() (any, error) {
// 		return d.AllRoomid()
// 	})
// 	var (
// 		rt []int64
// 		ok bool
// 	)
// 	if rt, ok = r.([]int64); !ok {
// 		return nil, fmt.Errorf("batch return type wrong")
// 	}
// 	return rt, err
// }

// 100ms
// func (d *SqlDB) doBatchAllRoomid() ([]int64, error) {
// 	d.MergeWriter.
// }

func (d *SqlDB) SetAllRoomid(roomids []int64) error {
	// 删除所有现有的记录
	d.DelAllRoomid()

	// 插入新的记录
	for _, roomID := range roomids {
		roomMgr := model.RoomMgr{RoomID: roomID}
		if err := d.DB.Create(&roomMgr).Error; err != nil {
			return err
		}
	}
	return nil
}

func (d *SqlDB) AddRoomid(roomid int64) error {
	roomMgr := model.RoomMgr{RoomID: roomid}
	err := d.DB.Save(&roomMgr).Error
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

func (d *SqlDB) DelRoomid(roomid int64) error {
	roomMgr := model.RoomMgr{RoomID: roomid}
	if err := d.DB.Where("room_id = ?", roomid).Delete(&roomMgr).Error; err != nil {
		return err
	}
	return nil
}

func (d *SqlDB) DelAllRoomid() error {
	if err := d.DB.Where("1 = 1").Delete(&model.RoomMgr{}).Error; err != nil {
		return err
	}
	return nil
}

// room:online
func (d *SqlDB) SetRoomOnlie(roomid int64, onlineCount int) error {
	err := d.DB.Save(&model.RoomOnline{RoomID: roomid, OnlineCount: onlineCount}).Error
	if err != nil {
		log.Fatalf("could not set room online: %v", err)
		return err
	}
	return nil
}

func (d *SqlDB) DelRoomOnline(roomid int64) error {
	err := d.DB.Where("room_id = ?", roomid).Delete(&model.RoomOnline{}).Error
	if err != nil {
		log.Fatalf("could not del room online: %v", err)
		return err
	}
	return nil
}

// room:user
func (d *SqlDB) RoomUserid(roomid int64) ([]int64, error) {
	// 查询所有行的 ID
	var userids []int64
	err := d.DB.Model(&model.RoomUserid{}).Where("room_id = ?", roomid).Pluck("user_id", &userids).Error
	if err != nil {
		log.Println("faild to query room userid", err)
	}
	return userids, nil
}

func (d *SqlDB) RoomCount(roomid int64) (int, error) {
	// 查询所有行的 ID
	var userids []int64
	err := d.DB.Model(&model.RoomUserid{}).Where("room_id = ?", roomid).Pluck("user_id", &userids).Error
	if err != nil {
		log.Println("faild to query room userid", err)
	}
	return len(userids), nil
}

func (d *SqlDB) AddRoomUser(roomid int64, userid int64) error {
	err := d.DB.Save(&model.RoomUserid{RoomID: roomid, UserID: userid}).Error
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

func (d *SqlDB) DelRoomUser(roomid int64, userid int64) error {
	err := d.DB.Where("room_id = ? and user_id = ?", roomid, userid).Delete(&model.RoomUserid{}).Error
	if err != nil {
		log.Println("faild to del room user", err)
		return err
	}
	return nil
}

func (d *SqlDB) DelRoomAllUser(roomid int64) error {
	err := d.DB.Where("room_id = ?", roomid).Delete(&model.RoomUserid{}).Error
	if err != nil {
		log.Println("faild to del all room user", err)
		return err
	}
	return nil
}

// room:comet
func (d *SqlDB) RoomComet(roomid int64) ([]string, error) {
	// 查询所有行的 ID
	var cometAddrs []string
	err := d.DB.Model(&model.RoomComet{}).Pluck("comet_addr", &cometAddrs).Error
	if err != nil {
		log.Println("faild to query room comet", err)
	}
	return cometAddrs, nil
}

func (d *SqlDB) AddRoomComet(roomid int64, comet string) error {
	err := d.DB.Create(&model.RoomComet{RoomID: roomid, CometAddr: comet}).Error
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

// func AddRoomCometWithUserid(  userid int64, cometAddr string) error {
// 	var roomIDs []int64
// 	err := d.DB.Model(&model.RoomUserid{}).
// 		Select("room_id").
// 		Where("user_id = ?", userid).
// 		Pluck("room_id", &roomIDs).Error
// 	if err != nil {
// 		// 处理查询错误
// 		log.Println("failed to query room IDs")
// 		return err
// 	}
// 	// d.DB.Model(&)
// 	return nil
// }

func (d *SqlDB) AddRoomCometWithUserid(userid int64, cometAddr string) error {

	tx := d.DB.Begin()

	// 插入网关地址到所有包含用户的房间
	err := tx.Exec(`
INSERT INTO
    room_comets (room_id, comet_addr)
SELECT rm.room_id, ?
FROM room_mgrs rm
    JOIN room_userids ru ON rm.room_id = ru.room_id
WHERE
    ru.user_id = ?
ON DUPLICATE KEY UPDATE
    comet_addr = VALUES(comet_addr);
`, cometAddr, userid).Error

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
func (d *SqlDB) DelRoomComet(roomid int64, comet string) error {
	err := d.DB.Where("room_id = ? AND comet = ?", roomid, comet).Delete(&model.RoomComet{}).Error
	if err != nil {
		log.Println("faild to del room comet", err)
		return err
	}
	return nil
}

func (d *SqlDB) DelRoomAllComet(roomid int64) error {
	err := d.DB.Where("room_id = ?", roomid).Delete(&model.RoomComet{}).Error
	if err != nil {
		log.Println("faild to del room comet", err)
		return err
	}
	return nil
}

func (d *SqlDB) NewRoom(roomid int64, userid int64, serverAddr string) error {
	d.AddRoomid(roomid)
	d.SetRoomOnlie(roomid, 1)
	d.AddRoomUser(roomid, userid)
	return nil
}

func (d *SqlDB) DelRoom(roomid int64) error {

	d.DelRoomid(roomid)
	d.DelRoomOnline(roomid)
	d.DelRoomAllComet(roomid)
	d.DelRoomAllUser(roomid)
	return nil
}
