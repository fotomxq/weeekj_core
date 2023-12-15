package IOTQrcode

import (
	"encoding/base64"
	"errors"
	"fmt"

	BaseWeixinWXXQRCodeCore "github.com/fotomxq/weeekj_core/v5/base/weixin/wxx/qrcode"
	IOTDevice "github.com/fotomxq/weeekj_core/v5/iot/device"
)

// ArgsMakeWeixinWXX 微信小程序二维码参数
type ArgsMakeWeixinWXX struct {
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID" check:"id"`
	//二维码类型
	// weixin_wxx 微信小程序设备group+code二维码; weixin_wxx_id 微信小程序设备ID二维码
	QrcodeType string `db:"qrcode_type" json:"qrcodeType" check:"mark"`
	//尺寸
	// eg: 430
	Size int `db:"size" json:"size" check:"intThan0"`
	//商户ID
	// 可以留空，则走平台微信小程序主体
	MerchantID int64 `json:"merchantID" check:"id" empty:"true"`
	//页面地址
	// eg: pages/index
	Page string `json:"page"`
	//是否需要透明底色
	IsHyaline bool `json:"isHyaline"`
	//自动配置线条颜色
	// 为 false 时生效, 使用 rgb 设置颜色 十进制表示
	AutoColor bool `json:"autoColor"`
	//色调
	// 50
	R string `json:"r"`
	G string `json:"g"`
	B string `json:"b"`
}

// MakeWeixinWXX 微信小程序二维码
func MakeWeixinWXX(args *ArgsMakeWeixinWXX) (data string, err error) {
	//尺寸不能太小
	if args.Size < 10 {
		err = errors.New("size too small")
		return
	}
	if args.Size > 1024 {
		err = errors.New("size too big")
		return
	}
	//获取设备信息
	var deviceData IOTDevice.FieldsDevice
	deviceData, err = IOTDevice.GetDeviceByID(&IOTDevice.ArgsGetDeviceByID{
		ID:    args.DeviceID,
		OrgID: -1,
	})
	if err != nil {
		err = errors.New("device not exist, " + err.Error())
		return
	}
	//设备组数据
	var groupData IOTDevice.FieldsGroup
	groupData, err = IOTDevice.GetGroupByID(&IOTDevice.ArgsGetGroupByID{
		ID: deviceData.GroupID,
	})
	if err != nil {
		err = errors.New("group not exist, " + err.Error())
		return
	}
	//根据类型输出数据
	switch args.QrcodeType {
	case "weixin_wxx":
		var dataByte []byte
		dataByte, err = BaseWeixinWXXQRCodeCore.GetQRByParam(&BaseWeixinWXXQRCodeCore.ArgsGetQRByParam{
			MerchantID: args.MerchantID,
			Page:       args.Page,
			Param:      fmt.Sprint(groupData.Mark, "@", deviceData.Code),
			Width:      args.Size,
			IsHyaline:  args.IsHyaline,
			AutoColor:  args.AutoColor,
			R:          args.R,
			G:          args.G,
			B:          args.B,
		})
		if err != nil {
			return
		}
		data = base64.StdEncoding.EncodeToString(dataByte)
		return
	case "weixin_wxx_id":
		var dataByte []byte
		dataByte, err = BaseWeixinWXXQRCodeCore.GetQRByParam(&BaseWeixinWXXQRCodeCore.ArgsGetQRByParam{
			MerchantID: args.MerchantID,
			Page:       args.Page,
			Param:      fmt.Sprint(deviceData.ID),
			Width:      args.Size,
			IsHyaline:  args.IsHyaline,
			AutoColor:  args.AutoColor,
			R:          args.R,
			G:          args.G,
			B:          args.B,
		})
		if err != nil {
			return
		}
		data = base64.StdEncoding.EncodeToString(dataByte)
		return
	default:
		err = errors.New("qrcode type error")
		return
	}
}
