package ServiceAD2

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

//FieldsAD 广告设置
type FieldsAD struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//分区标识码
	// 同一个组织下唯一，具体行为交给mode识别和前端处理
	Mark string `db:"mark" json:"mark"`
	//结构
	Data FieldsADChildList `json:"data"`
}

type FieldsADChildList []FieldsADChild

// Value sql底层处理器
func (t FieldsADChildList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsADChildList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsADChild struct {
	//广告模式
	// none 展示类广告
	// blog_content 跳转文章; mall_core_product 跳转商品; mall_core_sort 跳转商品分类; user_ticket 跳转用户票据及优惠券; user_sub 用户会员; finance_user_deposit 用户充值
	// weixin_wxx 微信小程序广告
	Mode string `db:"mode" json:"mode" check:"mark"`
	//绑定ID
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//用于一些第三方识别系统
	BindMark string `db:"bind_mark" json:"bindMark"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id"`
}

// Value sql底层处理器
func (t FieldsADChild) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsADChild) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
