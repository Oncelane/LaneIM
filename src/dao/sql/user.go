package sql

import (
	"fmt"
	"laneIM/src/model"
	"laneIM/src/pkg/laneLog"
	"log"
	"strconv"

	"gorm.io/gorm"
)

// --------------User------------
func (d *SqlDB) NewUser(userid int64) error {
	user := model.UserMgr{
		UserID: userid,
	}
	d.DB.Save(&user)
	return nil
}

func (d *SqlDB) NewUserBatch(_ *gorm.DB, userids []int64) error {

	var users []model.UserMgr
	i := 0
	for {
		tx := d.DB.Begin()
		left := len(userids) - i
		if left > 2000 {
			users = make([]model.UserMgr, 2000)
		} else if left != 0 {
			users = make([]model.UserMgr, left)
		} else {
			break
		}
		for ; i < len(users); i++ {
			users[i].UserID = userids[i]
		}
		tx.Save(&users)
		err := tx.Commit().Error
		if err != nil {
			laneLog.Logger.Fatalln("[server] faild to commit new user comet", err.Error())
		}
	}
	return nil
}

func (d *SqlDB) DelUser(userid int64) error {
	user := model.UserMgr{
		UserID: userid,
	}
	d.DB.Delete(&user)
	return nil
}

func (d *SqlDB) DelUserBatch(userids []int64) error {
	users := make([]model.UserMgr, len(userids))
	for i := range userids {
		users[i] = model.UserMgr{
			UserID: userids[i],
		}
	}
	d.DB.Delete(&users)
	return nil
}

func (d *SqlDB) UserMgr(userid int64) (*model.UserMgr, error) {
	user := model.UserMgr{}
	err := d.DB.Preload("Rooms").First(&user, userid).Error
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to sql userMgr", err)
		return nil, err
	}
	return &user, nil
}

// user:mgr
func (d *SqlDB) AllUserid() ([]int64, error) {
	// 查询所有行的 ID
	var userids []int64
	err := d.DB.Model(&model.UserMgr{}).Pluck("user_id", &userids).Error
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to query userid", err)
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

func (d *SqlDB) SetUserOnlineBatch(tx *gorm.DB, userids []int64, cometAddr string) error {
	var end bool
	if tx == nil {
		end = true
		tx = d.DB.Begin()
	}

	// 构建 SQL 查询
	if len(userids) == 0 {
		return nil // 如果没有 userids，直接返回
	}

	// 使用 GORM 的 `Where` 和 `Updates` 方法进行批量更新
	err := tx.Model(&model.UserMgr{}).
		Where("user_id IN ?", userids).
		Update("online", true).Error
	if err != nil {
		tx.Rollback()
		laneLog.Logger.Infoln("failed to set user online rollback", err)
		return err
	}

	// 调用 SetUserCometBatch 函数
	d.SetUserCometBatch(tx, userids, cometAddr)

	if end {
		err := tx.Commit().Error
		if err != nil {
			laneLog.Logger.Infoln("failed to commit set user online", err.Error())
		}
	}
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
func (d *SqlDB) SetUserOfflineBatch(tx *gorm.DB, userids []int64) error {
	var end bool
	if tx == nil {
		end = true
		tx = d.DB.Begin()
	}
	for i := range userids {
		err := tx.Model(&model.UserMgr{UserID: userids[i]}).Update("online", false).Error
		if err != nil {
			tx.Rollback()
			laneLog.Logger.Fatalln("[server] faild to set user online rollback", err)
			return err
		}
	}
	if end {
		err := tx.Commit().Error
		if err != nil {
			laneLog.Logger.Fatalln("[server] faild to commit set user online", err.Error())
		}
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

func (d *SqlDB) UserRoomBatch(userids []int64) ([][]int64, error) {

	var users []model.UserMgr
	if err := d.DB.Where("user_id IN ?", userids).Preload("Rooms").Find(&users).Error; err != nil {
		return nil, err
	}
	results := make([][]int64, len(userids))
	for i, user := range users {
		rt := make([]int64, len(user.Rooms))
		for j, room := range user.Rooms {
			rt[j] = room.RoomID
		}
		results[i] = rt
	}
	return results, nil
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

func (d *SqlDB) DelUserRoom(userid int64, roomid int64) error {
	err := d.DB.Model(&model.UserMgr{UserID: userid}).Association("Rooms").Delete(&model.RoomMgr{RoomID: roomid})
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to del user user", err)
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
		laneLog.Logger.Infoln("failed to query user comet:", err)
		return "", err
	}
	return user.CometAddr, nil
}

func (d *SqlDB) SetUserComet(userid int64, comet string) error {

	err := d.DB.Model(&model.UserMgr{UserID: userid}).Update("comet_addr", comet).Error
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to add user comet", err)
		return err
	}
	return err

}

func (d *SqlDB) SetUserCometBatch(tx *gorm.DB, userids []int64, comet string) error {
	var end bool
	if tx == nil {
		end = true
		tx = d.DB.Begin()
	}

	for i := range userids {
		err := tx.Model(&model.UserMgr{UserID: userids[i]}).Update("comet_addr", comet).Error
		if err != nil {
			laneLog.Logger.Fatalln("[server] faild rollback set user comet batch")
			tx.Rollback()
			return err
		}
	}

	if end {
		err := tx.Commit().Error
		if err != nil {
			laneLog.Logger.Fatalln("[server] faild to commit set user comet", err.Error())
		}
	}
	return nil

}

func (d *SqlDB) DelUserComet(userid int64, comet string) error {
	err := d.DB.Model(&model.UserMgr{UserID: userid}).Update("comet_addr", "").Error
	if err != nil {
		laneLog.Logger.Fatalln("[server] faild to del user comet", err)
		return err
	}
	return nil
}
