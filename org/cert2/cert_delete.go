package OrgCert2

import (
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsDeleteCert 删除证书参数
type ArgsDeleteCert struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteCert 删除证书
func DeleteCert(args *ArgsDeleteCert) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "org_cert2", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	deleteCertCache(args.ID)
	return
}

// ArgsDeleteCerts 批量删除证书参数
type ArgsDeleteCerts struct {
	//ID
	IDs pq.Int64Array `db:"ids" json:"ids" check:"ids"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteCerts 批量删除证书
func DeleteCerts(args *ArgsDeleteCerts) (err error) {
	for _, v := range args.IDs {
		err = DeleteCert(&ArgsDeleteCert{
			ID:    v,
			OrgID: args.OrgID,
		})
		if err != nil {
			return
		}
	}
	return
}

// DeleteAllCertByBindID 删除指定用户的所有证件
func DeleteAllCertByBindID(bindFrom string, bindID int64) (err error) {
	var dataList []FieldsCert
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT c.id as id FROM org_cert2 as c, org_cert_config2 as g WHERE g.bind_from = $1 AND c.bind_id = $2 AND g.id = c.config_id", bindFrom, bindID)
	if err != nil || len(dataList) < 1 {
		err = nil
		return
	}
	for _, vCert := range dataList {
		_ = DeleteCert(&ArgsDeleteCert{
			ID:    vCert.ID,
			OrgID: -1,
		})
	}
	return
}
