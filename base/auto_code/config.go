package BaseAutoCode

import (
	"errors"
	"fmt"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

// ArgsGetConfigList 获取编码配置列表参数
type ArgsGetConfigList struct {
	//分页参数
	Pages CoreSQL2.ArgsPages `json:"pages"`
	//系统模块标识码
	// eq: user_core 用户模块
	ModuleCode string `db:"module_code" json:"moduleCode" check:"code" empty:"true"`
	//模块内分支标识码
	// eq: user_core_core 用户模块的用户表
	BranchCode string `db:"branch_code" json:"branchCode" check:"code" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetConfigList 获取编码配置列表
func GetConfigList(args *ArgsGetConfigList) (dataList []FieldsConfig, dataCount int64, err error) {
	dataCount, err = configDB.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id", "create_at", "update_at", "delete_at", "name"}).SetPages(args.Pages).SetDeleteQuery("delete_at", args.IsRemove).SetStringQuery("module_code", args.ModuleCode).SetStringQuery("branch_code", args.BranchCode).SetSearchQuery([]string{"name"}, args.Search).SelectList("").ResultAndCount(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		vData := getConfigByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList[k] = vData
	}
	return
}

// ArgsGetConfigByID 获取配置参数
type ArgsGetConfigByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// GetConfigByID 获取配置
func GetConfigByID(args *ArgsGetConfigByID) (data FieldsConfig, err error) {
	data = getConfigByID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// GetConfigNameByID 获取配置名称
func GetConfigNameByID(id int64) (name string) {
	data := getConfigByID(id)
	if data.ID < 1 {
		return
	}
	return data.Name
}

// ArgsCreateConfig 创建配置参数
type ArgsCreateConfig struct {
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
	// 对应的字段名用{}包裹，支持多个字段组合；原则上仅支持英文字符（自动大写）、数字、下划线；不支持特殊字符
	// eq: {prefix}{auto_number}
	CustomRule string `db:"custom_rule" json:"customRule" check:"des" min:"1" max:"255" empty:"true"`
}

// CreateConfig 创建配置
func CreateConfig(args *ArgsCreateConfig) (id int64, err error) {
	//检查标识码是否存在
	data := getConfigByCode(args.ModuleCode, args.BranchCode)
	if data.ID > 0 {
		err = errors.New("code is exist")
		return
	}
	//创建数据
	id, err = configDB.Insert().SetFields([]string{"name", "module_code", "branch_code", "prefix", "auto_number", "is_global_unique", "is_branch_unique", "is_log", "is_enable", "custom_rule"}).Add(map[string]any{
		"name":             args.Name,
		"module_code":      args.ModuleCode,
		"branch_code":      args.BranchCode,
		"prefix":           args.Prefix,
		"auto_number":      args.AutoNumber,
		"is_global_unique": args.IsGlobalUnique,
		"is_branch_unique": args.IsBranchUnique,
		"is_log":           args.IsLog,
		"is_enable":        args.IsEnable,
		"custom_rule":      args.CustomRule,
	}).ExecAndResultID()
	if err != nil {
		return
	}
	//反馈
	return
}

// ArgsUpdateConfig 修改配置参数
type ArgsUpdateConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
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
	// 对应的字段名用{}包裹，支持多个字段组合；原则上仅支持英文字符（自动大写）、数字、下划线；不支持特殊字符
	// eq: {prefix}{auto_number}
	CustomRule string `db:"custom_rule" json:"customRule" check:"des" min:"1" max:"255" empty:"true"`
}

// UpdateConfig 修改配置
func UpdateConfig(args *ArgsUpdateConfig) (err error) {
	//检查标识码是否存在
	data := getConfigByCode(args.ModuleCode, args.BranchCode)
	if data.ID > 0 && args.ID != data.ID {
		err = errors.New("code is exist")
		return
	}
	//更新数据
	err = configDB.Update().SetFields([]string{"name", "module_code", "branch_code", "prefix", "auto_number", "is_global_unique", "is_branch_unique", "is_log", "is_enable", "custom_rule"}).NeedUpdateTime().AddWhereID(args.ID).NamedExec(map[string]any{
		"name":             args.Name,
		"module_code":      args.ModuleCode,
		"branch_code":      args.BranchCode,
		"prefix":           args.Prefix,
		"auto_number":      args.AutoNumber,
		"is_global_unique": args.IsGlobalUnique,
		"is_branch_unique": args.IsBranchUnique,
		"is_log":           args.IsLog,
		"is_enable":        args.IsEnable,
		"custom_rule":      args.CustomRule,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteConfigCache(args.ID)
	//反馈
	return
}

// ArgsDeleteConfig 删除配置参数
type ArgsDeleteConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteConfig 删除配置
func DeleteConfig(args *ArgsDeleteConfig) (err error) {
	//删除数据
	err = configDB.Delete().NeedSoft(true).AddWhereID(args.ID).ExecNamed(nil)
	if err != nil {
		return
	}
	//删除缓冲
	deleteConfigCache(args.ID)
	//反馈
	return
}

// getConfigByCode 通过标识获取配置
func getConfigByCode(moduleCode string, branchCode string) (data FieldsConfig) {
	_ = configDB.Get().SetFieldsOne([]string{"id"}).SetDeleteQuery("delete_at", false).SetStringQuery("module_code", moduleCode).SetStringQuery("branch_code", branchCode).NeedLimit().Result(&data)
	if data.ID < 1 {
		return
	}
	data = getConfigByID(data.ID)
	if data.ID < 1 {
		return
	}
	return

}

// getConfigByID 通过ID获取配置
func getConfigByID(id int64) (data FieldsConfig) {
	cacheMark := getConfigCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := configDB.Get().SetFieldsOne([]string{"id", "create_at", "update_at", "delete_at", "name", "module_code", "branch_code", "prefix", "auto_number", "is_global_unique", "is_branch_unique", "is_log", "is_enable", "custom_rule"}).GetByID(id).NeedLimit().Result(&data)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheConfigTime)
	return
}

// 缓冲
func getConfigCacheMark(id int64) string {
	return fmt.Sprint("base:auto:code:config.", id)
}

func deleteConfigCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getConfigCacheMark(id))
}
