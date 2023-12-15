package ERPDocument

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	ERPCore "github.com/fotomxq/weeekj_core/v5/erp/core"
	"time"
)

type FieldsConfig struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//发布状态
	PublishAt time.Time `db:"publish_at" json:"publishAt"`
	//hash
	// 如果hash和提交hash不同，服务端将自动拒绝更新，避免流处理异常
	Hash string `db:"hash" json:"hash" check:"sha1"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"name" min:"1" max:"300"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"6000" empty:"true"`
	//封面ID
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//文档类型
	// custom 自定义; doc 普通文稿; excel 表格
	DocType string `db:"doc_type" json:"docType"`
	//节点组件
	ComponentList ERPCore.FieldsComponentDefineList `db:"component_list" json:"componentList"`
	//列表展示数据
	ListShow FieldsConfigListShows `db:"list_show" json:"listShow"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

type FieldsConfigListShows []FieldsConfigListShow

// Value sql底层处理器
func (t FieldsConfigListShows) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsConfigListShows) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsConfigListShow struct {
	//内容类型
	// field 字段; component 组件
	Mode string `db:"mode" json:"mode"`
	//数据转化方式
	// none 直接展示; auto 根据组件类型自动识别
	DataType string `db:"data_type" json:"dataType"`
	//组件key
	Key string `db:"key" json:"key"`
	//默认宽度
	DefaultWidth int `db:"default_width" json:"defaultWidth"`
}

// Value sql底层处理器
func (t FieldsConfigListShow) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsConfigListShow) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
