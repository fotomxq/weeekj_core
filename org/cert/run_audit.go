package OrgCert

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

func runAudit() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("org cert audit run, ", r)
		}
	}()
	//遍历数据
	limit := 100
	step := 0
	for {
		var dataList []FieldsCert
		if err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id FROM org_cert WHERE delete_at < to_timestamp(1000000) AND pay_at > to_timestamp(1000000) AND audit_at < to_timestamp(1000000) AND audit_ban_at < to_timestamp(1000000) ORDER BY id LIMIT $1 OFFSET $2", limit, step); err != nil {
			break
		}
		if len(dataList) < 1 {
			break
		}
		for _, v := range dataList {
			if err := UpdateCertAudit(&ArgsUpdateCertAudit{
				ID:          v.ID,
				OrgID:       -1,
				ChildOrgID:  -1,
				AuditBindID: 0,
				IsBan:       false,
				AuditDes:    "系统自动审核",
			}); err != nil {
				CoreLog.Error("org cert audit run, auto audit, ", err)
			}
		}
		//下一页继续
		step += limit
	}
}
