package sql

import (
	"laneIM/src/model"
	"log"
)

func (d *SqlDB) AddComet(addr string) error {
	comet := model.CometMgr{
		CometAddr: addr,
	}
	err := d.DB.Save(&comet).Error
	if err != nil {
		log.Println("faild to save comet", err)
		return err
	}
	return nil
}

func (d *SqlDB) DelComet(addr string) error {
	comet := model.CometMgr{
		CometAddr: addr,
	}
	err := d.DB.Delete(&comet).Error
	if err != nil {
		log.Println("faild to delete comet", err)
		return err
	}
	return nil
}
