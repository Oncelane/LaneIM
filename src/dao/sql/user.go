package sql

import (
	"laneIM/src/model"
	"log"

	"gorm.io/gorm"

	lane "laneIM/src/common"
)

//--------------User------------

// user:mgr
func AllUserid(db *gorm.DB) ([]int64, error) {
	// 查询所有行的 ID
	var userids []int64
	err := db.Model(&model.UserMgr{}).Pluck("user_id", &userids).Error
	if err != nil {
		log.Println("faild to query userid", err)
	}
	return userids, nil
}

func SetAllUserid(db *gorm.DB, userids []int64) error {
	// 删除所有现有的记录
	DelAllUserid(db)

	// 插入新的记录
	for _, userID := range userids {
		userMgr := model.UserMgr{UserID: int64(userID)}
		if err := db.Create(&userMgr).Error; err != nil {
			return err
		}
	}
	return nil
}

func AddUserid(db *gorm.DB, userid lane.Int64) error {
	userMgr := model.UserMgr{UserID: int64(userid)}
	if err := db.Create(&userMgr).Error; err != nil {
		return err
	}
	return nil
}

func DelUserid(db *gorm.DB, userid lane.Int64) error {
	userMgr := model.UserMgr{UserID: int64(userid)}
	if err := db.Where("user_id = ?", userid).Delete(&userMgr).Error; err != nil {
		return err
	}
	return nil
}

func DelAllUserid(db *gorm.DB) error {
	if err := db.Where("1 = 1").Delete(&model.UserMgr{}).Error; err != nil {
		return err
	}
	return nil
}

// user:online
func SetUserOnlie(db *gorm.DB, userid lane.Int64, cometAddr string) error {
	err := db.Save(&model.UserOnline{UserID: int64(userid), Online: true}).Error
	if err != nil {
		log.Fatalf("could not set user online: %v", err)
		return err
	}
	AddUserComet(db, userid, cometAddr)
	return nil
}

func SetUserofflie(db *gorm.DB, userid lane.Int64) error {
	err := db.Save(&model.UserOnline{UserID: int64(userid), Online: false}).Error
	if err != nil {
		log.Fatalf("could not set user online: %v", err)
		return err
	}
	return nil
}

func DelUserOnline(db *gorm.DB, userid lane.Int64) error {
	err := db.Where("user_id = ?", int64(userid)).Delete(&model.UserOnline{}).Error
	if err != nil {
		log.Fatalf("could not set user offline: %v", err)
		return err
	}
	return nil
}

// user:room
func UserRoomid(db *gorm.DB, userid lane.Int64) ([]int64, error) {
	// 查询所有行的 ID
	var roomids []int64
	err := db.Model(&model.UserRoom{}).Where("user_id = ?", int64(userid)).Pluck("room_id", &roomids).Error
	if err != nil {
		log.Println("faild to query user roomid", err)
	}
	return roomids, nil
}

func UserRoomCount(db *gorm.DB, userid lane.Int64) (int, error) {
	// 查询所有行的 ID
	var roomids []int64
	err := db.Model(&model.UserRoom{}).Where("user_id = ?", int64(userid)).Pluck("room_id", &roomids).Error
	if err != nil {
		log.Println("faild to query user roomid", err)
	}
	return len(roomids), nil
}

func AddUserRoom(db *gorm.DB, userid lane.Int64, roomid lane.Int64) error {
	err := db.Create(&model.UserRoom{UserID: int64(userid), RoomID: int64(roomid)}).Error
	if err != nil {
		log.Println("faild to add user user", err)
		return err
	}
	return nil
}

func DelUserRoom(db *gorm.DB, userid lane.Int64, roomid lane.Int64) error {
	err := db.Where("user_id = ? and room_id = ?", int64(userid), int64(roomid)).Delete(&model.UserRoom{}).Error
	if err != nil {
		log.Println("faild to del user user", err)
		return err
	}
	return nil
}

func DelUserAllRoom(db *gorm.DB, userid lane.Int64) error {
	err := db.Where("user_id = ?", userid).Delete(&model.UserRoom{}).Error
	if err != nil {
		log.Println("faild to del all user user", err)
		return err
	}
	return nil
}

// user:comet
func UserComet(db *gorm.DB, userid lane.Int64) ([]string, error) {
	// 查询所有行的 ID
	var cometAddrs []string
	err := db.Model(&model.UserComet{}).Pluck("comet_addr", &cometAddrs).Error
	if err != nil {
		log.Println("faild to query user comet", err)
	}
	return cometAddrs, nil
}

func AddUserComet(db *gorm.DB, userid lane.Int64, comet string) error {
	err := db.Create(&model.UserComet{UserID: int64(userid), CometAddr: comet}).Error
	if err != nil {
		log.Println("faild to add user comet", err)
		return err
	}
	return err

}

func DelUserComet(db *gorm.DB, userid lane.Int64, comet string) error {
	err := db.Where("user_id = ? AND comet = ?", int64(userid), comet).Delete(&model.UserComet{}).Error
	if err != nil {
		log.Println("faild to del user comet", err)
		return err
	}
	return nil
}

func DelUserAllComet(db *gorm.DB, userid lane.Int64) error {
	err := db.Where("user_id = ?", int64(userid)).Delete(&model.UserComet{}).Error
	if err != nil {
		log.Println("faild to del user comet", err)
		return err
	}
	return nil
}

func NewUser(db *gorm.DB, userid lane.Int64, roomid lane.Int64, serverAddr string) error {
	AddUserid(db, userid)
	SetUserOnlie(db, userid, serverAddr)
	AddUserRoom(db, userid, roomid)
	return nil
}

func DelUser(db *gorm.DB, userid lane.Int64) error {
	DelUserid(db, userid)
	DelUserOnline(db, userid)
	DelUserAllComet(db, userid)
	DelUserAllRoom(db, userid)
	return nil
}
