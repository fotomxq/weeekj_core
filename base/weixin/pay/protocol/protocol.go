package BaseWeixinPayProtocol

import (
	"errors"
	"fmt"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	FinancePay "gitee.com/weeekj/weeekj_core/v5/finance/pay"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

//本方法集合用于构建自动扣费、签约操作
/**
1\ 在小程序或其他端发起签约后，系统将保存签约ID
2\ 服务端巡逻检查所有签约数据，在24小时到达之前自动触发续约请求
3\ 如果发现续约被第三方解除，将自动删除扣费
*/

//获取签约情况

// ArgsGetProtocolList 获取签约列表参数
type ArgsGetProtocolList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//签约模版ID
	TemplateID int64 `db:"template_id" json:"templateID" check:"id" empty:"true"`
	//签约模块
	// 0 会员模块; 1 平台组织会员模块
	ConfigSystem int `db:"config_system" json:"configSystem" check:"intThan0" empty:"true"`
	//签约模块ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id" empty:"true"`
}

// GetProtocolList 获取签约列表
func GetProtocolList(args *ArgsGetProtocolList) (dataList []FieldsProtocol, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.OrgID > -1 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.UserID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.TemplateID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "template_id = :template_id"
		maps["template_id"] = args.TemplateID
	}
	if args.ConfigSystem > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "config_system = :config_system"
		maps["config_system"] = args.ConfigSystem
		where = where + " AND config_id = :config_id"
		maps["config_id"] = args.ConfigID
	}
	if where != "" {
		where = "true"
	}
	tableName := "core_weixin_pay_protocol_template"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, org_id, user_id, template_id, config_system, config_id, next_at FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// ArgsCreateProtocol 记录新的协议参数
type ArgsCreateProtocol struct {
	//组织ID
	// 也是商户ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//签约模版ID
	TemplateID int64 `db:"template_id" json:"templateID" check:"id"`
	//签约模块
	// 0 会员模块; 1 平台组织会员模块
	ConfigSystem int `db:"config_system" json:"configSystem" check:"intThan0" empty:"true"`
	//签约模块ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
}

// CreateProtocol 记录新的协议
func CreateProtocol(args *ArgsCreateProtocol) (data FieldsProtocol, err error) {
	var templateData FieldsTemplate
	err = Router2SystemConfig.MainDB.Get(&templateData, "SELECT id FROM core_weixin_pay_protocol_template WHERE id = $1 AND ($2 < 1 OR org_id = $2)", args.TemplateID, args.OrgID)
	if err != nil || templateData.ID < 1 {
		err = errors.New(fmt.Sprint("template not exist, ", err))
		return
	}
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "core_weixin_pay_protocol", "INSERT INTO core_weixin_pay_protocol (org_id, user_id, template_id, config_system, config_id, next_at) VALUES (:org_id, :user_id, :template_id, :config_system, :config_id, NOW())", args, &data)
	return
}

// ArgsDeleteProtocolByConfig 为会员配置解约参数
type ArgsDeleteProtocolByConfig struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 也是商户ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//签约模块
	// 0 会员模块; 1 平台组织会员模块
	ConfigSystem int `db:"config_system" json:"configSystem" check:"intThan0" empty:"true"`
	//签约模块ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id"`
}

// DeleteProtocolByConfig 为会员配置解约
func DeleteProtocolByConfig(args *ArgsDeleteProtocolByConfig) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_weixin_pay_protocol", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// ArgsDeleteProtocol 解约指定ID参数
type ArgsDeleteProtocol struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 也是商户ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteProtocol 解约指定ID
func DeleteProtocol(args *ArgsDeleteProtocol) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "core_weixin_pay_protocol", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// ArgsCreatePay 发起支付请求参数
// 请注意，请在24小时之外发起该请求，否则可能无法前街续订
type ArgsCreatePay struct {
}

// CreatePay 发起支付请求
func CreatePay(args *ArgsCreatePay) (data FinancePay.FieldsPayType, err error) {
	return
}
