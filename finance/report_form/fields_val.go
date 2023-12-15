package FinanceReportForm

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	"time"
)

type FieldsVal struct {
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
	//所属文件
	FileID int64 `db:"file_id" json:"fileID"`
	//列ID
	ColID int64 `db:"col_id" json:"colID"`
	//标识码
	// 位置标识码：A1\B1，字母代表列、数字代表行
	Mark string `db:"mark" json:"mark"`
	//值
	Val      string  `db:"val" json:"val"`
	ValFloat float64 `db:"val_float" json:"valFloat"`
	ValInt   int64   `db:"val_int" json:"valInt"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}
