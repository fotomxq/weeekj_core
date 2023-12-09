package UserLogin

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	OrgCoreCore "gitee.com/weeekj/weeekj_core/v5/org/core"
	UserCore "gitee.com/weeekj/weeekj_core/v5/user/core"
	"time"
)

// FieldsSave 临时存储数据包
// 统一一分钟有效期，提取失效
type FieldsSave struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//过期时间
	ExpireAt time.Time `db:"expire_at" json:"expireAt"`
	//提取密钥
	Key string `db:"key" json:"key"`
	//数据
	Data FieldsSaveReportData `db:"data" json:"data"`
}

// FieldsSaveReportData 反馈的数据结构
type FieldsSaveReportData struct {
	//会话
	Token FieldsSaveReportDataToken `db:"token" json:"token"`
	//用户脱敏数据
	UserData UserCore.DataUserDataType `db:"user_data" json:"userData"`
	//绑定关系和组织关系
	OrgBindData []OrgCoreCore.DataGetBindByUserMarge `db:"org_bind_data" json:"orgBindData"`
	//文件集
	// id => url
	FileList map[int64]string `db:"file_list" json:"fileList"`
}

// Value sql底层处理器
func (t FieldsSaveReportData) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsSaveReportData) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsSaveReportDataToken struct {
	Token int64  `db:"token" json:"token"`
	Key   string `db:"key" json:"key"`
}

// Value sql底层处理器
func (t FieldsSaveReportDataToken) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsSaveReportDataToken) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
