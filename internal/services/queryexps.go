package services

import (
	"context"
	"encoding/json"

	pb "github.com/BravoLu/grpc_idl"
	"github.com/grpc_exp_info/internal/mysql"
	"github.com/grpc_plugins/log"
)

func (s *ExpInfoServerImpl) QueryExps(
	ctx context.Context,
	req *pb.QueryExpsReq,
) (*pb.QueryExpsRsp, error) {
	log.Infof("QueryExps req: %+v", req)

	qry := &mysql.QueryCond{
		Rid:       req.Rid,
		Offset:    int((req.PageIndex - 1) * req.PageSize),
		Limit:     int(req.PageSize),
		OrderType: req.OrderType,
		State:     int32(req.State),
		// MinPrice:      req.MinPrice,
		// OnlySeeMyself: req.OnlySeeMyself,
	}
	// DB查數據
	dao := &mysql.ExpInfoDaoImpl{}
	res, num, stats, err := dao.QueryExps(ctx, qry)
	if err != nil {
		return nil, err
	}

	// 構造回包數據
	var exps []*pb.ExpInfo
	for _, v := range res {
		tags := make([]string, 0)
		err := json.Unmarshal(v.Tags, &tags)
		if err == nil {
			log.Infof("tags: ", tags)
		}
		if err != nil {
			log.Errorf("convert tags error: %+v", v.Tags)
			// continue
			tags = make([]string, 0)
		}

		t := &pb.ExpInfo{
			Eid:        int64(v.ID),
			Title:      v.Title,
			Desc:       v.Desc,
			Rid:        v.Rid,
			Ctime:      v.Ctime,
			Pnum:       v.Pnum,
			CreateTime: v.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdateTime: v.UpdatedAt.Format("2006-01-02 15:04:05"),
			Price:      v.Price,
			CurType:    v.CurType,
			Url:        v.Url,
			Tags:       tags,
			State:      pb.ExpState(v.State),
			Deadline: v.EndTime.Format("2006-01-02 15:04:05"),
		}
		exps = append(exps, t)
	}
	resp := &pb.QueryExpsRsp{
		Code:     0,
		Msg:      "ok",
		TotalNum: int32(num),
		Exps:     exps,
		Stats: &pb.ExpStats{
			OngoingNum:  int32(stats.OnGoingNum),
			ClosedNum:   int32(stats.ClosedNum),
			AllNum:      int32(stats.AllNum),
			FinishedNum: int32(stats.FinishedNum),
		},
	}
	return resp, err
}
