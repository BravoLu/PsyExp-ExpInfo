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

func (s *ExpInfoServerImpl) UpdateExp(
	ctx context.Context,
	req *pb.UpdateExpReq,
) (*pb.UpdateExpRsp, error) {
	log.Infof("req: %+v", req)
	dao := &mysql.ExpInfoDaoImpl{}
	// precheck ~
	d, err := dao.QueryExp(ctx, req.Eid)
	if err != nil {
		return &pb.UpdateExpRsp{Code: 2, Msg: "invalid experiment id."}, nil
	}
	if d.Rid != req.Rid {
		return &pb.UpdateExpRsp{Code: 3, Msg: "invalid operation."}, nil
	}
	// publish
	if req.State == pb.ExpState_Published {
		if d.State != 1 {
			return &pb.UpdateExpRsp{Code: 4, Msg: "invalid state."}, nil
		}
	}

	if req.State == pb.ExpState_Closed {
		if d.State != 1 && d.State != 2 {
			return &pb.UpdateExpRsp{Code: 4, Msg: "invalid state."}, nil
		}
	}

	// close
	jsonTags, err := json.Marshal(req.Tags)
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
		// ID: req.Eid,
		Rid:     req.Rid,
		Title:   req.Title,
		Desc:    req.Desc,
		Tags:    datatypes.JSON(jsonTags),
		Url:     req.Url,
		Pnum:    req.Pnum,
		Ctime:   req.Ctime,
		Price:   req.Price,
		EndTime: et,
		State:   int32(req.State),
	}
	log.Infof("tags: %+v", e.Tags.String())
	// if e.Tags.IsZero() != 0 {
	// 	e.Tags = datatypes.JSON(jsonTags)
	// }
	err = dao.UpdateExp(ctx, req.Eid, e)
	if err != nil {
		return nil, err
	}
	return &pb.UpdateExpRsp{Code: 0, Msg: "ok"}, nil
}
