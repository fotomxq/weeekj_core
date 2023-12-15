package UserSubscription

import (
	BaseExpireTip "github.com/fotomxq/weeekj_core/v5/base/expire_tip"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	OrgUserMod "github.com/fotomxq/weeekj_core/v5/org/user/mod"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//标记会员续费
	CoreNats.SubDataByteNoErr("/user/sub/set_add", subNatsSubSetAdd)
	//设置订阅
	CoreNats.SubDataByteNoErr("/user/sub/set", subNatsSubSet)
	//使用订阅
	CoreNats.SubDataByteNoErr("/user/sub/use", subNatsSubUse)
	//通知过期数据包
	CoreNats.SubDataByteNoErr("/base/expire_tip/expire", subNatsExpireTip)
}

// 标记会员续费
func subNatsSubSetAdd(_ *nats.Msg, _ string, _ int64, _ string, data []byte) {
	//日志
	logAppend := "user sub sub nats order pay, "
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

// 设置订阅
func subNatsSubSet(_ *nats.Msg, _ string, _ int64, _ string, data []byte) {
	logAppend := "user sub sub nats sub set, "
	//通过args解析数据
	var args ArgsSetSub
	if err := CoreNats.ReflectDataByte(data, &args); err != nil {
		CoreLog.Error(logAppend, "get args, ", err)
		return
	}
	//设置数据
	if err := SetSub(&args); err != nil {
		CoreLog.Error(logAppend, "set sub, ", err)
		return
	}
}

// 使用订阅
func subNatsSubUse(_ *nats.Msg, _ string, _ int64, _ string, data []byte) {
	logAppend := "user sub sub nats sub use, "
	//通过args解析数据
	var args ArgsUseSub
	if err := CoreNats.ReflectDataByte(data, &args); err != nil {
		CoreLog.Error(logAppend, "get args, ", err)
		return
	}
	//设置数据
	if err := UseSub(&args); err != nil {
		CoreLog.Error(logAppend, "use sub, ", err)
		return
	}
}

// 通知过期数据包
func subNatsExpireTip(_ *nats.Msg, action string, id int64, mark string, data []byte) {
	//如果系统不符合，跳出
	if action != "user_sub" {
		return
	}
	//日志
	logAppend := "user sub sub nats expire tip, "
	//解析数据
	rawData, err := BaseExpireTip.GetExpireData(data)
	if err != nil {
		CoreLog.Error(logAppend, "get expire data, ", err)
		return
	}
	//找到绑定ID
	subData, err := GetSub(&ArgsGetSub{
		ConfigID: id,
		UserID:   rawData.UserID,
	})
	if err != nil {
		CoreLog.Error(logAppend, "get sub data, ", err)
		return
	}
	//检查过期时间是否一致？
	if getSubHash(&subData) != mark {
		CoreLog.Error(logAppend, "sub no hash")
		return
	}
	//获取配置
	configData, err := getConfigByID(subData.ConfigID)
	if err != nil {
		//CoreLog.Error(logAppend, "get sub config, ", err)
		//return
	} else {
		//收尾工作
		expireSubLast(&configData, &subData)
	}
	//删除过期数据
	if err = ClearSubByUser(&ArgsClearSubByUser{
		OrgID:    subData.OrgID,
		UserID:   subData.UserID,
		ConfigID: subData.ConfigID,
	}); err != nil {
		CoreLog.Error(logAppend, "clear sub by user, ", err)
		return
	}
	//强制更新组织用户数据
	OrgUserMod.PushUpdateUserData(subData.OrgID, subData.UserID)
}
