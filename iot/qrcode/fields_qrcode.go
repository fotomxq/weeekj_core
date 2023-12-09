package IOTQrcode

import "time"

//FieldsQrcode 二维码存储
// 该表设计暂时不用，未来系统承载量巨大时，再考虑启动
type FieldsQrcode struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//设备ID
	DeviceID int64 `db:"device_id" json:"deviceID"`
	//二维码类型
	// normal 普通设备group+code二维码; normal_key 带key的二维码; normal_connect 用于快速连接的二维码（带所有必要信息）; normal_id 普通设备ID二维码
	// weixin_wxx 微信小程序设备group+code二维码; weixin_wxx_id 微信小程序设备ID二维码
	QrcodeType string `db:"qrcode_type" json:"qrcodeType"`
	//二维码内容
	// 图形base数据集
	Data string `db:"data" json:"data"`
}