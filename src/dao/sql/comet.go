package sql

import (
	"laneIM/src/model"
	"laneIM/src/pkg/laneLog"
)

func (d *SqlDB) AddComet(addr string) error {
	comet := model.CometMgr{
		CometAddr: addr,
		Valid:     true,
	}
	err := d.DB.Save(&comet).Error
	if err != nil {
		laneLog.Logger.Infoln("faild to save comet", err)
		return err
	}
	return nil
}

func (d *SqlDB) DelComet(addr string) error {
	comet := model.CometMgr{
		CometAddr: addr,
		Valid:     false,
	}
	err := d.DB.Save(&comet).Error
	if err != nil {
		laneLog.Logger.Infoln("faild to delete comet", err)
		return err
	}
	return nil
}
