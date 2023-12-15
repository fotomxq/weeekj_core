package OrgActivity

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLIDs "github.com/fotomxq/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

//TODO：注意，所有接口实际没有对接！！！注意修正后再测试

// ArgsGetConfigList 获取配置列表参数
type ArgsGetConfigList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetConfigList 获取配置列表
func GetConfigList(args *ArgsGetConfigList) (dataList []FieldsConfig, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.Search != "" {
		where = where + " AND (title ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "nursing_info_pay_config"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, org_id, title, cover_file_id, des_files, style_id, currency, price, price_old, time_type, time_n, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// ArgsGetConfig 获取指定配置参数
type ArgsGetConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetConfig 获取指定配置
func GetConfig(args *ArgsGetConfig) (data FieldsConfig, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, title, cover_file_id, des_files, style_id, currency, price, price_old, time_type, time_n, params FROM nursing_info_pay_config WHERE id = $1 AND ($2 < 1 OR org_id = $2)", args.ID, args.OrgID)
	return
}

// ArgsGetConfigs 获取多个配置参数
type ArgsGetConfigs struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

func GetConfigs(args *ArgsGetConfigs) (dataList []FieldsConfig, err error) {
	err = CoreSQLIDs.GetIDsOrgAndDelete(&dataList, "nursing_info_pay_config", "id, create_at, update_at, delete_at, org_id, title, cover_file_id, des_files, style_id, currency, price, price_old, time_type, time_n, params", args.IDs, args.OrgID, args.HaveRemove)
	return
}

// GetConfigsName 获取多个配置名称
func GetConfigsName(args *ArgsGetConfigs) (data map[int64]string, err error) {
	data, err = CoreSQLIDs.GetIDsOrgTitleAndDelete("nursing_info_pay_config", args.IDs, args.OrgID, args.HaveRemove)
	return
}

// ArgsCreateConfig 创建配置参数
type ArgsCreateConfig struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//标题
	Title string `db:"title" json:"title" check:"title"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//描述文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//样式ID
	// 关联到样式库后，本记录的图片和文本将交给样式库布局实现
	StyleID int64 `db:"style_id" json:"styleID" check:"id" empty:"true"`
	//交易货币类型
	// 采用CoreCurrency匹配
	// 86 CNY
	Currency int `db:"currency" json:"currency" check:"currency"`
	//缴费金额
	Price int64 `db:"price" json:"price" check:"price"`
	//折扣前费用，用于展示
	PriceOld int64 `db:"price_old" json:"priceOld" check:"price"`
	//时间类型
	// 0 小时 1 天 2 周 3 月 4 年
	TimeType int `db:"time_type" json:"timeType" check:"intThan0" empty:"true"`
	//时间长度
	TimeN int `db:"time_n" json:"timeN" check:"intThan0" empty:"true"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateConfig 创建配置
func CreateConfig(args *ArgsCreateConfig) (data FieldsConfig, err error) {
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "nursing_info_pay_config", "INSERT INTO nursing_info_pay_config (org_id, title, des, cover_file_id, des_files, style_id, currency, price, price_old, time_type, time_n, params) VALUES (:org_id,:title,:des,:cover_file_id,:des_files,:style_id,:currency,:price,:price_old,:time_type,:time_n,:params)", args, &data)
	return
}

// ArgsUpdateConfig 修改配置参数
type ArgsUpdateConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//标题
	Title string `db:"title" json:"title" check:"title"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//描述文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//样式ID
	// 关联到样式库后，本记录的图片和文本将交给样式库布局实现
	StyleID int64 `db:"style_id" json:"styleID" check:"id" empty:"true"`
	//交易货币类型
	// 采用CoreCurrency匹配
	// 86 CNY
	Currency int `db:"currency" json:"currency" check:"currency"`
	//缴费金额
	Price int64 `db:"price" json:"price" check:"price"`
	//折扣前费用，用于展示
	PriceOld int64 `db:"price_old" json:"priceOld" check:"price"`
	//时间类型
	// 0 小时 1 天 2 周 3 月 4 年
	TimeType int `db:"time_type" json:"timeType" check:"intThan0" empty:"true"`
	//时间长度
	TimeN int `db:"time_n" json:"timeN" check:"intThan0" empty:"true"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateConfig 修改配置
func UpdateConfig(args *ArgsUpdateConfig) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE nursing_info_pay_config SET update_at = NOW(), title = :title, des = :des, cover_file_id = :cover_file_id, des_files = :des_files, style_id = :style_id, currency = :currency, price = :price, price_old = :price_old, time_type = :time_type, time_n = :time_n, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// ArgsDeleteConfig 删除配置参数
type ArgsDeleteConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteConfig 删除配置
func DeleteConfig(args *ArgsDeleteConfig) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "nursing_info_pay_config", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}
