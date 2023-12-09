package Router2Mid

import (
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	IOTDevice "gitee.com/weeekj/weeekj_core/v5/iot/device"
	"github.com/gin-gonic/gin"
)

// ArgsIOTData 检查设备是否存在且授权参数
type ArgsIOTData struct {
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

// checkDeviceAndOrg 检查设备是否存在且授权
func checkDeviceAndOrg(c *gin.Context, args *ArgsIOTData) (b bool) {
	_, b = checkDeviceAndOrgReturnDevice(c, args)
	return
}

func checkDeviceAndOrgReturnDevice(c *gin.Context, args *ArgsIOTData) (deviceID int64, b bool) {
	appendLog := "check device and ord return device, "
	var err error
	deviceID, err = IOTDevice.CheckDeviceKeyAndDeviceID(&IOTDevice.ArgsCheckDeviceKey{
		GroupMark: args.GroupMark,
		Code:      args.Code,
		NowTime:   args.NowTime,
		Rand:      args.Rand,
		Key:       args.Key,
	})
	if err != nil || deviceID < 1 {
		//reportGin(c, false, 0, err, "", false, "err_device_lost", 0, nil)
		CoreLog.Warn(appendLog, "check device data, ", err)
		return
	}
	if err = IOTDevice.CheckOperateNoData(&IOTDevice.ArgsCheckOperate{
		DeviceID: deviceID,
		OrgID:    args.OrgID,
	}); err != nil {
		//reportGin(c, false, 0, err, "", false, "err_device_lost", 0, nil)
		CoreLog.Warn(appendLog, "check device operate, ", err)
		return
	}
	b = true
	return
}
