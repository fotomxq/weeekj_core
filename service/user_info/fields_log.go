package ServiceUserInfo

import "time"

// FieldsLog 人员修改日志
type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//组织ID
	// 允许平台方的0数据，该数据可能来源于其他领域
	OrgID int64 `db:"org_id" json:"orgID"`
	//档案ID
	InfoID int64 `db:"info_id" json:"infoID"`
	//修改的位置
	// 1. 字段
	// 2. 或扩展参数指定的内容，例如params.[mark]
	// 3. 其他内容采用.形式跨越记录
	// 4. room.in 入驻房间变更
	ChangeMark string `db:"change_mark" json:"changeMark"`
	ChangeDes  string `db:"change_des" json:"changeDes"`
	//修改前描述
	OldDes string `db:"old_des" json:"oldDes"`
	//修改后描述
	NewDes string `db:"new_des" json:"newDes"`
}
