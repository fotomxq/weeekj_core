package FinanceMargePay

import "time"

//FieldsQrcode 聚合支付二维码
type FieldsQrcode struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//交易过期时间
	// 如果提交空的时间，将直接按照过期处理
	ExpireAt time.Time `db:"expire_at" json:"expireAt" check:"isoTime"`
	//支付Key
	PayKey string `db:"pay_key" json:"payKey"`
	//支付ID
	// 二维码将识别对应的支付ID，并根据端的区别，自动修正支付请求的端数据包，方便调用支付请求处理
	PayID int64 `db:"pay_id" json:"payID" check:"id"`
}
