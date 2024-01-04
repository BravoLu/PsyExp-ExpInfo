package mysql

import (
	"context"
	"fmt"

	"github.com/grpc_exp_info/internal/entity"
	"github.com/grpc_plugins/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	masterDB *gorm.DB
	slaveDB  *gorm.DB
)

const (
	Database     = "PsyExp"
	ExpTableName = "exp_infos"
	SubTableName = "sub_infos"
)

func MasterClient() (*gorm.DB, error) {
	if masterDB == nil {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.Config.Db.Master.User, config.Config.Db.Master.Passwd,
			config.Config.Db.Master.IP, config.Config.Db.Master.Port, Database)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, err
		}
		masterDB = db
	}
	return masterDB, nil
}

func SlaveClient() (*gorm.DB, error) {
	if slaveDB == nil {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.Config.Db.Slave.User, config.Config.Db.Slave.Passwd,
			config.Config.Db.Slave.IP, config.Config.Db.Slave.Port, Database)
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, err
		}
		slaveDB = db
	}
	return slaveDB, nil
}

type QueryCond struct {
	Rid           int64
	Limit         int
	Offset        int
	State         int32
	OrderType     int32
	OnlySeeMyself int32
}

type Stats struct {
	OnGoingNum  int64
	AllNum      int64
	FinishedNum int64
	ClosedNum   int64
}

type ExpInfoDao interface {
	AddExp(context.Context, *entity.ExpInfo) (int64, error)
	UpdateExp(context.Context, *entity.ExpInfo) (error)
	QueryExp(context.Context, int64) (*entity.ExpInfo, error)
	QueryExps(context.Context, *QueryCond) ([]*entity.ExpInfo, int64, *Stats, error)

	AddSub(context.Context, *entity.SubInfo) (string, error)
	UpdateSub(context.Context, *entity.SubInfo) (string, error)
	QuerySub(context.Context, *entity.SubInfo) (*entity.SubInfo, error)
	QuerySubs(context.Context, int64) ([]*entity.SubInfo, int32, error)
	QueryUserSubs(context.Context, *QueryCond) ([]*entity.SubInfo, int32, error)
	GetSnum(context.Context, int64) (int64, error)
	GetFinishedNum(context.Context, int64) (int64, error)
}
