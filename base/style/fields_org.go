package BaseStyle

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"github.com/lib/pq"
	"time"
)

// FieldsOrg 组织定义表
type FieldsOrg struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//样式表ID
	StyleID int64 `db:"style_id" json:"styleID"`
	//组件列
	Components FieldsOrgComponents `db:"components" json:"components"`
	//标题
	Title string `db:"title" json:"title"`
	//描述
	Des string `db:"des" json:"des"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID"`
	//描述文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles"`
	//附加参数
	// 布局相关的自定会参数，都将在此处定义
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

type FieldsOrgComponents []FieldsOrgComponent

// Value sql底层处理器
func (t FieldsOrgComponents) Value() (driver.Value, error) {
	/**
	buf := &bytes.Buffer{}
	buf.Write([]byte{'{'})
	l := len(t)
	for i, k := range t {
		_, _ = fmt.Fprintf(buf, "\"%s\": \"%v\"", fmt.Sprint(i), k)
		if i < l-1 {
			buf.WriteByte(',')
		}
	}
	buf.Write([]byte{'}'})
	return buf.Bytes(), nil
	*/
	return json.Marshal(t)
}

func (t *FieldsOrgComponents) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsOrgComponent struct {
	//组件ID
	ComponentID int64 `db:"component_id" json:"componentID" check:"id"`
	//附加参数
	// 布局相关的自定会参数，都将在此处定义
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// Value sql底层处理器
func (t FieldsOrgComponent) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsOrgComponent) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
