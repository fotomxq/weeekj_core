package BaseAutoCode

import "time"

// FieldsLog 自动编码字段配置
type FieldsLog struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//系统模块标识码
	// eq: user_core 用户模块
	ModuleCode string `db:"module_code" json:"moduleCode" check:"code"`
	//模块内分支标识码
	// eq: user_core_core 用户模块的用户表
	BranchCode string `db:"branch_code" json:"branchCode" check:"code"`
	//采用配置ID
	ConfigID int64 `db:"config_id" json:"configId" check:"id"`
	//编码
	Code string `db:"code" json:"code" check:"code"`
}
