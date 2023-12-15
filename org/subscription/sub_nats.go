package OrgSubscription

import (
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	OrgCore "github.com/fotomxq/weeekj_core/v5/org/core"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//标记会员续费
	CoreNats.SubDataByteNoErr("/org/sub/set_add", subNatsSubSetAdd)
	//通知过期数据包
	CoreNats.SubDataByteNoErr("/base/expire_tip/expire", subNatsExpireTip)
}

// 标记会员续费
func subNatsSubSetAdd(_ *nats.Msg, _ string, _ int64, _ string, data []byte) {
	//日志
	logAppend := "org sub sub nats order pay, "
	//解析数据包
	var args argsSetSubAdd
	if err := CoreNats.ReflectDataByte(data, &args); err != nil {
		CoreLog.Error(logAppend, "get data, ", err)
		return
	}
	//设置订阅
	if err := setSubAdd(&args); err != nil {
		CoreLog.Error(logAppend, "set sub by order, err: ", err)
		return
	}
}

// 通知过期数据包
func subNatsExpireTip(_ *nats.Msg, action string, id int64, mark string, _ []byte) {
	//如果系统不符合，跳出
	if action != "org_sub" {
		return
	}
	//日志
	logAppend := "user sub sub nats expire tip, "
	//找到绑定ID
	subData, err := getSub(id)
	if err != nil {
		CoreLog.Error(logAppend, "get sub data, ", err)
		return
	}
	//检查过期时间是否一致？
	if getSubHash(&subData) != mark {
		CoreLog.Error(logAppend, "sub no hash")
		return
	}
	//收尾工作
	configData, err := GetConfigByID(&ArgsGetConfigByID{
		ID: subData.ConfigID,
	})
	if err == nil {
		orgData, err := OrgCore.GetOrg(&OrgCore.ArgsGetOrg{
			ID: subData.OrgID,
		})
		if err == nil {
			//处理商户的权益
			expireSubLast(&configData, &orgData)
		}
	}
	//删除过期数据
	if err = DeleteSub(&ArgsDeleteSub{
		ID: id,
	}); err != nil {
		CoreLog.Error(logAppend, "clear sub by user, ", err)
		return
	}
}
