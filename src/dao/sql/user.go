package sql

import (
	"laneIM/src/model"
	"log"

	mysql2 "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
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
	for _, userid := range userids {
		userMgr := model.UserMgr{UserID: userid}
		if err := db.Create(&userMgr).Error; err != nil {
			return err
		}
	}
	return nil
}

func NewUserid(db *gorm.DB, userid int64) error {
	userMgr := model.UserMgr{UserID: userid}
	err := db.Create(&userMgr).Error
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

func DelUserid(db *gorm.DB, userid int64) error {
	userMgr := model.UserMgr{UserID: userid}
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

func UserOnlie(db *gorm.DB, userid int64) (bool, string, error) {
	// 查询所有行的 ID
	var exist bool
	if err := db.Model(&model.UserOnline{}).Where("user_id = ?", userid).Select("online").First(&exist).Error; err != nil {
		return false, "", err
	}
	cometAddr, err := UserComet(db, userid)
	return exist, cometAddr, err
}

func SetUserOnline(db *gorm.DB, userid int64, cometAddr string) error {
	err := db.Model(&model.UserOnline{}).Where("user_id = ?", userid).Update("online", true).Error
	if err != nil {
		log.Fatalf("could not set user online: %v", err)
		return err
	}
	SetUserComet(db, userid, cometAddr)
	return nil
}

func SetUseroffline(db *gorm.DB, userid int64) error {
	err := db.Model(&model.UserOnline{}).Select("user_id = ?", userid).Update("online", false).Error
	if err != nil {
		log.Fatalf("could not set user online: %v", err)
		return err
	}
	return nil
}

func DelUserOnline(db *gorm.DB, userid int64) error {
	err := db.Where("user_id = ?", userid).Delete(&model.UserOnline{}).Error
	if err != nil {
		log.Fatalf("could not set user offline: %v", err)
		return err
	}
	return nil
}

// user:room
func UserRoom(db *gorm.DB, userid int64) ([]int64, error) {
	// 查询所有行的 ID
	var roomids []int64
	err := db.Model(&model.UserRoom{}).Where("user_id = ?", userid).Pluck("room_id", &roomids).Error
	if err != nil {
		log.Println("faild to query user roomid", err)
	}
	return roomids, nil
}

func UserRoomCount(db *gorm.DB, userid int64) (int, error) {
	// 查询所有行的 ID
	var roomids []int64
	err := db.Model(&model.UserRoom{}).Where("user_id = ?", userid).Pluck("room_id", &roomids).Error
	if err != nil {
		log.Println("faild to query user roomid", err)
	}
	return len(roomids), nil
}

func AddUserRoom(db *gorm.DB, userid int64, roomid int64) error {
	err := db.Save(&model.UserRoom{UserID: userid, RoomID: roomid}).Error
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

func DelUserRoom(db *gorm.DB, userid int64, roomid int64) error {
	err := db.Where("user_id = ? and room_id = ?", userid, roomid).Delete(&model.UserRoom{}).Error
	if err != nil {
		log.Println("faild to del user user", err)
		return err
	}
	return nil
}

func DelUserAllRoom(db *gorm.DB, userid int64) error {
	err := db.Where("user_id = ?", userid).Delete(&model.UserRoom{}).Error
	if err != nil {
		log.Println("faild to del all user user", err)
		return err
	}
	return nil
}

// user:comet
func UserComet(db *gorm.DB, userid int64) (string, error) {
	// 查询所有行的 ID
	var cometAddr string
	err := db.Model(&model.UserComet{}).Where("user_id = ?", userid).Select("comet_addr").First(&cometAddr).Error
	if err != nil {
		log.Println("faild to query user comet", err)
	}
	return cometAddr, nil
}

func SetUserComet(db *gorm.DB, userid int64, comet string) error {
	err := db.Save(&model.UserComet{UserID: userid, CometAddr: comet}).Error
	if err != nil {
		log.Println("faild to add user comet", err)
		return err
	}
	return err

}

func DelUserComet(db *gorm.DB, userid int64, comet string) error {
	err := db.Where("user_id = ? AND comet = ?", userid, comet).Delete(&model.UserComet{}).Error
	if err != nil {
		log.Println("faild to del user comet", err)
		return err
	}
	return nil
}

func DelUserAllComet(db *gorm.DB, userid int64) error {
	err := db.Where("user_id = ?", userid).Delete(&model.UserComet{}).Error
	if err != nil {
		log.Println("faild to del user comet", err)
		return err
	}
	return nil
}

func NewUser(db *gorm.DB, userid int64) error {
	NewUserid(db, userid)
	// SetUserOnlie(db, userid, serverAddr)
	// AddUserRoom(db, userid, roomid)
	return nil
}

func DelUser(db *gorm.DB, userid int64) error {
	DelUserid(db, userid)
	DelUserOnline(db, userid)
	DelUserAllComet(db, userid)
	DelUserAllRoom(db, userid)
	return nil
}
