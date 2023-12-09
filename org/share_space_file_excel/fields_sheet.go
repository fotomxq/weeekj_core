package OrgShareSpaceFileExcel

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type FieldsSheetList []FieldsSheet

// Value sql底层处理器
func (t FieldsSheetList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsSheetList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsSheet struct {
	//表名称
	SheetName string `db:"sheet_name" json:"sheetName"`
	//单元格内容列
	Data FieldsSheetDataList `db:"data" json:"data"`
}

type FieldsSheetDataList []FieldsSheetData

// Value sql底层处理器
func (t FieldsSheetDataList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsSheetDataList) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}

type FieldsSheetData struct {
	//位置
	Key string `db:"key" json:"key"`
	//值
	Val string `db:"val" json:"val"`
	//合并附近单元格
	// 当前单元格为原点，向右下侧开始衍生，如果为0则不合并
	// 行
	MargeRow int `db:"marge_row" json:"margeRow"`
	// 列
	MargeCel int `db:"marge_cel" json:"margeCel"`
	//样式约定
	Style string `db:"style" json:"style"`
}

// Value sql底层处理器
func (t FieldsSheetData) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *FieldsSheetData) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &t)
}
