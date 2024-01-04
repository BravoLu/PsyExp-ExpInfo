package services

import (
	"context"

	pb "github.com/BravoLu/grpc_idl"
	"github.com/grpc_exp_info/internal/entity"
	"github.com/grpc_exp_info/internal/mysql"
	"github.com/grpc_plugins/log"
)

func (s *ExpInfoServerImpl) UpdateSub(
	ctx context.Context,
	req *pb.UpdateSubReq,
) (*pb.UpdateSubRsp, error) {
	log.Infof("req: ", req)

	dao := &mysql.ExpInfoDaoImpl{}

	// check
	sub, err := dao.QuerySub(ctx, &entity.SubInfo{Sid: req.Sid})
	if err != nil {
		return nil, err
	}

	u := &entity.SubInfo{
		Sid:   req.Sid,
		State: int32(req.State),
	}

	sid, err := dao.UpdateSub(ctx, u)
	if err != nil {
		return nil, err
	}

	// Check The whether the snum of exp reach the snum -> change state .
	exp, err := dao.QueryExp(ctx, sub.Eid)
	if err != nil {
		return nil, err
	}

	fnum, err := dao.GetFinishedNum(ctx, sub.Eid)
	if err != nil {
		return nil, err
	}
	log.Infof("pnum:[%+v], snum:[%+v]", exp.Pnum, fnum)

	if exp.Pnum == int32(fnum) {
		// update exp
		log.Infof("update state of [%s]", exp.ID)
		exp.State = 3
		err := dao.UpdateExp(ctx, int64(exp.ID), exp)
		if err != nil {
			return nil, err
		}
	}

	rsp := &pb.UpdateSubRsp{
		Code: 0,
		Msg:  "ok",
		Sid:  sid,
	}

	return rsp, nil
}
