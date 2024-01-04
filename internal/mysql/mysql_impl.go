package mysql

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/grpc_exp_info/internal/entity"
	"github.com/grpc_plugins/log"
)

type ExpInfoDaoImpl struct {
}

func (s *ExpInfoDaoImpl) AddExp(
	ctx context.Context,
	e *entity.ExpInfo,
) (int64, error) {
	tx, err := MasterClient()
	if err != nil {
		return 0, err
	}
	if err := tx.WithContext(ctx).
		Table(ExpTableName).
		Debug().
		Create(e).Error; err != nil {
		log.Errorf("AddExp error: %+v", err)
		return 0, err
	}
	return int64(e.ID), nil
}

func (s *ExpInfoDaoImpl) UpdateExp(
	ctx context.Context,
	id int64,
	e *entity.ExpInfo,
) error {
	tx, err := MasterClient()
	if err != nil {
		return err
	}
	tx = tx.Table(ExpTableName).
		WithContext(ctx).
		Where("id = ?", id)
	if e.Tags.String() == "null" {
		tx = tx.Omit("tags")
	}
	if err := tx.Debug().Updates(e).Error; err != nil {
		return err
	}
	return nil
}

func (s *ExpInfoDaoImpl) QueryExp(
	ctx context.Context,
	eid int64,
) (*entity.ExpInfo, error) {
	tx, err := SlaveClient()
	if err != nil {
		return nil, err
	}
	var res *entity.ExpInfo
	if err := tx.WithContext(ctx).Table(ExpTableName).
		Debug().
		Where("id = ?", eid).
		Find(&res).Error; err != nil {
		log.Errorf("query db error: %+v", err)
		return nil, err
	}

	return res, nil
}

func (s *ExpInfoDaoImpl) QueryExps(
	ctx context.Context,
	cond *QueryCond,
) ([]*entity.ExpInfo, int64, *Stats, error) {
	tx, err := SlaveClient()
	if err != nil {
		return nil, 0, nil, err
	}
	tx = tx.Table(ExpTableName).Debug().Where("state != 0")

	if cond.Rid != 0 {
		tx = tx.Where("rid = ?", cond.Rid)
	}

	// get the stats
	var onGoingNum, allNum, finishedNum, closedNum int64
	if err := tx.WithContext(ctx).
		Count(&allNum).Error; err != nil {
		log.Errorf("FindExperiments get count error: %+v", err)
		return nil, 0, nil, err
	}

	if err := tx.WithContext(ctx).
		Where("state = 1").
		Count(&onGoingNum).Error; err != nil {
		log.Errorf("FindExperiments get count error: %+v", err)
		return nil, 0, nil, err
	}

	if err := tx.WithContext(ctx).
		Where("state = 2").
		Count(&finishedNum).Error; err != nil {
		log.Errorf("FindExperiments get count error: %+v", err)
		return nil, 0, nil, err
	}

	if err := tx.WithContext(ctx).
		Where("state = 3").
		Count(&closedNum).Error; err != nil {
		log.Errorf("FindExperiments get count error: %+v", err)
		return nil, 0, nil, err
	}

	stats := &Stats{
		OnGoingNum:  onGoingNum,
		ClosedNum:   closedNum,
		AllNum:      allNum,
		FinishedNum: finishedNum,
	}

	if cond.State != 0 {
		tx = tx.Where("state = ?", cond.State)
	}
	// TODO: 新增条件在这里加
	var cnt int64
	if err := tx.WithContext(ctx).Count(&cnt).Error; err != nil {
		log.Errorf("FindExperiments get count error: %+v", err)
		return nil, 0, nil, err
	}
	if cond.Limit == 0 {
		tx = tx.Limit(10)
	} else {
		tx = tx.Limit(cond.Limit)
	}
	if cond.Offset != 0 {
		tx = tx.Offset(cond.Offset)
	}

	// Order Type
	if cond.OrderType == 1 {
		tx = tx.Order("price DESC")
	} else {
		tx = tx.Order("updated_at DESC")
	}
	var res []*entity.ExpInfo
	if err := tx.WithContext(ctx).
		Find(&res).Error; err != nil {
		log.Errorf("db error: %+v", err)
		return nil, 0, nil, err
	}

	return res, cnt, stats, nil
}

