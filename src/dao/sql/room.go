package sql

import (
	"fmt"
	"laneIM/src/config"
	"laneIM/src/dao/rds"
	"laneIM/src/model"
	"laneIM/src/pkg/laneLog"
	"laneIM/src/pkg/mergewrite"
	"log"
	"strconv"

	"github.com/go-redis/redis/v8"
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
		laneLog.Logger.Infoln("faild to add room user", err)
		return err
	}
	// laneLog.Logger.Warnf("find the time of room create roomid[%d] userid[%d] server[%d]", roomid, userid, serverAddr)
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
func (d *SqlDB) RoomMgr(roomid int64) (*model.RoomMgr, error) {
	room := model.RoomMgr{}
	err := d.DB.Preload("Users").Preload("Comets").First(&room, roomid).Error
	if err != nil {
		laneLog.Logger.Infoln("faild to sql roomMgr")
		return nil, err
	}
	return &room, nil
}

func (d *SqlDB) AllRoomid() ([]int64, error) {
	// 查询所有行的 ID
	var roomids []int64
	result := d.DB.Model(&model.RoomMgr{}).Pluck("room_id", &roomids)
	if result.Error != nil {
		laneLog.Logger.Infoln("faild to query room:mgr:k", result.Error)
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
// 		laneLog.Logger.Infoln("faild to add room user", err)
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

// singleflight
func (d *SqlDB) RoomUserSingleflight(roomid int64) ([]int64, error) {
	key := "room:user:" + strconv.FormatInt(roomid, 36)
	r, err, _ := d.SingleFlightGroup.Do(key, func() (any, error) {
		return d.RoomUserid(roomid)
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

func (d *SqlDB) AddRoomUser(roomid int64, userid int64) error {
	err := d.DB.Model(&model.RoomMgr{RoomID: roomid}).Association("Users").Append(&model.UserMgr{UserID: userid})
	if err != nil {
		if sqlerr, ok := err.(*mysql2.MySQLError); ok {
			if sqlerr.Number == 1062 {
				return nil
			}
		}
		laneLog.Logger.Infoln("faild to add room user", err)
		return err
	}
	return nil
}

func (d *SqlDB) AddRoomUserBatch(tx *gorm.DB, roomids []int64, userids []int64) error {
	var end bool
	if tx == nil {
		end = true
		tx = d.DB.Begin()
	}
	// laneLog.Logger.Warnf("add rooms[%v] user[%v] ", roomids, userids)
	for i := range roomids {

		err := tx.Model(&model.RoomMgr{RoomID: roomids[i]}).Association("Users").Append(&model.UserMgr{UserID: userids[i]})
		if err != nil {
			tx.Rollback()
			laneLog.Logger.Errorln("faild commit ,roll back", err)
			return err
		}
	}

	if end {
		err := tx.Commit().Error
		if err != nil {
			laneLog.Logger.Infoln("faild to commit batch room user")
		}
		return err
	}

	return nil
}

func (d *SqlDB) DelRoomUser(roomid int64, userid int64) error {
	err := d.DB.Model(&model.RoomMgr{RoomID: roomid}).Association("Users").Delete(&model.UserMgr{UserID: userid})
	if err != nil {
		laneLog.Logger.Infoln("faild to del room user", err)
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
	err := d.DB.Preload("Comets").First(&room, roomid).Error
	if err != nil {
		laneLog.Logger.Infoln("faild to query room comet", err)
	}
	rt := make([]string, len(room.Comets))
	for i, u := range room.Comets {
		rt[i] = u.CometAddr
	}
	return rt, nil
}

func (d *SqlDB) RoomCometBatch(roomids []int64) ([][]string, error) {
	var rooms []model.RoomMgr
	err := d.DB.Preload("Comets").Where("room_id IN ?", roomids).Find(&rooms).Error
	if err != nil {
		laneLog.Logger.Infoln("failed to query room comets", err)
		return nil, err
	}

	rt := make([][]string, len(rooms))
	for i, room := range rooms {
		rt[i] = make([]string, len(room.Comets))
		for j, comet := range room.Comets {
			rt[i][j] = comet.CometAddr
		}
	}
	return rt, nil
}

func (d *SqlDB) AddRoomComet(roomid int64, comet string) error {
	err := d.DB.Model(&model.RoomMgr{RoomID: roomid}).Association("Comets").Append(&model.CometMgr{CometAddr: comet})
	if err != nil {
		if sqlerr, ok := err.(*mysql2.MySQLError); ok {
			if sqlerr.Number == 1062 {
				return nil
			}
		}
		laneLog.Logger.Infoln("faild to add room comet", err)
		return err
	}
	return err
}

func (d *SqlDB) AddRoomCometBatch(roomids []int64, comet string) error {
	tx := d.DB.Begin()
	for _, roomid := range roomids {

		err := tx.Model(&model.RoomMgr{RoomID: roomid}).Association("Comets").Append(&model.CometMgr{CometAddr: comet})
		if err != nil {
			if sqlerr, ok := err.(*mysql2.MySQLError); ok {
				if sqlerr.Number == 1062 {
					// 发生重复键错误时跳过
					continue
				}
			}
			laneLog.Logger.Fatalln("failed to add room comet for room ID:", roomid, "comet:", comet, "error:", err)
			return err
		}

	}
	err := tx.Commit().Error
	return err
}

// func (d *SqlDB) AddRoomCometWithUserid(userid int64, cometAddr string) error {
// 	query := `
// INSERT IGNORE INTO
//     room_comets (
//         room_mgr_room_id,
//         comet_mgr_comet_addr
//     )
// SELECT room_mgr_room_id, ?
// FROM room_users
//     JOIN room_mgrs ON room_users.room_mgr_room_id = room_mgrs.room_id
// WHERE
//     room_users.user_mgr_user_id = ?
// `
// 	err := d.DB.Exec(query, cometAddr, userid).Error
// 	if err != nil {
// 		if sqlerr, ok := err.(*mysql2.MySQLError); ok {
// 			if sqlerr.Number == 1062 {
// 				return nil
// 			}
// 		}
// 		laneLog.Logger.Infoln("faild to add comet to room", err)
// 		return err
// 	}

// 	return nil

// }

func (d *SqlDB) AddRoomCometWithUseridBatch(rdb *redis.ClusterClient, userids []int64, cometAddr string) error {
	// start := time.Now()

	rooms, err := d.UserRoomBatch(userids)
	// laneLog.Logger.Debugln("userids", userids)
	// laneLog.Logger.Debugln("rooms:", rooms)
	if err != nil {
		laneLog.Logger.Fatalln("faild to UserRoomBatch ", err)
	}
	r := make(map[int64]struct{})
	for _, roomss := range rooms {
		for _, room := range roomss {
			r[room] = struct{}{}
		}
	}
	rr := make([]int64, 0)
	for roomid := range r {
		rr = append(rr, roomid)
	}
	// laneLog.Logger.Debugln("r", r)
	// laneLog.Logger.Debugln("rr", rr)
	rds.AddRoomcometBatch(rdb, rr, cometAddr)
	if err != nil {
		laneLog.Logger.Fatalln("faild to commit add room user's comet", err)
		return err
	}

	err = d.AddRoomCometBatch(rr, cometAddr)
	if err != nil {
		laneLog.Logger.Fatalln("faild to commit add room user's comet", err)
		return err
	}
	return nil

}

func (d *SqlDB) DelRoomComet(roomid int64, comet string) error {
	err := d.DB.Model(&model.RoomMgr{RoomID: roomid}).Association("Comets").Delete(&model.CometMgr{CometAddr: comet})
	if err != nil {
		laneLog.Logger.Infoln("faild to del room comet", err)
		return err
	}
	return nil
}
