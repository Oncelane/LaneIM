package sql

import (
	"fmt"
	"laneIM/src/model"
	"log"
	"strconv"

	mysql2 "github.com/go-sql-driver/mysql"
)

// --------------User------------

// user:mgr
func (d *SqlDB) AllUserid() ([]int64, error) {
	// 查询所有行的 ID
	var userids []int64
	err := d.DB.Model(&model.UserMgr{}).Pluck("user_id", &userids).Error
	if err != nil {
		log.Println("faild to query userid", err)
	}
	return userids, nil
}

func (d *SqlDB) SetAllUserid(userids []int64) error {
	// 删除所有现有的记录
	d.DelAllUserid()

	// 插入新的记录
	for _, userid := range userids {
		userMgr := model.UserMgr{UserID: userid}
		if err := d.DB.Create(&userMgr).Error; err != nil {
			return err
		}
	}
	return nil
}

func (d *SqlDB) NewUserid(userid int64) error {
	userMgr := model.UserMgr{UserID: userid}
	err := d.DB.Create(&userMgr).Error
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

func (d *SqlDB) DelUserid(userid int64) error {
	userMgr := model.UserMgr{UserID: userid}
	if err := d.DB.Where("user_id = ?", userid).Delete(&userMgr).Error; err != nil {
		return err
	}
	return nil
}

func (d *SqlDB) DelAllUserid() error {
	if err := d.DB.Where("1 = 1").Delete(&model.UserMgr{}).Error; err != nil {
		return err
	}
	return nil
}

// user:online

func (d *SqlDB) UserOnlie(userid int64) (bool, string, error) {
	// 查询所有行的 ID
	var exist bool
	if err := d.DB.Model(&model.UserOnline{}).Where("user_id = ?", userid).Select("online").First(&exist).Error; err != nil {
		return false, "", err
	}
	cometAddr, err := d.UserComet(userid)
	return exist, cometAddr, err
}

func (d *SqlDB) SetUserOnline(userid int64, cometAddr string) error {
	err := d.DB.Model(&model.UserOnline{}).Where("user_id = ?", userid).Update("online", true).Error
	if err != nil {
		log.Fatalf("could not set user online: %v", err)
		return err
	}
	d.SetUserComet(userid, cometAddr)
	return nil
}

func (d *SqlDB) SetUseroffline(userid int64) error {
	err := d.DB.Model(&model.UserOnline{}).Select("user_id = ?", userid).Update("online", false).Error
	if err != nil {
		log.Fatalf("could not set user online: %v", err)
		return err
	}
	return nil
}

func (d *SqlDB) DelUserOnline(userid int64) error {
	err := d.DB.Where("user_id = ?", userid).Delete(&model.UserOnline{}).Error
	if err != nil {
		log.Fatalf("could not set user offline: %v", err)
		return err
	}
	return nil
}

// user:room
func (d *SqlDB) UserRoom(userid int64) ([]int64, error) {
	// 查询所有行的 ID
	var roomids []int64
	err := d.DB.Model(&model.UserRoom{}).Where("user_id = ?", userid).Pluck("room_id", &roomids).Error
	if err != nil {
		log.Println("faild to query user roomid", err)
	}
	return roomids, nil
}

// singleflight
func (d *SqlDB) UserRoomSingleflight(userid int64) ([]int64, error) {
	key := "user:room:" + strconv.FormatInt(userid, 36)
	r, err, _ := d.SingleFlightGroup.Do(key, func() (any, error) {
		return d.UserRoom(userid)
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

func (d *SqlDB) UserRoomCount(userid int64) (int, error) {
	// 查询所有行的 ID
	var roomids []int64
	err := d.DB.Model(&model.UserRoom{}).Where("user_id = ?", userid).Pluck("room_id", &roomids).Error
	if err != nil {
		log.Println("faild to query user roomid", err)
	}
	return len(roomids), nil
}

func (d *SqlDB) AddUserRoom(userid int64, roomid int64) error {
	err := d.DB.Save(&model.UserRoom{UserID: userid, RoomID: roomid}).Error
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

func (d *SqlDB) DelUserRoom(userid int64, roomid int64) error {
	err := d.DB.Where("user_id = ? and room_id = ?", userid, roomid).Delete(&model.UserRoom{}).Error
	if err != nil {
		log.Println("faild to del user user", err)
		return err
	}
	return nil
}

func (d *SqlDB) DelUserAllRoom(userid int64) error {
	err := d.DB.Where("user_id = ?", userid).Delete(&model.UserRoom{}).Error
	if err != nil {
		log.Println("faild to del all user user", err)
		return err
	}
	return nil
}

// user:comet
func (d *SqlDB) UserComet(userid int64) (string, error) {
	// 查询所有行的 ID
	var cometAddr string
	err := d.DB.Model(&model.UserComet{}).Where("user_id = ?", userid).Select("comet_addr").First(&cometAddr).Error
	if err != nil {
		log.Println("faild to query user comet", err)
	}
	return cometAddr, nil
}

func (d *SqlDB) SetUserComet(userid int64, comet string) error {
	err := d.DB.Save(&model.UserComet{UserID: userid, CometAddr: comet}).Error
	if err != nil {
		log.Println("faild to add user comet", err)
		return err
	}
	return err

}

func (d *SqlDB) elUserComet(userid int64, comet string) error {
	err := d.DB.Where("user_id = ? AND comet = ?", userid, comet).Delete(&model.UserComet{}).Error
	if err != nil {
		log.Println("faild to del user comet", err)
		return err
	}
	return nil
}

func (d *SqlDB) DelUserAllComet(userid int64) error {
	err := d.DB.Where("user_id = ?", userid).Delete(&model.UserComet{}).Error
	if err != nil {
		log.Println("faild to del user comet", err)
		return err
	}
	return nil
}

func (d *SqlDB) NewUser(userid int64) error {
	d.NewUserid(userid)
	// SetUserOnlie(db, userid, serverAddr)
	// AddUserRoom(db, userid, roomid)
	return nil
}

func (d *SqlDB) DelUser(userid int64) error {
	d.DelUserid(userid)
	d.DelUserOnline(userid)
	d.DelUserAllComet(userid)
	d.DelUserAllRoom(userid)
	return nil
}
