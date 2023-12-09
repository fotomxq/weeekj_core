package OrgCert

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// 支付状态检查维护
func runPay() {
	//捕捉异常
	defer func() {
		if r := recover(); r != nil {
			CoreLog.Error("org cert pay run, ", r)
		}
	}()
	//遍历数据
	limit := 100
	step := 0
	for {
		var dataList []FieldsCert
		if err := Router2SystemConfig.MainDB.Select(&dataList, "SELECT id FROM org_cert WHERE pay_failed = false AND pay_id > 0 AND delete_at < to_timestamp(1000000) ORDER BY id LIMIT $1 OFFSET $2", limit, step); err != nil {
			break
		}
		if len(dataList) < 1 {
			break
		}
		for _, v := range dataList {
			_, err := CheckPay(&ArgsCheckPay{
				ID:     v.ID,
				OrgID:  -1,
				BindID: -1,
			})
			if err != nil {
				CoreLog.Error("org cert pay run, update pay status, id: ", v.ID, ", err: ", err)
			}
		}
		//下一页继续
		step += limit
	}
}