func (s *ExpInfoDaoImpl) AddSub(
	ctx context.Context,
	e *entity.SubInfo,
) (string, error) {
	tx, err := MasterClient()
	if err != nil {
		return "", err
	}
	uuid, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	e.Sid = uuid.String()
	if err := tx.WithContext(ctx).
		Table(SubTableName).
		Debug().
		Create(e).Error; err != nil {
		log.Errorf("AddExp error: %+v", err)
		return "", err
	}
	return e.Sid, nil
}

func (s *ExpInfoDaoImpl) UpdateSub(
	ctx context.Context,
	sub *entity.SubInfo,
) (string, error) {
	tx, err := MasterClient()
	if err != nil {
		return "", err
	}

	if err := tx.WithContext(ctx).
		Where("sid = ?", sub.Sid).
		Omit("created_at", "researcher_id").
		Updates(sub).Error; err != nil {
		return "", err
	}

	return sub.Sid, nil
}

func (s *ExpInfoDaoImpl) QuerySubs(
	ctx context.Context,
	eid int64,
) ([]*entity.SubInfo, int64, int64, error) {
	tx, err := SlaveClient()
	if err != nil {
		return nil, 0, 0, err
	}

	var res []*entity.SubInfo

	tx = tx.WithContext(ctx).
		Table(SubTableName).
		Where("eid = ?", eid)

	if err := tx.
		Order("updated_at DESC").
		Debug().
		Find(&res).Error; err != nil {
		return nil, 0, 0, err
	}

	// get cnt
	cnt, err := s.GetSnum(ctx, eid)
	if err != nil {
		return nil, 0, 0, err
	}

	// get finished num
	fcnt, err := s.GetFinishedNum(ctx, eid)
	if err != nil {
		return nil, 0, 0, err
	}

	return res, cnt, fcnt, nil
}

func (s *ExpInfoDaoImpl) QueryUserSubs(
	ctx context.Context,
	cond *QueryCond,
) ([]*entity.SubInfo, int32, error) {
	tx, err := SlaveClient()
	if err != nil {
		return nil, 0, err
	}

	var res []*entity.SubInfo
	tx = tx.WithContext(ctx).
		Table(SubTableName).
		Debug()

	if cond.State != 0 {
		tx = tx.Where("state = ?", cond.State)
	}
	var cnt int64

	if err := tx.Count(&cnt).Error; err != nil {
		return nil, 0, err
	}

	if cond.Limit == 0 {
		tx = tx.Limit(10)
	} else {
		tx = tx.Limit(cond.Limit)
	}
	if cond.Offset != 0 {
		tx = tx.Offset(cond.Offset)
	}

	if err := tx.
		Order("updated_at DESC").
		Where("pid = ?", cond.Rid).
		Find(&res).Error; err != nil {
		return nil, 0, err
	}

	return res, int32(cnt), nil
}

func (s *ExpInfoDaoImpl) GetSnum(
	ctx context.Context,
	eid int64,
) (int64, error) {
	tx, err := SlaveClient()
	if err != nil {
		return 0, err
	}
	var cnt int64

	if err := tx.WithContext(ctx).
		Table(SubTableName).
		Where("eid = ?", eid).Debug().Count(&cnt).Error; err != nil {
		return 0, err
	}
	return cnt, nil
}

func (s *ExpInfoDaoImpl) GetFinishedNum(
	ctx context.Context,
	eid int64,
) (int64, error) {
	tx, err := SlaveClient()
	if err != nil {
		return 0, err
	}
	var cnt int64

	if err := tx.WithContext(ctx).
		Table(SubTableName).
		Where("eid = ? and state = 2", eid).Debug().Count(&cnt).Error; err != nil {
		return 0, err
	}
	return cnt, nil
}

func (s *ExpInfoDaoImpl) QuerySub(ctx context.Context,
	sub *entity.SubInfo,
) (*entity.SubInfo, error) {
	tx, err := SlaveClient()
	if err != nil {
		return nil, err
	}
	var res *entity.SubInfo
	tx = tx.WithContext(ctx).
		Table(SubTableName).
		Debug()
	if len(sub.Sid) != 0 {
		tx = tx.Where("sid = ?", sub.Sid)
	}
	if sub.Pid != 0 {
		tx = tx.Where("pid = ?", sub.Pid)
	}
	if sub.Eid != 0 {
		tx = tx.Where("eid = ?", sub.Eid)
	}
	sql := tx.Find(&res)
	if err := sql.Error; err != nil {
		return nil, err
	}
	if sql.RowsAffected == 0 {
		return nil, fmt.Errorf("subject not found.")
	}
	return res, nil
}
