package UserSub2

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	"time"
)

// FieldsConfig 配置文件
type FieldsConfig struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	// 留空则为平台级别会员设置
	// 每个组织或平台只有一条数据
	OrgID int64 `db:"org_id" json:"orgID"`
	//价格等信息
	ConfigData FieldsConfigDataList `db:"config_data" json:"configData"`
	//样式标识码
	StyleMark string `db:"style_mark" json:"styleMark"`
	//总描述
	Des string `db:"des" json:"des"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

type FieldsConfigDataList []FieldsConfigData

// Value sql底层处理器
func (t FieldsConfigDataList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsConfigDataList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsConfigData struct {
	//会员标识码
	// 指定多个标识码，用逗号分隔，前端识别
	// eg: mark0
	Mark string `db:"mark" json:"mark"`
	//名称
	Name string `db:"name" json:"name"`
	//子描述
	// 如果存在，用户点击将覆盖总描述
	Des string `db:"des" json:"des"`
	//会员对应时间长度
	// 1y 1年 / 1m 1月 / 1w 一周 / 1d 一天
	// eg: 1y
	AddTime string `db:"add_time" json:"addTime"`
	//会员价格
	// 单位：分
	Price int64 `db:"price" json:"price"`
	//折扣价格
	// 如果设置为0，则不展示折扣价格
	OldPrice int64 `db:"old_price" json:"oldPrice"`
}

// Value sql底层处理器
func (t FieldsConfigData) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsConfigData) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
