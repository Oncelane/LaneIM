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

// user:online

func (d *SqlDB) UserOnlie(userid int64) (bool, string, error) {
	// 查询所有行的 ID
	var user model.UserMgr
	if err := d.DB.Model(&model.UserMgr{}).Select("online, comet_addr").Where("user_id = ?", userid).First(&user).Error; err != nil {
		return false, "", err
	}
	return user.Online, user.CometAddr, nil
}

func (d *SqlDB) SetUserOnline(userid int64, cometAddr string) error {
	err := d.DB.Model(&model.UserMgr{UserID: userid}).Update("online", true).Error
	if err != nil {
		log.Fatalf("could not set user online: %v", err)
		return err
	}
	d.SetUserComet(userid, cometAddr)
	return nil
}

func (d *SqlDB) SetUseroffline(userid int64) error {
	err := d.DB.Model(&model.UserMgr{UserID: userid}).Update("online", false).Error
	if err != nil {
		log.Fatalf("could not set user online: %v", err)
		return err
	}
	return nil
}

// user:room
func (d *SqlDB) UserRoom(userid int64) ([]int64, error) {
	// 查询所有行的 ID
	var user model.UserMgr
	if err := d.DB.Preload("Rooms").First(&user, userid).Error; err != nil {
		return nil, err
	}
	rt := make([]int64, len(user.Rooms))
	for i, u := range user.Rooms {
		rt[i] = u.RoomID
	}
	return rt, nil
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

func (d *SqlDB) AddUserRoom(userid int64, roomid int64) error {
	err := d.DB.Model(&model.UserMgr{UserID: userid}).Association("Rooms").Append(&model.RoomMgr{RoomID: roomid})
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
	err := d.DB.Model(&model.UserMgr{UserID: userid}).Association("Rooms").Delete(&model.RoomMgr{RoomID: roomid})
	if err != nil {
		log.Println("faild to del user user", err)
		return err
	}
	return nil
}

func (d *SqlDB) DelUserAllRoom(userid int64) error {
	var user model.UserMgr
	if err := d.DB.Preload("Rooms").First(&user, userid).Error; err != nil {
		return err
	}
	// Remove all users from the user
	return d.DB.Model(&user).Association("Rooms").Clear()
}

// user:comet
func (d *SqlDB) UserComet(userid int64) (string, error) {

	user := model.UserMgr{}
	err := d.DB.Model(&model.UserMgr{}).Where("user_id = ?", userid).First(&user).Error
	if err != nil {
		log.Println("failed to query user comet:", err)
		return "", err
	}
	return user.CometAddr, nil
}

func (d *SqlDB) SetUserComet(userid int64, comet string) error {

	err := d.DB.Model(&model.UserMgr{UserID: userid}).Update("comet_addr", comet).Error
	if err != nil {
		log.Println("faild to add user comet", err)
		return err
	}
	return err

}

func (d *SqlDB) DelUserComet(userid int64, comet string) error {
	err := d.DB.Model(&model.UserMgr{UserID: userid}).Update("comet_addr", "").Error
	if err != nil {
		log.Println("faild to del user comet", err)
		return err
	}
	return nil
}

func (d *SqlDB) NewUser(userid int64) error {
	user := model.UserMgr{
		UserID: userid,
	}
	d.DB.Save(&user)
	return nil
}

func (d *SqlDB) DelUser(userid int64) error {
	user := model.UserMgr{
		UserID: userid,
	}
	d.DB.Delete(&user)
	return nil
}
