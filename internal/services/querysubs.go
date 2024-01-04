package services

import (
	"context"

	pb "github.com/BravoLu/grpc_idl"
	"github.com/grpc_exp_info/internal/mysql"
	"github.com/grpc_plugins/log"
)

func (s *ExpInfoServerImpl) QuerySubs(
	ctx context.Context,
	req *pb.QuerySubsReq,
) (*pb.QuerySubsRsp, error) {
	log.Infof("req: ", req)
	dao := &mysql.ExpInfoDaoImpl{}
	cond := &mysql.QueryCond{
		Rid:    req.Pid,
		Offset: int((req.PageIndex - 1) * req.PageSize),
		Limit:  int(req.PageSize),
		State: int32(req.State),
	}
	subs, cnt, err := dao.QueryUserSubs(ctx, cond)
	if err != nil {
		return nil, err
	}
	var res []*pb.SubInfo
	for _, v := range subs {
		subInfo := &pb.SubInfo{
			Sid:        v.Sid,
			Eid:        v.Eid,
			Pid:        v.Pid,
			FinishedAt: v.FinishedAt.Format("2006-01-02 15:04:05"),
			CreateTime: v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdateTime: v.UpdatedAt.Format("2006-01-02 15:04:05"),
			State:      pb.SubState(v.State),
		}
		res = append(res, subInfo)
	}
	return &pb.QuerySubsRsp{Code: 0, Msg: "success", Subs: res, TotalNum: cnt}, nil
}
