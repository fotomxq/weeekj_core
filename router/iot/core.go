package RouterIOT

import (
	IOTDevice "github.com/fotomxq/weeekj_core/v5/iot/device"
	RouterReport "github.com/fotomxq/weeekj_core/v5/router/report"
	"github.com/gin-gonic/gin"
)

// ArgsCheckDeviceAndOrg 检查设备是否存在且授权参数
type ArgsCheckDeviceAndOrg struct {
	//设备分组
	GroupMark string `db:"group_mark" json:"groupMark" check:"mark"`
	//设备编号
	// 同一个分组下，必须唯一
	Code string `db:"code" json:"code" check:"mark"`
	//时间戳
	NowTime int64 `db:"now_time" json:"nowTime"`
	//随机码
	Rand string `db:"rand" json:"rand"`
	//key计算结果
	// key+code+时间戳+随机码
	Key string `db:"key" json:"key"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// CheckDeviceAndOrg 检查设备是否存在且授权
func CheckDeviceAndOrg(c *gin.Context, args *ArgsCheckDeviceAndOrg) (b bool) {
	_, b = CheckDeviceAndOrgReturnDevice(c, args)
	return
}

func CheckDeviceAndOrgReturnDevice(c *gin.Context, args *ArgsCheckDeviceAndOrg) (deviceID int64, b bool) {
	var err error
	deviceID, err = IOTDevice.CheckDeviceKeyAndDeviceID(&IOTDevice.ArgsCheckDeviceKey{
		GroupMark: args.GroupMark,
		Code:      args.Code,
		NowTime:   args.NowTime,
		Rand:      args.Rand,
		Key:       args.Key,
	})
	if err != nil || deviceID < 1 {
		RouterReport.BaseBool(c, "device", false, "设备不存在")
		//CoreLog.Warn("check device and org error, ", err)
		return
	}
	if err = IOTDevice.CheckOperateNoData(&IOTDevice.ArgsCheckOperate{
		DeviceID: deviceID,
		OrgID:    args.OrgID,
	}); err != nil {
		RouterReport.BaseBool(c, "device", false, "设备不存在")
		//CoreLog.Warn("check device and org error, no operate data, ", err)
		return
	}
	b = true
	return
}
