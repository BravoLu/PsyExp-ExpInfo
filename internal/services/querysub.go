package services

import (
	"context"

	pb "github.com/BravoLu/grpc_idl"
	"github.com/grpc_exp_info/internal/entity"
	"github.com/grpc_exp_info/internal/mysql"
	"github.com/grpc_plugins/log"
)

func (s *ExpInfoServerImpl) QuerySub(
	ctx context.Context,
	req *pb.QuerySubReq,
) (*pb.QuerySubRsp, error) {
	log.Infof("req: %+v", req)

	dao := &mysql.ExpInfoDaoImpl{}
	_, err := dao.QuerySub(ctx, &entity.SubInfo{
		Pid: req.Pid,
		Sid: req.Sid,
		Eid: req.Eid,
	})
	if err != nil {
		return &pb.QuerySubRsp{Code: 1, Msg: "sub doesn't exist"}, nil
	}

	return &pb.QuerySubRsp{Code: 0, Msg: "success"}, nil
}
