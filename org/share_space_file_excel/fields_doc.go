package OrgShareSpaceFileExcel

import "time"

type FieldsDoc struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//采用模板
	TemplateID int64 `db:"template_id" json:"templateID"`
	//结构体设计
	SheetData FieldsSheetList `db:"sheet_data" json:"sheetData"`
}
