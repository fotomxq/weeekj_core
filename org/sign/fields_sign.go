package OrgSign

import "time"

// FieldsSign 签名存储
type FieldsSign struct {
	// ID
	ID int64 `db:"id" json:"id" check:"id" unique:"true"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt" default:"now()"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt" default:"now()"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt" default:"0" index:"true"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgId" check:"id" index:"true"`
	//组织成员ID
	OrgBindID int64 `db:"org_bind_id" json:"orgBindId" check:"id" index:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userId" check:"id" index:"true"`
	//是否默认
	// 一个客体可拥有多个签名，但只能有一个默认签名
	IsDefault bool `db:"is_default" json:"isDefault" index:"true"`
	//签名类型
	// base64 Base64文件; file 文件系统ID
	SignType string `db:"sign_type" json:"signType" check:"des" min:"1" max:"50" index:"true"`
	//是否临时传递
	// 临时传递的签名将在使用后立即删除
	IsTemp bool `db:"is_temp" json:"isTemp" index:"true"`
	//签名数据
	SignData string `db:"sign_data" json:"signData" check:"des" min:"1" max:"-1" empty:"true"`
	//文件ID
	FileID int64 `db:"file_id" json:"fileID" check:"id" index:"true"`
}
