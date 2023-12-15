package AnalysisOrg

import (
	TMSTransport "github.com/fotomxq/weeekj_core/v5/tms/transport"
)

// ArgsGetBindMarge 获取聚合统计参数
type ArgsGetBindMarge struct {
	//多个标识码组
	Marks []ArgsMargeMark `json:"marks"`
	//商户ID
	OrgID int64 `json:"orgID"`
	//成员ID
	BindID int64 `json:"bindID"`
}

// GetBindMarge 获取聚合统计
func GetBindMarge(args *ArgsGetBindMarge) (result DataMarge, err error) {
	//遍历聚合函数处理
	for _, vMark := range args.Marks {
		if vMark.Limit < 1 || vMark.Limit > 1000 {
			continue
		}
		vResult := DataMargeMark{
			Mark:        vMark.Mark,
			TimeBetween: vMark.TimeBetween,
			Limit:       vMark.Limit,
		}
		switch vMark.Mark {
		case "tms_wait_count":
			//尚未完成的配送任务
			vResult.Data, err = TMSTransport.GetAnalysisBindWaitCount(&TMSTransport.ArgsGetAnalysisBindWaitCount{
				OrgID:       args.OrgID,
				BindID:      args.BindID,
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
				vResult.Data = 0
			}
		case "tms_finish_count":
			//已完成配送任务
			vResult.Data, err = TMSTransport.GetAnalysisBindFinishCount(&TMSTransport.ArgsGetAnalysisBindWaitCount{
				OrgID:       args.OrgID,
				BindID:      args.BindID,
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
				vResult.Data = 0
			}
		case "tms_all_count":
			//全部配送任务总数
			vResult.Data, err = TMSTransport.GetAnalysisBindCount(&TMSTransport.ArgsGetAnalysisBindCount{
				OrgID:       args.OrgID,
				BindID:      args.BindID,
				TimeBetween: vMark.TimeBetween,
			})
			if err != nil {
				err = nil
				vResult.Data = 0
			}
		}
		result.Marks = append(result.Marks, vResult)
	}
	return
}
