package OrgCert2

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"time"
)

// ArgsUpdateCert 修改证书信息参数
type ArgsUpdateCert struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//绑定ID
	// 根据配置决定，可能是用户ID\组织ID\或其他任意主体
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//证件序列号
	SN string `db:"sn" json:"sn"`
	//过期时间
	ExpireAt string `db:"expire_at" json:"expireAt" check:"defaultTime" empty:"true"`
	//拍照文件ID序列
	FileIDs pq.Int64Array `db:"file_ids" json:"fileIDs" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateCert 修改证书信息参数
func UpdateCert(args *ArgsUpdateCert) (err error) {
	//修正参数
	if args.Params == nil {
		args.Params = CoreSQLConfig.FieldsConfigsType{}
	}
	//获取证件
	var certData FieldsCert
	certData, err = GetCert(&ArgsGetCert{
		ID:    args.ID,
		OrgID: args.OrgID,
	})
	if err != nil {
		return
	}
	//获取证件配置
	configData := getConfigByID(certData.ConfigID)
	if configData.ID < 1 {
		err = errors.New("no data")
		return
	}
	//检查SN长度
	if len(args.SN) > configData.SNLen {
		err = errors.New("sn too more")
		return
	}
	//检查过期时间
	var expireAt time.Time
	if args.ExpireAt != "" {
		expireAt, err = CoreFilter.GetTimeByDefault(args.ExpireAt)
		if err != nil {
			return
		}
	} else {
		expireAt, err = CoreFilter.GetTimeByAdd(configData.DefaultExpire)
		if err != nil {
			return
		}
	}
	if args.OrgID < 1 && args.BindID > 0 {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_cert2 SET update_at = NOW(), audit_at = to_timestamp(0), name = :name, sn = :sn, expire_at = :expire_at, file_ids = :file_ids, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id) AND (:bind_id < 1 OR bind_id = :bind_id)", map[string]interface{}{
			"id":        args.ID,
			"org_id":    args.OrgID,
			"bind_id":   args.BindID,
			"name":      args.Name,
			"sn":        args.SN,
			"expire_at": expireAt,
			"file_ids":  args.FileIDs,
			"params":    args.Params,
		})
		if err != nil {
			err = errors.New(fmt.Sprint("update cert id: ", args.ID, ", expire at: ", expireAt, ", org id and bind id mix, org id: ", args.OrgID, ", bind id: ", args.BindID, ", err: ", err))
			return
		}
	} else {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_cert2 SET update_at = NOW(), name = :name, sn = :sn, expire_at = :expire_at, file_ids = :file_ids, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id) AND (:bind_id < 1 OR bind_id = :bind_id)", map[string]interface{}{
			"id":        args.ID,
			"org_id":    args.OrgID,
			"bind_id":   args.BindID,
			"name":      args.Name,
			"sn":        args.SN,
			"expire_at": expireAt,
			"file_ids":  args.FileIDs,
			"params":    args.Params,
		})
		if err != nil {
			err = errors.New(fmt.Sprint("update cert id: ", args.ID, ", expire at: ", expireAt, ", err: ", err))
			return
		}
	}
	deleteCertCache(args.ID)
	return
}

// UpdateCertOrg 移动证件归属
func UpdateCertOrg(id int64, orgID int64) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_cert2 SET update_at = NOW(), org_id = :org_id WHERE id = :id", map[string]interface{}{
		"id":     id,
		"org_id": orgID,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("update cert org by id: ", id, ", org id: ", orgID, ", err: ", err))
		return
	}
	deleteCertCache(id)
	return
}

// ArgsUpdateCertExpire 修改证件过期时间参数
type ArgsUpdateCertExpire struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//过期时间
	ExpireAt string `db:"expire_at" json:"expireAt" check:"defaultTime" empty:"true"`
}

// UpdateCertExpire 修改证件过期时间
func UpdateCertExpire(args *ArgsUpdateCertExpire) (err error) {
	//获取过期时间
	var expireAt time.Time
	if args.ExpireAt != "" {
		expireAt, err = CoreFilter.GetTimeByISO(args.ExpireAt)
		if err != nil {
			return
		}
	}
	//更新数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_cert2 SET update_at = NOW(), expire_at = :expire_at WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", map[string]interface{}{
		"id":        args.ID,
		"org_id":    args.OrgID,
		"expire_at": expireAt,
	})
	if err != nil {
		return
	}
	//删除缓冲
	deleteCertCache(args.ID)
	//如果存在过期时间，则通知处理
	checkCertExpireSend(args.ID, expireAt)
	//反馈
	return
}
