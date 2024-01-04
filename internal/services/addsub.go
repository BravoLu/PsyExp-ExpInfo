package services


import (
	"context"
	"time"

	pb "github.com/BravoLu/grpc_idl"
	"github.com/grpc_exp_info/internal/entity"
	"github.com/grpc_exp_info/internal/mysql"
	"github.com/grpc_plugins/log"
)


func (s *ExpInfoServerImpl) AddSub(
	ctx context.Context,
	req *pb.AddSubReq,
) (*pb.AddSubRsp, error) {
	log.Infof("req: %+v", req)

	// TODO: check if the eid is valid

	e := &entity.SubInfo{
		Eid: req.Eid,
		Pid: req.Pid,
		State: 1,
		FinishedAt: time.Now().Add(24 * time.Hour),  // Mock
	}

	dao := &mysql.ExpInfoDaoImpl{}
	sid, err := dao.AddSub(ctx, e)
	if err != nil {
		return nil, err
	}
	rsp := &pb.AddSubRsp{
		Code: 0,
		Msg: "ok",
		Sid: sid, 
	}
	return rsp, nil
}