package IOTDevice

import (
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreLog "github.com/fotomxq/weeekj_core/v5/core/log"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/nats-io/nats.go"
)

func subNats() {
	//删除旧的日志数据
	CoreNats.SubDataByteNoErr("/iot/device/auto_log", subNatsDeleteAutoLog)
	//删除操作权限过期
	CoreNats.SubDataByteNoErr("/base/expire_tip/expire", subNatsOperateExpire)
	//标记设备掉线处理
	CoreNats.SubDataByteNoErr("/base/expire_tip/expire", subNatsDeviceOnlineExpire)
	//删除自动化模板
	CoreNats.SubDataByteNoErr("/iot/device/auto_info_template", subNatsAutoInfoTemplate)
}

// 删除旧的日志数据
func subNatsDeleteAutoLog(_ *nats.Msg, _ string, _ int64, _ string, _ []byte) {
	autoLogDeleteBlocker.CheckWait(0, "", func(_ int64, _ string) {
		//删除旧的数据
		_, _ = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "iot_core_auto_log", "create_at < :create_at", map[string]interface{}{
			"create_at": CoreFilter.GetNowTimeCarbon().SubHour(),
		})
	})
}

// 删除操作权限过期
func subNatsOperateExpire(_ *nats.Msg, action string, operateID int64, _ string, _ []byte) {
	if action != "iot_device_operate" {
		return
	}
	if err := deleteOperateByID(operateID); err != nil {
		CoreLog.Error("iot device sub nats operate expire, ", err)
	}
}

// 设备掉线到期处理
func subNatsDeviceOnlineExpire(_ *nats.Msg, action string, deviceID int64, _ string, _ []byte) {
	if action != "iot_device_group_online" {
		return
	}
	deviceData := getDeviceByID(deviceID)
	if deviceData.ID < 1 {
		return
	}
	if CoreSQL.CheckTimeThanNow(deviceData.LastAt) {
		return
	}
	if err := UpdateDeviceOnline(&ArgsUpdateDeviceOnline{
		ID:       deviceID,
		IsOnline: false,
	}); err != nil {
		CoreLog.Error("iot device sub nats device online expire, ", err)
	}
}

// 自动化模板变动
func subNatsAutoInfoTemplate(_ *nats.Msg, action string, templateID int64, _ string, _ []byte) {
	appendLog := "iot device sub nats auto info template, "
	switch action {
	case "delete":
		//删除模板
		var dataList []FieldsAutoInfo
		_ = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id FROM iot_core_auto_info WHERE template_id = $1", templateID)
		if len(dataList) > 0 {
			for _, v := range dataList {
				if err := DeleteAutoInfo(&ArgsDeleteAutoInfo{
					ID:    v.ID,
					OrgID: -1,
				}); err != nil {
					CoreLog.Error(appendLog, "delete auto info by id: ", v.ID, ", err: ", err)
				}
			}
		}
	}
}
