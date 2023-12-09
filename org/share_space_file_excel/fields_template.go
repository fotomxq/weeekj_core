package OrgShareSpaceFileExcel

import "time"

type FieldsTemplate struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//名称
	Name string `db:"name" json:"name"`
	//结构体设计
	SheetData FieldsSheetList `db:"sheet_data" json:"sheetData"`
}
