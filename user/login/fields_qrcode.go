package UserLogin

import (
	"encoding/json"
	"time"
)

//FieldsQrcode 登陆专用二维码
type FieldsQrcode struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//会话ID
	TokenID int64 `db:"token_id" json:"tokenID"`
	//来源系统
	SystemMark string `db:"system_mark" json:"systemMark"`
	//所属用户
	UserID int64 `db:"user_id" json:"userID"`
	//密钥
	// 二维码根据ID\Key混合计算，APP扫码后识别匹配
	// 二维码包含前缀设计: user_login、APP名称
	Key string `db:"key" json:"key"`
}

//DataQrcode 二维码结构设计
type DataQrcode struct {
	//APP前缀
	// xxx_v2_user_login
	Mark string `json:"mark"`
	//ID
	ID int64 `json:"id"`
	//密钥
	// 二维码根据ID\Key混合计算，APP扫码后识别匹配
	// 二维码包含前缀设计: user_login、APP名称
	Key string `json:"key"`
}

//GetJSON 转化DataQrcode为JSON字符串
func (t *DataQrcode) GetJSON() (string, error) {
	data, err := json.Marshal(t)
	return string(data), err
}

//GetData 转化JSON为DataQrcode
func (t *DataQrcode) GetData(jsonData string) error {
	return json.Unmarshal([]byte(jsonData), &t)
}