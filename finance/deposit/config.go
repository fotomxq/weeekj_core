package FinanceDeposit

import (
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLMarks "gitee.com/weeekj/weeekj_core/v5/core/sql/marks"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetConfigList 获取列表参数
type ArgsGetConfigList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList
	//搜索
	Search string
}

// GetConfigList 获取列表
func GetConfigList(args *ArgsGetConfigList) (dataList []FieldsConfigType, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.Search != "" {
		where = where + "(des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	} else {
		maps = nil
	}
	if where == "" {
		where = "true"
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"finance_deposit_config",
		"mark",
		"SELECT mark, name, des, currency, take_out, take_limit, once_save_min_limit, once_save_max_limit, once_take_min_limit, once_take_max_limit, configs FROM finance_deposit_config WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"mark", "name"},
	)
	return
}

// ArgsGetConfigByMark 获取mark参数
type ArgsGetConfigByMark struct {
	//标识码
	Mark string `json:"mark" check:"mark"`
}

// GetConfigByMark 获取mark
func GetConfigByMark(args *ArgsGetConfigByMark) (data FieldsConfigType, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT mark, name, des, currency, take_out, take_limit, once_save_min_limit, once_save_max_limit, once_take_min_limit, once_take_max_limit, configs FROM finance_deposit_config WHERE mark = $1", args.Mark)
	return
}

// ArgsGetConfigMore 获取一组配置参数
type ArgsGetConfigMore struct {
	//Mark列
	Marks pq.StringArray `json:"marks" check:"marks"`
}

// GetConfigMore 获取一组配置
func GetConfigMore(args *ArgsGetConfigMore) (dataList []FieldsConfigType, err error) {
	err = CoreSQLMarks.GetMarks(&dataList, "finance_deposit_config", "mark, name, des, currency, take_out, take_limit, once_save_min_limit, once_save_max_limit, once_take_min_limit, once_take_max_limit, configs", args.Marks)
	return
}

// GetConfigMoreMap 获取一组配置名称组
func GetConfigMoreMap(args *ArgsGetConfigMore) (data map[string]string, err error) {
	data, err = CoreSQLMarks.GetMarksName("finance_deposit_config", args.Marks)
	return
}

// GetConfigName 获取配置名称
func GetConfigName(mark string) (name string) {
	_ = Router2SystemConfig.MainDB.Get(&name, "SELECT name FROM finance_deposit_config WHERE mark = $1", mark)
	return
}

// ArgsSetConfig 保存或设置新的参数
type ArgsSetConfig struct {
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

// SetConfig 保存或设置新的
func SetConfig(args *ArgsSetConfig) (data FieldsConfigType, err error) {
	data, err = GetConfigByMark(&ArgsGetConfigByMark{
		Mark: args.Mark,
	})
	if err != nil {
		_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO finance_deposit_config (mark, name, des, currency, take_out, take_limit, once_save_min_limit, once_save_max_limit, once_take_min_limit, once_take_max_limit, configs) VALUES (:mark,:name,:des,:currency,:take_out,:take_limit,:once_save_min_limit,:once_save_max_limit,:once_take_min_limit,:once_take_max_limit,:configs)", args)
	} else {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE finance_deposit_config SET  name = :name, des = :des, currency = :currency, take_out = :take_out, take_limit = :take_limit, once_save_min_limit = :once_save_min_limit, once_save_max_limit = :once_save_max_limit, once_take_min_limit = :once_take_min_limit, once_take_max_limit = :once_take_max_limit, configs = :configs WHERE mark = :mark", args)
	}
	return
}

// ArgsDeleteConfigByMark 删除参数
type ArgsDeleteConfigByMark struct {
	//标识码
	Mark string `db:"mark" json:"mark"`
}

// DeleteConfigByMark 删除
func DeleteConfigByMark(args *ArgsDeleteConfigByMark) (err error) {
	_, err = CoreSQL.DeleteOne(Router2SystemConfig.MainDB.DB, "finance_deposit_config", "mark", args)
	return
}
