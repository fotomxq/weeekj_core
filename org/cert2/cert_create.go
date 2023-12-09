package OrgCert2

import (
	"errors"
	"fmt"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"time"
)

// ArgsCreateCert 创建新的证书请求参数
type ArgsCreateCert struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//绑定来源
	// user 用户 / org 商户 / org_bind 商户成员 / finance_assets 财务资产 /
	BindFrom string `db:"bind_from" json:"bindFrom" check:"mark"`
	//配置ID
	ConfigID   int64  `db:"config_id" json:"configID" check:"id" empty:"true"`
	ConfigMark string `db:"config_mark" json:"configMark" check:"mark" empty:"true"`
	//绑定ID
	// 根据配置决定，可能是用户ID\组织ID\或其他任意主体
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//证件序列号
	SN string `db:"sn" json:"sn"`
	//拍照文件ID序列
	FileIDs pq.Int64Array `db:"file_ids" json:"fileIDs" check:"ids" empty:"true"`
	//过期时间
	ExpireAt string `db:"expire_at" json:"expireAt" check:"defaultTime" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateCert 创建新的证书请求
func CreateCert(args *ArgsCreateCert) (data FieldsCert, errCode string, err error) {
	if args.Params == nil {
		args.Params = CoreSQLConfig.FieldsConfigsType{}
	}
	var configData FieldsConfig
	if args.ConfigID > 0 {
		configData = getConfigByID(args.ConfigID)
	} else {
		configData = getConfigByMark(args.OrgID, args.ConfigMark)
		if configData.ID < 1 {
			configData = getConfigByMark(0, args.ConfigMark)
		}
	}
	if configData.ID < 1 || CoreSQL.CheckTimeHaveData(configData.DeleteAt) {
		errCode = "err_config"
		err = errors.New(fmt.Sprint("get config data, ", err))
		return
	}
	if configData.BindFrom != args.BindFrom {
		errCode = "err_own"
		err = errors.New("bind from error")
		return
	}
	if len(args.SN) > configData.SNLen {
		errCode = "err_sn"
		err = errors.New("sn too more")
		return
	}
	if len(args.FileIDs) > 100 {
		errCode = "err_too_many"
		err = errors.New("too many files")
		return
	}
	var price int64
	price = configData.Price
	var payAt time.Time
	if price < 1 {
		payAt = CoreFilter.GetNowTime()
	}
	var expireAt time.Time
	if args.ExpireAt == "" {
		expireAt, err = CoreFilter.GetTimeByAdd(configData.DefaultExpire)
	} else {
		expireAt, err = CoreFilter.GetTimeByDefault(args.ExpireAt)
	}
	if err != nil {
		errCode = "err_time"
		return
	}
	var auditAt time.Time
	if configData.AuditType == "none" {
		auditAt = CoreFilter.GetNowTime()
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_cert2 WHERE bind_id = $1 AND org_id = $2 AND config_id = $3 AND delete_at < to_timestamp(1000000) LIMIT 1", args.BindID, args.OrgID, configData.ID)
	if err == nil && data.ID > 0 {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_cert2 SET update_at = NOW(), audit_at = :audit_at, name = :name, sn = :sn, expire_at = :expire_at, file_ids = :file_ids, params = :params WHERE id = :id", map[string]interface{}{
			"id":        data.ID,
			"audit_at":  auditAt,
			"name":      args.Name,
			"sn":        args.SN,
			"expire_at": expireAt,
			"file_ids":  args.FileIDs,
			"params":    args.Params,
		})
		if err != nil {
			errCode = "report_update_failed"
			err = errors.New(fmt.Sprint("update cert id: ", data.ID, ", expireAt: ", expireAt, ", err: ", err))
			return
		}
		//删除缓冲
		deleteCertCache(data.ID)
		//重新获取证件数据
		data = getCertByID(data.ID)
		if data.ID < 1 {
			errCode = "err_no_data"
			err = errors.New("no data")
			return
		}
	} else {
		err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "org_cert2", "INSERT INTO org_cert2 (expire_at, audit_at, audit_bind_id, audit_ban_at, audit_des, org_id, config_id, bind_id, name, sn, file_ids, pay_at, pay_failed, pay_id, currency, price, params) VALUES (:expire_at,:audit_at,:audit_bind_id,to_timestamp(0),:audit_des,:org_id,:config_id,:bind_id,:name,:sn,:file_ids,:pay_at,false,:pay_id,:currency,:price,:params)", map[string]interface{}{
			"expire_at":     expireAt,
			"audit_at":      auditAt,
			"audit_bind_id": 0,
			"audit_des":     "",
			"org_id":        args.OrgID,
			"config_id":     configData.ID,
			"bind_id":       args.BindID,
			"name":          args.Name,
			"sn":            args.SN,
			"file_ids":      args.FileIDs,
			"pay_at":        payAt,
			"pay_id":        0,
			"currency":      configData.Currency,
			"price":         price,
			"params":        args.Params,
		}, &data)
		if err != nil {
			errCode = "err_insert"
			return
		}
	}
	//如果不存在缴费需求，则触发自动审核
	if configData.AuditType == "auto" && CoreSQL.CheckTimeHaveData(payAt) {
		pushNatsAutoAudit(data.ID)
	}
	//过期处理
	checkCertExpireSend(data.ID, data.ExpireAt)
	//反馈
	return
}
