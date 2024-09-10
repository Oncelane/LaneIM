package sql

import (
	"errors"
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

func (d *SqlDB) NewRoom(roomid int64, userid int64, serverAddr string) error {
	roommgr := model.RoomMgr{
		RoomID:      roomid,
		OnlineCount: 1,
	}
	err := d.DB.Create(&roommgr).Error
	if err != nil {
		if sqlerr, ok := err.(*mysql2.MySQLError); ok {
			if sqlerr.Number == 1062 {
				return nil
			}
		}
		log.Println("faild to add room user", err)
		return err
	}
	err = d.AddRoomComet(roomid, serverAddr)
	if err != nil {
		return err
	}
	err = d.AddRoomUser(roomid, userid)
	if err != nil {
		return err
	}
	return nil

}

func (d *SqlDB) DelRoom(roomid int64) (int64, error) {

	// 使用事务处理删除操作
	err := d.DB.Delete(&model.RoomMgr{RoomID: roomid}).Error
	if err != nil {
		return 0, err
	}
	return 0, nil
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

// func (d *SqlDB) SetAllRoomid(roomids []int64) error {
// 	// 删除所有现有的记录
// 	d.DelAllRoomid()

// 	// 插入新的记录
// 	for _, roomID := range roomids {
// 		roomMgr := model.RoomMgr{RoomID: roomID}
// 		if err := d.DB.Create(&roomMgr).Error; err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (d *SqlDB) AddRoomid(roomid int64) error {
// 	roomMgr := model.RoomMgr{RoomID: roomid}
// 	err := d.DB.Save(&roomMgr).Error
// 	if err != nil {
// 		if sqlerr, ok := err.(*mysql2.MySQLError); ok {
// 			if sqlerr.Number == 1062 {
// 				return nil
// 			}
// 		}
// 		log.Println("faild to add room user", err)
// 		return err
// 	}
// 	return nil
// }

// func (d *SqlDB) DelRoomid(roomid int64) error {
// 	roomMgr := model.RoomMgr{RoomID: roomid}
// 	if err := d.DB.Where("room_id = ?", roomid).Delete(&roomMgr).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

func (d *SqlDB) DelAllRoomid() error {
	if err := d.DB.Where("1 = 1").Delete(&model.RoomMgr{}).Error; err != nil {
		return err
	}
	return nil
}

// room:online
func (d *SqlDB) SetRoomOnlie(roomid int64, onlineCount int) error {
	err := d.DB.Model(&model.RoomMgr{RoomID: roomid}).Update("online_count", onlineCount).Error
	if err != nil {
		log.Fatalf("could not set room online count: %v", err)
		return err
	}
	return nil
}

// room:user
func (d *SqlDB) RoomUserid(roomid int64) ([]int64, error) {
	// 查询所有行的 ID
	var room model.RoomMgr
	if err := d.DB.Preload("Users").First(&room, roomid).Error; err != nil {
		return nil, err
	}
	rt := make([]int64, len(room.Users))
	for i, u := range room.Users {
		rt[i] = u.UserID
	}
	return rt, nil
}

// func (d *SqlDB) RoomCount(roomid int64) (int, error) {
// 	// 查询所有行的 ID
// 	var userids []int64
// 	err := d.DB.Model(&model.RoomUserid{}).Where("room_id = ?", roomid).Pluck("user_id", &userids).Error
// 	if err != nil {
// 		log.Println("faild to query room userid", err)
// 	}
// 	return len(userids), nil
// }

func (d *SqlDB) AddRoomUser(roomid int64, userid int64) error {
	err := d.DB.Model(&model.RoomMgr{RoomID: roomid}).Association("Users").Append(&model.UserMgr{UserID: userid})
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
	err := d.DB.Model(&model.RoomMgr{RoomID: roomid}).Association("Users").Delete(&model.UserMgr{UserID: userid})
	if err != nil {
		log.Println("faild to del room user", err)
		return err
	}
	return nil
}

func (d *SqlDB) DelRoomAllUser(roomid int64) error {
	var room model.RoomMgr
	if err := d.DB.Preload("Users").First(&room, roomid).Error; err != nil {
		return err
	}
	// Remove all users from the room
	return d.DB.Model(&room).Association("Users").Clear()
}

// room:comet
func (d *SqlDB) RoomComet(roomid int64) ([]string, error) {
	// 查询所有行的 ID
	room := model.RoomMgr{}
	var cometAddrs []string
	err := d.DB.Preload("Comets").First(&room, roomid).Error
	if err != nil {
		log.Println("faild to query room comet", err)
	}
	return cometAddrs, nil
}

func (d *SqlDB) AddRoomComet(roomid int64, comet string) error {
	err := d.DB.Model(&model.RoomMgr{RoomID: roomid}).Association("Comets").Append(&model.RoomComet{CometAddr: comet})
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

func (d *SqlDB) AddRoomCometWithUserid(userid int64, cometAddr string) error {

	// 开始事务
	tx := d.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 查询所有包含该用户的房间
	var rooms []model.RoomMgr
	if err := tx.Joins("JOIN room_users ON room_id = room_users.room_mgr_room_id").
		Where("room_users.user_mgr_user_id = ?", userid).
		Find(&rooms).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 遍历所有房间，插入或更新网关地址
	for _, room := range rooms {
		// 检查是否已存在该网关地址的记录
		var existingRoomComet model.RoomComet
		err := tx.Where("room_id = ? AND comet_addr = ?", room.RoomID, cometAddr).First(&existingRoomComet).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			// 查询出错
			tx.Rollback()
			return err
		}

		if err == gorm.ErrRecordNotFound {
			// 如果记录不存在，则插入新记录
			if err := tx.Create(&model.RoomComet{
				RoomID:    room.RoomID,
				CometAddr: cometAddr,
			}).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (d *SqlDB) DelRoomComet(roomid int64, comet string) error {
	err := d.DB.Model(&model.RoomMgr{RoomID: roomid}).Association("Comets").Delete(&model.RoomComet{CometAddr: comet})
	if err != nil {
		log.Println("faild to del room comet", err)
		return err
	}
	return nil
}
