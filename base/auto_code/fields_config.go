package BaseAutoCode

import "time"

// FieldsConfig 自动编码字段配置
type FieldsConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//删除时间
	DeleteAt time.Time `db:"delete_at" json:"deleteAt"`
	//名称
	Name string `db:"name" json:"name" check:"name" min:"1" max:"300"`
	//系统模块标识码
	// eq: user_core 用户模块
	ModuleCode string `db:"module_code" json:"moduleCode" check:"code"`
	//模块内分支标识码
	// eq: user_core_core 用户模块的用户表
	BranchCode string `db:"branch_code" json:"branchCode" check:"code"`
	//编码前缀
	// eq: UC 用户中心
	Prefix string `db:"prefix" json:"prefix" check:"code" empty:"true"`
	//是否自动按序号生成
	AutoNumber bool `db:"auto_number" json:"autoNumber" check:"bool"`
	//自动生成序号预留位数
	AutoNumberLen int `db:"auto_number_len" json:"autoNumberLen" check:"int" min:"1" max:"10" empty:"true"`
	//是否全局强制唯一
	IsGlobalUnique bool `db:"is_global_unique" json:"isGlobalUnique" check:"bool"`
	//模块内是否强制唯一
	IsBranchUnique bool `db:"is_branch_unique" json:"isBranchUnique" check:"bool"`
	//是否记录日志
	// 如果不记录日志，将无法实现上述排重功能
	IsLog bool `db:"is_log" json:"isLog" check:"bool"`
	//是否启用
	IsEnable bool `db:"is_enable" json:"isEnable" check:"bool"`
	//自定义生成规则
	// 对应的字段名用","分割，支持多个字段组合；原则上仅支持英文字符（自动大写）、数字、下划线；不支持特殊字符
	// eq: prefix,auto_number
	CustomRule string `db:"custom_rule" json:"customRule" check:"des" min:"1" max:"255" empty:"true"`
}
