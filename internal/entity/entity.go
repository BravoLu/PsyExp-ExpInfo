package entity

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ExpInfo struct {
	gorm.Model
	Title   string         `gorm:"column:title;unique;no null"`
	Desc    string         `gorm:"column:detail;no null"`
	Rid     int64          `gorm:"column:rid"`
	Ctime   int32          `gorm:"column:ctime;no null"`
	Pnum    int32          `gorm:"column:pnum;no null"`
	CurType int32          `gorm:"column:cur_type"`
	Price   int64          `gorm:"column:price"`
	State   int32          `gorm:"column:state"`
	EndTime time.Time      `gorm:"column:end_time"`
	Url     string         `gorm:"column:url"`
	Tags    datatypes.JSON `gorm:"column:tags;type:jsonb,omitempty"`
}

type SubInfo struct {
	gorm.Model
	Eid        int64     `gorm:"column:eid"`
	Sid        string    `gorm:"column:sid"`
	Pid        int64     `gorm:"column:pid"`
	State      int32     `gorm:"column:state"`
	FinishedAt time.Time `gorm:"column:finished_at"`
}
