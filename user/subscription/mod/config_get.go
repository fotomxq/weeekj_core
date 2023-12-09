package UserSubscriptionMod

import (
	"errors"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

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
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, mark, time_type, time_n, currency, price, price_old, title, des, cover_file_id, des_files, user_groups, exemption_discount, exemption_price, exemption_min_price, limits, exemption_time, style_id, params FROM user_sub_config WHERE id = $1 AND ($2 < 1 OR org_id = $2)", args.ID, args.OrgID)
	return
}

// ArgsCheckConfigAndOrg 检查配置和商户是否关联参数
type ArgsCheckConfigAndOrg struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// CheckConfigAndOrg 检查配置和商户是否关联
func CheckConfigAndOrg(args *ArgsCheckConfigAndOrg) (err error) {
	var id int64
	err = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM user_sub_config WHERE id = $1 AND ($2 < 1 OR org_id = $2)", args.ID, args.OrgID)
	if err == nil && id < 1 {
		err = errors.New("not exist")
		return
	}
	return
}
