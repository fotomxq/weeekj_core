package ServiceCompany

import (
	"errors"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
)

type ArgsUpdateCompanyAudit struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否审核通过
	IsAudit bool `json:"isAudit" check:"bool"`
}

func UpdateCompanyAudit(args *ArgsUpdateCompanyAudit) (errCode string, err error) {
	data := getCompanyAudit(args.ID)
	if data.ID < 1 || CoreSQL.CheckTimeHaveData(data.DeleteAt) || !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		errCode = "err_no_data"
		err = errors.New("no data")
		return
	}
	if CoreSQL.CheckTimeHaveData(data.AuditAt) {
		errCode = "err_no_data"
		err = errors.New("data have audit")
		return
	}
	if args.IsAudit {
		findCompanyData := GetCompanyBySN(args.OrgID, data.SN, data.UseType)
		if findCompanyData.ID > 0 {
			_, err = SetBind(&ArgsSetBind{
				OrgID:      data.OrgID,
				UserID:     data.BindUserID,
				NationCode: "",
				Phone:      "",
				CompanyID:  findCompanyData.ID,
				Managers:   []string{},
			})
			if err != nil {
				errCode = "err_update"
				return
			}
		} else {
			errCode, err = CreateCompany(&ArgsCreateCompany{
				Hash:          data.Hash,
				OrgID:         data.OrgID,
				UseType:       data.UseType,
				BindOrgID:     data.BindOrgID,
				BindUserID:    data.BindUserID,
				Name:          data.Name,
				SN:            data.SN,
				Des:           data.Des,
				Country:       data.Country,
				City:          data.City,
				Address:       data.Address,
				MapType:       data.MapType,
				Longitude:     data.Longitude,
				Latitude:      data.Latitude,
				Params:        data.Params,
				NoNameReplace: false,
			})
			if err != nil {
				return
			}
		}
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE service_company_audit SET audit_at = NOW() WHERE id = :id", map[string]interface{}{
			"id": data.ID,
		})
		if err != nil {
			errCode = "err_update"
			return
		}
	} else {
		_, err = CoreSQL.DeleteOneSoft(Router2SystemConfig.MainDB.DB, "service_company_audit", "id", map[string]interface{}{
			"id": data.ID,
		})
		if err != nil {
			errCode = "err_update"
			return
		}
	}
	//删除缓冲
	deleteCompanyAuditCache(data.ID)
	//反馈
	return
}
