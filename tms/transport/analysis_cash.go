package TMSTransport

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetAnalysisCashSum 获取配送员
type ArgsGetAnalysisCashSum struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//配送人员
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//收付款类型
	// 0 收款(配送员收到款项) 1 付款(配送员付出款项)
	PayType int `db:"pay_type" json:"payType"`
	//时间范围
	// 部分统计支持
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
}

func GetAnalysisCashSum(args *ArgsGetAnalysisCashSum) (count int64, err error) {
	var timeBetween CoreSQLTime.FieldsCoreTime
	timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
	if err != nil {
		return
	}
	count, err = CoreSQL.GetAllSumMap(Router2SystemConfig.MainDB.DB, "tms_transport_cash", "price", "org_id = :org_id AND (:bind_id < 1 OR bind_id = :bind_id) AND create_at >= :start_at AND create_at <= :end_at AND (:pay_type < 0 OR pay_type = :pay_type)", map[string]interface{}{
		"org_id":   args.OrgID,
		"bind_id":  args.BindID,
		"pay_type": args.PayType,
		"start_at": timeBetween.MinTime,
		"end_at":   timeBetween.MaxTime,
	})
	return
}
