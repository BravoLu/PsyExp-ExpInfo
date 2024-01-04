package services

import (
	"context"
	"encoding/json"

	pb "github.com/BravoLu/grpc_idl"
	"github.com/grpc_exp_info/internal/mysql"
	"github.com/grpc_plugins/log"
)

func (s *ExpInfoServerImpl) QueryExp(
	ctx context.Context,
	req *pb.QueryExpReq,
) (*pb.QueryExpRsp, error) {
	log.Info("QueryExp: ", req)
	defer func() {
	}()
	// DB查數據
	dao := &mysql.ExpInfoDaoImpl{}
	res, err := dao.QueryExp(ctx, req.Eid)
	if err != nil {
		return nil, err
	}

	eSubs, subsCnt, fCnt, err := dao.QuerySubs(ctx, req.Eid)
	if err != nil {
		return nil, err
	}
	var subs []*pb.SubInfo
	for _, sub := range eSubs {
		t := &pb.SubInfo{
			Eid: sub.Eid,
			Pid: sub.Pid,
			Sid: sub.Sid,
			FinishedAt: sub.FinishedAt.Format("2006-01-02 15:04:05"),
			CreateTime: sub.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdateTime: sub.UpdatedAt.Format("2006-01-02 15:04:05"),
			State: pb.SubState(sub.State),
			
		}
		subs = append(subs, t)
	}
	log.Infof("subs: ", subs)

	var tags []string
	err = json.Unmarshal(res.Tags, &tags)
	if err != nil {
		panic(err)
	}

	resp := &pb.QueryExpRsp{
		Code: 0,
		Msg: "ok",
		Exp: &pb.ExpInfo{
			Eid: int64(res.ID),
			Title: res.Title,
			Desc: res.Desc,
			Rid: res.Rid,
			Ctime: res.Ctime,
			Pnum: res.Pnum,
			Price: res.Price,
			CreateTime: res.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdateTime: res.UpdatedAt.Format("2006-01-02 15:04:05"),
			State: pb.ExpState(res.State),
			Url: res.Url,
			Tags: tags,
			Deadline: res.EndTime.Format("2006-01-02 15:04:05"),
		},
		Subs: subs,
		SubsNum: int32(subsCnt),
		FinishedNum: int32(fCnt), 
	}
	return resp, err
}