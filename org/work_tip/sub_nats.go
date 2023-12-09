package OrgWorkTip

import (
	"encoding/json"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreNats "gitee.com/weeekj/weeekj_core/v5/core/nats"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//处理新增通知
	CoreNats.SubDataByteNoErr("/org/work_tip", subNatsNew)
}

func subNatsNew(_ *nats.Msg, action string, _ int64, _ string, data []byte) {
	//跳出非添加数据
	if action != "new" {
		return
	}
	//获取参数
	var args argsAppendTip
	if err := json.Unmarshal(data, &args); err != nil {
		CoreLog.Error("org work tip, get params, ", err)
		return
	}
	//添加数据
	if err := appendTip(args); err != nil {
		CoreLog.Error("org work tip, insert, ", err)
		return
	}
}
