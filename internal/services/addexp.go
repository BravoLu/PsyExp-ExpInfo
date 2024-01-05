package services

import (
	"context"
	"encoding/json"
	"gorm.io/datatypes"
	"time"

	pb "github.com/BravoLu/grpc_idl"
	"github.com/grpc_exp_info/internal/entity"
	"github.com/grpc_exp_info/internal/mysql"
	"github.com/grpc_plugins/log"
)

func (s *ExpInfoServerImpl) AddExp(
	ctx context.Context,
	req *pb.AddExpReq,
) (*pb.AddExpRsp, error) {
	log.Infof("req: %+v", req)
	var tags []string
	for _, t := range req.Tags {
		tags = append(tags, t)
	}

	jsonTags, err := json.Marshal(tags)
	if err != nil {
		panic(err)
	}
	var et time.Time
	if len(req.Deadline) != 0 {
		et, err = time.Parse("2006-01-02 15:04:05", req.Deadline)
		if err != nil {
			log.Errorf("parse time error: %+v", err)
			return nil, err
		}
	}

	e := &entity.ExpInfo{
		Title:   req.Title,
		Desc:    req.Desc,
		Rid:     req.Rid,
		Ctime:   req.Ctime,
		Pnum:    req.Pnum,
		EndTime: et,
		CurType: req.CurType,
		Price:   req.Price,
		Url:     req.Url,
		State:   1,
		Tags:    datatypes.JSON(jsonTags),
	}
	dao := &mysql.ExpInfoDaoImpl{}
	id, err := dao.AddExp(ctx, e)
	if err != nil {
		return nil, err
	}
	rsp := &pb.AddExpRsp{
		Code: 0,
		Msg:  "ok",
		Eid:  id,
	}
	return rsp, nil
}
