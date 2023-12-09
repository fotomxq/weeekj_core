package TMSTransport

import (
	CoreSQLTime "gitee.com/weeekj/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetAnalysisBindGoods 反馈指定配送员的统计数据参数
type ArgsGetAnalysisBindGoods struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//多个成员ID
	BindIDs pq.Int64Array `db:"bind_ids" json:"bindIDs" check:"ids"`
	//是否为退单
	NeedIsRefund bool `json:"needIsRefund"`
	IsRefund     bool `json:"isRefund"`
	//查询时间范围
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
}

// GetAnalysisBindGoods 反馈指定配送员的统计数据
func GetAnalysisBindGoods(args *ArgsGetAnalysisBindGoods) (dataList []FieldsBindAnalysisGoods, err error) {
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, org_id, bind_id, goods FROM tms_transport_bind_analysis_goods WHERE org_id = $1 AND bind_id = ANY($2) AND create_at >= $3 AND create_at <= $4 AND ($5 = false OR is_refund = $6)", args.OrgID, args.BindIDs, timeBetween.MinTime, timeBetween.MaxTime, args.NeedIsRefund, args.IsRefund)
	if err != nil {
		return
	}
	//合并数据
	var newDataList []FieldsBindAnalysisGoods
	for _, v := range dataList {
		isFind := false
		for k2, v2 := range newDataList {
			if v.BindID == v2.BindID {
				if v.CreateAt.Unix() > v2.CreateAt.Unix() {
					newDataList[k2].CreateAt = v.CreateAt
				}
				if len(newDataList[k2].Goods) < 1 {
					newDataList[k2].Goods = FieldsBindAnalysisGoodsGoods{}
				}
				for _, v3 := range v.Goods {
					isFindGood := false
					for k4, v4 := range v2.Goods {
						if v3.System == v4.System && v3.ID == v4.ID {
							newDataList[k2].Goods[k4].Count += v3.Count
							isFindGood = true
							break
						}
					}
					if !isFindGood {
						newDataList[k2].Goods = append(newDataList[k2].Goods, v3)
					}
				}
				isFind = true
				break
			}
		}
		if !isFind {
			newDataList = append(newDataList, v)
		}
	}
	dataList = newDataList
	//反馈
	return
}
