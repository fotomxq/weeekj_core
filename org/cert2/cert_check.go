package OrgCert2

import (
	BaseExpireTip "github.com/fotomxq/weeekj_core/v5/base/expire_tip"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"time"
)

// ArgsCheckCertByMarks 通过一组mark查询证件到期情况参数
type ArgsCheckCertByMarks struct {
	//配置组
	Marks pq.StringArray `db:"marks" json:"marks" check:"marks"`
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//绑定ID
	// 根据配置决定，可能是用户ID\组织ID\或其他任意主体
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
}

// DataCheckCertByMarks 通过一组mark查询证件到期情况
type DataCheckCertByMarks struct {
	Mark string `json:"mark"`
	IsOK bool   `json:"isOK"`
	SN   string `json:"sn"`
}

// CheckCertByMarks 通过一组mark查询证件到期情况
func CheckCertByMarks(args *ArgsCheckCertByMarks) (dataList []DataCheckCertByMarks, err error) {
	var configList []FieldsConfig
	err = Router2SystemConfig.MainDB.Select(&configList, "SELECT id, mark FROM org_cert_config2 WHERE mark = ANY($1) AND delete_at < to_timestamp(1000000)", args.Marks)
	if err != nil {
		return
	}
	var configIDs pq.Int64Array
	for _, v := range configList {
		configIDs = append(configIDs, v.ID)
	}
	var certList []FieldsCert
	err = Router2SystemConfig.MainDB.Select(&certList, "SELECT id, config_id, sn FROM org_cert2 WHERE bind_id = $1 AND config_id = ANY($2) AND ($3 < 1 OR org_id = $3) AND delete_at < to_timestamp(1000000) AND expire_at >= NOW() AND audit_at >= to_timestamp(1000000)", args.BindID, configIDs, args.OrgID)
	if err != nil {
		return
	}
	for _, v := range configList {
		isFind := false
		for _, v2 := range certList {
			if v.ID == v2.ConfigID {
				isFind = true
				dataList = append(dataList, DataCheckCertByMarks{
					Mark: v.Mark,
					IsOK: true,
					SN:   v2.SN,
				})
				break
			}
		}
		if !isFind {
			dataList = append(dataList, DataCheckCertByMarks{
				Mark: v.Mark,
				IsOK: false,
				SN:   "",
			})
		}
	}
	return
}

// 处理证件过期提醒
func checkCertExpireSend(certID int64, expireAt time.Time) {
	if !CoreSQL.CheckTimeHaveData(expireAt) {
		return
	}
	//提前30天过期提醒
	BaseExpireTip.AppendTipNoErr(&BaseExpireTip.ArgsAppendTip{
		OrgID:      0,
		UserID:     0,
		SystemMark: "org_cert",
		BindID:     certID,
		Hash:       "",
		ExpireAt:   CoreFilter.GetCarbonByTime(expireAt).SubDays(30).Time,
	})
	//普通过期提醒
	BaseExpireTip.AppendTipNoErr(&BaseExpireTip.ArgsAppendTip{
		OrgID:      0,
		UserID:     0,
		SystemMark: "org_cert",
		BindID:     certID,
		Hash:       "",
		ExpireAt:   expireAt,
	})
}
