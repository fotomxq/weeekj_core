package ServiceCompany

import (
	"errors"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

// ArgsUpdateCompany 修改公司信息参数
type ArgsUpdateCompany struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//用途
	// client 客户; supplier 供应商; partners 合作商; service 服务商
	UseType string `db:"use_type" json:"useType"`
	//绑定组织ID
	BindOrgID int64 `db:"bind_org_id" json:"bindOrgID" check:"id" empty:"true"`
	//绑定用户ID
	BindUserID int64 `db:"bind_user_id" json:"bindUserID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//公司营业执照编号
	SN string `db:"sn" json:"sn"`
	//描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country" check:"country"`
	//所属城市
	City int `db:"city" json:"city" check:"city" empty:"true"`
	//街道详细信息
	Address string `db:"address" json:"address" check:"address" empty:"true"`
	//地图制式
	// 0 WGS-84 / 1 GCJ-02 / 2 BD-09
	MapType int `db:"map_type" json:"mapType" check:"mapType"`
	//坐标位置
	Longitude float64 `db:"longitude" json:"longitude"`
	Latitude  float64 `db:"latitude" json:"latitude"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
	//是否不排查名称
	NoNameReplace bool `json:"noNameReplace"`
}

// UpdateCompany 修改公司信息
func UpdateCompany(args *ArgsUpdateCompany) (errCode string, err error) {
	if !checkCompanyUseType(args.UseType) {
		errCode = "err_erp_company_use_type"
		err = errors.New("not support use type")
		return
	}
	if args.SN != "" {
		findData := GetCompanyBySN(args.OrgID, args.SN, args.UseType)
		if findData.ID > 0 && args.ID != findData.ID {
			errCode = "err_erp_company_sn_replace"
			err = errors.New("company sn replace")
			return
		}
	}
	if !args.NoNameReplace {
		findNameData := GetCompanyByName(args.OrgID, args.Name, args.UseType)
		if findNameData.ID > 0 && args.ID != findNameData.ID {
			errCode = "err_erp_company_name_replace"
			err = errors.New("company name replace")
			return
		}
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_company SET update_at = NOW(), use_type = :use_type, bind_org_id = :bind_org_id, bind_user_id = :bind_user_id, name = :name, sn = :sn, des = :des, country = :country, city = :city, address = :address, map_type = :map_type, longitude = :longitude, latitude = :latitude, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		errCode = "err_update"
		return
	}
	deleteCompanyCache(args.ID)
	return
}
