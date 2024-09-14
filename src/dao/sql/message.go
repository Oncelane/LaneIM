package sql

import "laneIM/src/model"

func (d *SqlDB) AddMessageBatch(ins []*model.GroupMessage) error {
	return d.DB.Create(ins).Error
}
