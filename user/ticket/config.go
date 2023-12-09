package UserTicket

import (
	"errors"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLIDs "gitee.com/weeekj/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetConfigList 获取配置列表参数
type ArgsGetConfigList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否被删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetConfigList 获取配置列表
func GetConfigList(args *ArgsGetConfigList) (dataList []FieldsConfig, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(title ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"user_ticket_config",
		"id",
		"SELECT id, create_at, update_at, delete_at, default_expire_time, org_id, mark, title, des, cover_file_id, des_files, exemption_price, exemption_discount, exemption_min_price, exemption_min_price, use_order, limit_time_type, limit_count, style_id, params FROM user_ticket_config WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// ArgsGetConfigByID 获取指定配置ID参数
type ArgsGetConfigByID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetConfigByID 获取指定配置ID
func GetConfigByID(args *ArgsGetConfigByID) (data FieldsConfig, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, default_expire_time, org_id, mark, title, des, cover_file_id, des_files, exemption_price, exemption_discount, exemption_min_price, exemption_min_price, use_order, limit_time_type, limit_count, style_id, params FROM user_ticket_config WHERE id = $1 AND ($2 < 1 OR org_id = $2)", args.ID, args.OrgID)
	if err == nil && data.ID < 1 {
		err = errors.New("config not exist")
		return
	}
	return
}

// CheckConfigOrg 检查票据是否可用于
func CheckConfigOrg(orgID int64, configID int64) (err error) {
	var id int64
	err = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM user_ticket_config WHERE id = $1 AND org_id = $2 AND delete_at < to_timestamp(1000000)", id, orgID)
	if err != nil || id < 1 {
		err = errors.New("not find config")
		return
	}
	return
}

// ArgsGetConfigMore 获取一组配置参数
type ArgsGetConfigMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove"`
}

// GetConfigMore 获取一组配置
func GetConfigMore(args *ArgsGetConfigMore) (dataList []FieldsConfig, err error) {
	err = CoreSQLIDs.GetIDsAndDelete(&dataList, "user_ticket_config", "id, create_at, update_at, delete_at, default_expire_time, org_id, mark, title, des, cover_file_id, des_files, exemption_price, exemption_discount, exemption_min_price, exemption_min_price, use_order, limit_time_type, limit_count, style_id, params", args.IDs, args.HaveRemove)
	return
}

// GetConfigMoreMap 获取一组配置名称组
func GetConfigMoreMap(args *ArgsGetConfigMore) (data map[int64]string, err error) {
	data, err = CoreSQLIDs.GetIDsTitleAndDelete("user_ticket_config", args.IDs, args.HaveRemove)
	return
}

// GetConfigNameByID 获取配置名称
func GetConfigNameByID(id int64) string {
	if id < 1 {
		return ""
	}
	data := getConfigByID(id)
	if data.ID < 1 {
		return ""
	}
	return data.Title
}

// ArgsCreateConfig 创建新的配置参数
type ArgsCreateConfig struct {
	//默认过期时间
	DefaultExpireTime int64 `db:"default_expire_time" json:"defaultExpireTime"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//标识码
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//标题
	Title string `db:"title" json:"title" check:"title" min:"1" max:"300"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//描述文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//默认减免的费用比例、费用金额
	ExemptionPrice int64 `db:"exemption_price" json:"exemptionPrice"`
	// 1-100% 百分比
	ExemptionDiscount int64 `db:"exemption_discount" json:"exemptionDiscount"`
	//费用低于多少时，将失效
	// 依赖于订单的总金额判断
	ExemptionMinPrice int64 `db:"exemption_min_price" json:"exemptionMinPrice"`
	//是否可用于订单抵扣
	// 否则只能用于一件商品的抵扣
	UseOrder bool `db:"use_order" json:"useOrder" check:"bool"`
	//领取周期类型
	// 0 不限制; 1 一次性; 2 每天限制; 3 每周限制; 4 每月限制; 5 每季度限制; 6 每年限制
	LimitTimeType int `db:"limit_time_type" json:"limitTimeType" check:"intThan0" empty:"true"`
	//领取次数
	LimitCount int `db:"limit_count" json:"limitCount" check:"intThan0" empty:"true"`
	//样式ID
	// 关联到样式库后，本记录的图片和文本将交给样式库布局实现
	StyleID int64 `db:"style_id" json:"styleID" check:"id" empty:"true"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateConfig 创建新的配置
func CreateConfig(args *ArgsCreateConfig) (data FieldsConfig, err error) {
	//检查组织下的mark是否重复？
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM user_ticket_config WHERE org_id = $1 AND mark = $2 AND delete_at < to_timestamp(1000000)", args.OrgID, args.Mark)
	if err == nil && data.ID > 0 {
		err = errors.New("mark is exist")
		return
	}
	//创建数据
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "user_ticket_config", "INSERT INTO user_ticket_config (default_expire_time, org_id, mark, title, des, cover_file_id, des_files, exemption_price, exemption_discount, exemption_min_price, use_order, limit_time_type, limit_count, style_id, params) VALUES (:default_expire_time, :org_id, :mark, :title, :des, :cover_file_id, :des_files, :exemption_price, :exemption_discount, :exemption_min_price, :use_order, :limit_time_type, :limit_count, :style_id, :params)", args, &data)
	return
}

// ArgsUpdateConfig 修改配置参数
type ArgsUpdateConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//默认过期时间
	DefaultExpireTime int64 `db:"default_expire_time" json:"defaultExpireTime"`
	//标识码
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//标题
	Title string `db:"title" json:"title" check:"title" min:"1" max:"300"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//描述文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//默认减免的费用比例、费用金额
	ExemptionPrice int64 `db:"exemption_price" json:"exemptionPrice"`
	// 1-100% 百分比
	ExemptionDiscount int64 `db:"exemption_discount" json:"exemptionDiscount"`
	//费用低于多少时，将失效
	// 依赖于订单的总金额判断
	ExemptionMinPrice int64 `db:"exemption_min_price" json:"exemptionMinPrice"`
	//是否可用于订单抵扣
	// 否则只能用于一件商品的抵扣
	UseOrder bool `db:"use_order" json:"useOrder" check:"bool"`
	//领取周期类型
	// 0 不限制; 1 一次性; 2 每天限制; 3 每周限制; 4 每月限制; 5 每季度限制; 6 每年限制
	LimitTimeType int `db:"limit_time_type" json:"limitTimeType" check:"intThan0" empty:"true"`
	//领取次数
	LimitCount int `db:"limit_count" json:"limitCount" check:"intThan0" empty:"true"`
	//样式ID
	// 关联到样式库后，本记录的图片和文本将交给样式库布局实现
	StyleID int64 `db:"style_id" json:"styleID" check:"id" empty:"true"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateConfig 修改配置
func UpdateConfig(args *ArgsUpdateConfig) (err error) {
	//更新的mark如果和之前不一致
	var data FieldsConfig
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM user_ticket_config WHERE org_id = $1 AND id != $2 AND mark = $3 AND delete_at < to_timestamp(1000000)", args.OrgID, args.ID, args.Mark)
	if err == nil && data.ID > 0 {
		err = errors.New("mark is exist")
		return
	}
	//更新数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE user_ticket_config SET update_at = NOW(), default_expire_time = :default_expire_time, mark = :mark, title = :title, des = :des, cover_file_id = :cover_file_id, des_files = :des_files, exemption_price = :exemption_price, exemption_discount = :exemption_discount, exemption_min_price = :exemption_min_price, use_order = :use_order, limit_time_type = :limit_time_type, limit_count = :limit_count, style_id = :style_id, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// ArgsDeleteConfig 删除配置参数
type ArgsDeleteConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

func DeleteConfig(args *ArgsDeleteConfig) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "user_ticket_config", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err == nil {
		_ = ClearTicket(&ArgsClearTicket{
			ConfigD: args.ID,
			OrgID:   args.OrgID,
		})
	}
	return
}

// 获取指定的配置
func getConfigByID(id int64) (data FieldsConfig) {
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, default_expire_time, org_id, mark, title, des, cover_file_id, des_files, exemption_price, exemption_discount, exemption_min_price, exemption_min_price, use_order, limit_time_type, limit_count, style_id, params FROM user_ticket_config WHERE id = $1", id)
	if err != nil {
		return
	}
	return
}
