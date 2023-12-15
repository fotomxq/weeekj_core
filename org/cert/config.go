package OrgCert

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLIDs "github.com/fotomxq/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetConfigList 获取证件配置列表参数
type ArgsGetConfigList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//绑定来源
	// user 用户 / org 商户 / org_bind 商户成员 / finance_assets 财务资产 /
	BindFrom string `db:"bind_from" json:"bindFrom" check:"mark" empty:"true"`
	//标识码
	// 用于程序化识别处理机制
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//是否被删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetConfigList 获取证件配置列表
func GetConfigList(args *ArgsGetConfigList) (dataList []FieldsConfig, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.BindFrom != "" {
		where = where + " AND bind_from = :bind_from"
		maps["bind_from"] = args.BindFrom
	}
	if args.Mark != "" {
		where = where + " AND mark = :mark"
		maps["mark"] = args.Mark
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "org_cert_config"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, default_expire, org_id, bind_from, mark, name, des, cover_file_id, des_files, audit_type, currency, price, sn_len, tip_type, params FROM "+tableName+" WHERE "+where,
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
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, default_expire, org_id, bind_from, mark, name, des, cover_file_id, des_files, audit_type, currency, price, sn_len, tip_type, params FROM org_cert_config WHERE id = $1 AND ($2 < 1 OR org_id = $2) AND delete_at < to_timestamp(1000000)", args.ID, args.OrgID)
	return
}

// ArgsGetConfigByMark 获取指定配置Mark参数
type ArgsGetConfigByMark struct {
	//mark
	Mark string `db:"mark" json:"mark" check:"mark"`
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetConfigByMark 获取指定配置Mark
func GetConfigByMark(args *ArgsGetConfigByMark) (data FieldsConfig, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, default_expire, org_id, bind_from, mark, name, des, cover_file_id, des_files, audit_type, currency, price, sn_len, tip_type, params FROM org_cert_config WHERE mark = $1 AND ($2 < 1 OR org_id = $2) AND delete_at < to_timestamp(1000000)", args.Mark, args.OrgID)
	return
}

// ArgsGetConfigMore 获取一组配置参数
type ArgsGetConfigMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
}

// GetConfigMore 获取一组配置
func GetConfigMore(args *ArgsGetConfigMore) (dataList []FieldsConfig, err error) {
	err = CoreSQLIDs.GetIDsAndDelete(&dataList, "org_cert_config", "id, create_at, update_at, delete_at, default_expire, org_id, bind_from, mark, name, des, cover_file_id, des_files, audit_type, currency, price, sn_len, tip_type, params", args.IDs, args.HaveRemove)
	return
}

// GetConfigMoreMap 获取一组配置名称组
func GetConfigMoreMap(args *ArgsGetConfigMore) (data map[int64]string, err error) {
	data, err = CoreSQLIDs.GetIDsNameAndDelete("org_cert_config", args.IDs, args.HaveRemove)
	return
}

// ArgsCreateConfig 创建新的配置参数
type ArgsCreateConfig struct {
	//默认过期时间长度
	DefaultExpire string `db:"default_expire" json:"defaultExpire"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//绑定来源
	// user 用户 / org 商户 / finance_assets 财务资产 /
	BindFrom string `db:"bind_from" json:"bindFrom" check:"mark"`
	//标识码
	// 用于程序化识别处理机制
	Mark string `db:"mark" json:"mark" check:"mark"`
	//证件名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//描述文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//审核模式
	// none 无需审核; wait 人工审核; auto 自动审核(依赖其他模块，根据扩展参数具体识别方案);
	AuditType string `db:"audit_type" json:"auditType" check:"mark"`
	//审核费用
	// 如果为0则无效
	Currency int   `db:"currency" json:"currency" check:"currency" empty:"true"`
	Price    int64 `db:"price" json:"price" check:"price" empty:"true"`
	//序列号长度
	// 0则不会限制，但数据表最多存储300位，超出需使用扩展参数存储
	SNLen int `db:"sn_len" json:"snLen" check:"intThan0" empty:"true"`
	//通知类型
	// none 无通知; audit 审核通过后通知; expire 过期前通知; all 全部通知;
	TipType string `db:"tip_type" json:"tipType" check:"mark" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateConfig 创建新的配置
func CreateConfig(args *ArgsCreateConfig) (data FieldsConfig, err error) {
	//核对auditType
	if args.AuditType == "" {
		args.AuditType = "none"
	}
	if err = checkAuditType(args.AuditType); err != nil {
		return
	}
	//核对TipType
	if args.TipType == "" {
		args.TipType = "none"
	}
	if err = checkTipType(args.TipType); err != nil {
		return
	}
	//检查mark
	if args.Mark != "" {
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_cert_config WHERE org_id = $1 AND mark = $2 AND delete_at < to_timestamp(1000000)", args.OrgID, args.Mark)
		if err == nil && data.ID > 0 {
			err = errors.New("config mark is exist")
			return
		}
	}
	//创建数据
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "org_cert_config", "INSERT INTO org_cert_config (default_expire, org_id, bind_from, mark, name, des, cover_file_id, des_files, audit_type, currency, price, sn_len, tip_type, params) VALUES (:default_expire,:org_id,:bind_from,:mark,:name,:des,:cover_file_id,:des_files,:audit_type,:currency,:price,:sn_len,:tip_type,:params)", args, &data)
	if err != nil {
		return
	}
	return
}

// ArgsUpdateConfig 修改配置参数
type ArgsUpdateConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//绑定来源
	// user 用户 / org 商户 / finance_assets 财务资产 /
	BindFrom string `db:"bind_from" json:"bindFrom" check:"mark"`
	//默认过期时间长度
	DefaultExpire string `db:"default_expire" json:"defaultExpire"`
	//标识码
	// 用于程序化识别处理机制
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//证件名称
	Name string `db:"name" json:"name" check:"name"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//描述文件组
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//审核模式
	// none 无需审核; wait 人工审核; auto 自动审核(依赖其他模块，根据扩展参数具体识别方案);
	AuditType string `db:"audit_type" json:"auditType" check:"mark" empty:"true"`
	//审核费用
	// 如果为0则无效
	Currency int   `db:"currency" json:"currency" check:"currency" empty:"true"`
	Price    int64 `db:"price" json:"price" check:"price" empty:"true"`
	//序列号长度
	// 0则不会限制，但数据表最多存储300位，超出需使用扩展参数存储
	SNLen int `db:"sn_len" json:"snLen" check:"intThan0" empty:"true"`
	//通知类型
	// none 无通知; audit 审核通过后通知; expire 过期前通知; all 全部通知;
	TipType string `db:"tip_type" json:"tipType" check:"mark" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateConfig 修改配置
func UpdateConfig(args *ArgsUpdateConfig) (err error) {
	//核对auditType
	if args.AuditType == "" {
		args.AuditType = "none"
	}
	if err = checkAuditType(args.AuditType); err != nil {
		return
	}
	//更新的mark如果和之前不一致
	var data FieldsConfig
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_cert_config WHERE org_id = $1 AND id != $2 AND mark = $3 AND delete_at < to_timestamp(1000000)", args.OrgID, args.ID, args.Mark)
	if err == nil && data.ID > 0 {
		err = errors.New("mark is exist")
		return
	}
	//更新数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_cert_config SET update_at = NOW(), bind_from = :bind_from, default_expire = :default_expire, mark = :mark, name = :name, des = :des, cover_file_id = :cover_file_id, des_files = :des_files, audit_type = :audit_type, currency = :currency, price = :price, sn_len = :sn_len, tip_type = :tip_type, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
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

// DeleteConfig 删除配置
func DeleteConfig(args *ArgsDeleteConfig) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "org_cert_config", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	return
}

// 检查审核类型
func checkAuditType(t string) (err error) {
	switch t {
	case "none":
	case "wait":
	case "auto":
	default:
		err = errors.New("unknown audit type")
		return
	}
	return
}

// 检查提醒类型
func checkTipType(t string) (err error) {
	switch t {
	case "none":
	case "audit":
	case "expire":
	case "all":
	default:
		err = errors.New("unknown audit type")
		return
	}
	return
}
