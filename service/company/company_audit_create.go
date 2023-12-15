package ServiceCompany

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
)

type ArgsCreateCompanyAudit struct {
	//hash值
	// 唯一的数据，可用于查询对应组织，取代ID直接查询；或用于第三方系统同步数据处理用
	Hash string `db:"hash" json:"hash"`
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
}

func CreateCompanyAudit(args *ArgsCreateCompanyAudit) (errCode string, err error) {
	if !checkCompanyUseType(args.UseType) {
		errCode = "err_erp_company_use_type"
		err = errors.New("not support use type")
		return
	}
	findData := GetCompanyBySN(args.OrgID, args.SN, args.UseType)
	if findData.ID > 0 {
		errCode = "err_erp_company_sn_replace"
		err = errors.New("company sn replace")
		return
	}
	if args.Hash == "" {
		args.Hash = CoreFilter.GetSha1Str(fmt.Sprint("org:", args.OrgID, ";name:", args.Name, ";sn:", args.SN))
		if args.Hash == "" {
			errCode = "err_hash"
			err = errors.New("hash error")
			return
		}
	}
	findHashData, _ := GetCompanyByHash(args.Hash, args.OrgID)
	if findHashData.ID > 0 {
		errCode = "err_erp_company_sn_replace"
		err = errors.New("hash replace")
		return
	}
	_, err = CoreSQL.CreateOneAndID(Router2SystemConfig.MainDB.DB, "INSERT INTO service_company_audit (hash, org_id, use_type, bind_org_id, bind_user_id, name, sn, des, country, city, address, map_type, longitude, latitude, params) VALUES (:hash, :org_id, :use_type, :bind_org_id, :bind_user_id, :name, :sn, :des, :country, :city, :address, :map_type, :longitude, :latitude, :params)", args)
	if err != nil {
		errCode = "err_insert"
		return
	}
	return
}
