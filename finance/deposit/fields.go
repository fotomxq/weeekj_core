package FinanceDeposit

import (
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	"time"
)

// FieldsDepositType 储蓄资金池
type FieldsDepositType struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//创建时间
	CreateAt time.Time `db:"create_at" json:"createAt"`
	//更新时间
	UpdateAt time.Time `db:"update_at" json:"updateAt"`
	//HASH
	UpdateHash string `db:"update_hash" json:"updateHash"`
	//来源
	CreateInfo CoreSQLFrom.FieldsFrom `db:"create_info" json:"createInfo"`
	//存储来源
	// 如不同的加盟商，储蓄均有所差异，也可以留空不指定，则为平台总的储蓄资金池
	FromInfo CoreSQLFrom.FieldsFrom `db:"from_info" json:"fromInfo"`
	//储蓄配置标识码
	// 可用于同一类货币下，多个用途，如赠送的储值额度、或用户自行充值的额度
	// user 用户自己储值 ; deposit 押金 ; free 免费赠送额度 ; ... 特定系统下的充值模块
	ConfigMark string `db:"config_mark" json:"configMark"`
	//储蓄金额
	SavePrice int64 `db:"save_price" json:"savePrice"`
}

// FieldsConfigType 存储配置
// 只有平台允许设计统一的储蓄配置池，组织只能使用相关储蓄池
type FieldsConfigType struct {
	//标识码
	// 可用于同一类货币下，多个用途，如赠送的储值额度、或用户自行充值的额度
	// user 用户自己储值 ; deposit 押金 ; free 免费赠送额度 ; ... 特定系统下的充值模块
	Mark string `db:"mark" json:"mark"`
	//显示名称
	Name string `db:"name" json:"name"`
	//备注
	Des string `db:"des" json:"des"`
	//储蓄货币类型
	// 采用CoreCurrency匹配
	Currency int `db:"currency" json:"currency"`
	//能否取出
	// 如果能，则允许用户使用取出接口
	TakeOut bool `db:"take_out" json:"takeOut"`
	//取款最低限额
	// 低于该资金禁止取款，同时需启动是否可取
	TakeLimit int64 `db:"take_limit" json:"takeLimit"`
	//单次存款最低限额
	OnceSaveMinLimit int64 `db:"once_save_min_limit" json:"onceSaveMinLimit"`
	//单次存款最大限额
	OnceSaveMaxLimit int64 `db:"once_save_max_limit" json:"onceSaveMaxLimit"`
	//单次取款最低限额
	OnceTakeMinLimit int64 `db:"once_take_min_limit" json:"onceTakeMinLimit"`
	//单次取款最大限额
	OnceTakeMaxLimit int64 `db:"once_take_max_limit" json:"onceTakeMaxLimit"`
	//扩展参数设计
	Configs CoreSQLConfig.FieldsConfigsType `db:"configs" json:"configs"`
}
