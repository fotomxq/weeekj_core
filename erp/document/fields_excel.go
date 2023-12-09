package ERPDocument

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
)

type FieldsExcelConfig struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//数据表
	Sheets FieldsExcelConfigSheetList `db:"sheets" json:"sheets"`
}

type FieldsExcelConfigSheetList []FieldsExcelConfigSheet

// Value sql底层处理器
func (t FieldsExcelConfigSheetList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsExcelConfigSheetList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsExcelConfigSheet struct {
	//名称
	Name string `db:"name" json:"name" check:"name" min:"1" max:"300"`
	//数据集合
	// 兼容性数据集合，可剔除掉数值类内容（或不需要统计的内容），放入本集合内
	Data string `db:"data" json:"data" check:"des" min:"1" max:"60000" empty:"true"`
	//数据集合
	RowCols FieldsExcelConfigRowColList `db:"row_cols" json:"rowCols"`
}

// Value sql底层处理器
func (t FieldsExcelConfigSheet) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsExcelConfigSheet) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsExcelConfigRowColList []FieldsExcelConfigRowCol

// Value sql底层处理器
func (t FieldsExcelConfigRowColList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsExcelConfigRowColList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsExcelConfigRowCol struct {
	//位置
	Row string `db:"row" json:"row"`
	Col string `db:"col" json:"col"`
	//默认样式
	ClassName string `db:"class_name" json:"className"`
	//样式
	StyleName string `db:"style_name" json:"styleName"`
	//组件默认值
	Val string `db:"val" json:"val"`
	//整数（内部记录用）
	ValInt64 int64 `db:"val_int64" json:"valInt64"`
	//浮点数（内部记录用）
	ValFloat64 float64 `db:"val_float64" json:"valFloat64"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// Value sql底层处理器
func (t FieldsExcelConfigRowCol) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsExcelConfigRowCol) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsExcelSheet struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//文档ID
	DocID int64 `db:"doc_id" json:"docID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"name" min:"1" max:"300"`
	//数据集合
	// 兼容性数据集合，可剔除掉数值类内容（或不需要统计的内容），放入本集合内
	Data string `db:"data" json:"data" check:"des" min:"1" max:"60000" empty:"true"`
}

type FieldsExcelRowCol struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
	//文档ID
	DocID int64 `db:"doc_id" json:"docID" check:"id"`
	//所属文档子表
	SheetID int64 `db:"sheet_id" json:"sheetID"`
	//位置
	Row string `db:"row" json:"row"`
	Col string `db:"col" json:"col"`
	//默认样式
	ClassName string `db:"class_name" json:"className"`
	//样式
	StyleName string `db:"style_name" json:"styleName"`
	//组件默认值
	Val string `db:"val" json:"val"`
	//整数（内部记录用）
	ValInt64 int64 `db:"val_int64" json:"valInt64"`
	//浮点数（内部记录用）
	ValFloat64 float64 `db:"val_float64" json:"valFloat64"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
