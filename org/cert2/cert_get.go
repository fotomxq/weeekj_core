package OrgCert2

import (
	"errors"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	CoreSQLTime2 "github.com/fotomxq/weeekj_core/v5/core/sql/time2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetCertList 获取证书列表参数
type ArgsGetCertList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id" empty:"true"`
	//配置标识码
	ConfigMark string `json:"configMark" check:"mark" empty:"true"`
	//审核人
	AuditBindID int64 `db:"audit_bind_id" json:"auditBindID" check:"id" empty:"true"`
	//绑定ID
	// 根据配置决定，可能是用户ID\组织ID\或其他任意主体
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//是否过期
	NeedIsExpire bool `json:"needIsExpire" check:"bool"`
	IsExpire     bool `json:"isExpire" check:"bool"`
	//是否审核
	NeedIsAudit bool `json:"needIsAudit" check:"bool"`
	IsAudit     bool `json:"isAudit" check:"bool"`
	//是否缴费
	NeedIsPay bool `json:"needIsPay" check:"bool"`
	IsPay     bool `json:"isPay" check:"bool"`
	//是否被删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetCertList 获取证书列表
func GetCertList(args *ArgsGetCertList) (dataList []FieldsCert, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.ConfigID > -1 {
		where = where + " AND config_id = :config_id"
		maps["config_id"] = args.ConfigID
	} else {
		if args.ConfigMark != "" {
			configData := getConfigByMark(args.OrgID, args.ConfigMark)
			if configData.ID < 1 {
				configData = getConfigByMark(0, args.ConfigMark)
			}
			if configData.ID > 0 {
				where = where + " AND config_id = :config_id"
				maps["config_id"] = configData.ID
			}
		}
	}
	if args.AuditBindID > -1 {
		where = where + " AND audit_bind_id = :audit_bind_id"
		maps["audit_bind_id"] = args.AuditBindID
	}
	if args.BindID > -1 {
		where = where + " AND bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	if args.NeedIsExpire {
		if args.IsExpire {
			where = where + " AND expire_at < NOW()"
		} else {
			where = where + " AND expire_at >= NOW()"
		}
	}
	if args.NeedIsAudit {
		if args.IsAudit {
			where = where + " AND audit_at >= to_timestamp(1000000)"
		} else {
			where = where + " AND audit_at < to_timestamp(1000000)"
		}
	}
	if args.NeedIsPay {
		if args.IsPay {
			where = where + " AND pay_at >= to_timestamp(1000000)"
		} else {
			where = where + " AND pay_at < to_timestamp(1000000)"
		}
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "org_cert2"
	var rawList []FieldsCert
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		"SELECT id FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "expire_at", "audit_at"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getCertByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetCert 获取证书参数
type ArgsGetCert struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// GetCert 获取证书
func GetCert(args *ArgsGetCert) (data FieldsCert, err error) {
	data = getCertByID(args.ID)
	if data.ID < 1 || !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		err = errors.New("no data")
		return
	}
	return
}

// GetCertByBind 获取指定标识码和渠道的证件
func GetCertByBind(configOrgID int64, configMark string, bindFrom string, bindID int64) (data FieldsCert) {
	//验证参数
	if bindID < 1 {
		return
	}
	//获取配置
	configData := getConfigByMark(configOrgID, configMark)
	if configData.ID < 1 {
		return
	}
	if configData.BindFrom != bindFrom {
		return
	}
	//获取证件
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_cert2 WHERE config_id = $1 AND bind_id = $2 AND delete_at < to_timestamp(1000000) LIMIT 1", configData.ID, bindID)
	if err != nil || data.ID < 1 {
		return
	}
	data = getCertByID(data.ID)
	if data.ID < 1 {
		return
	}
	return
}

// ArgsGetCertMoreByBind 获取多个证件参数
type ArgsGetCertMoreByBind struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
	//标识码
	// 用于程序化识别处理机制
	Mark string `db:"mark" json:"mark" check:"mark"`
}

// GetCertMoreByBind 获取多个证件
func GetCertMoreByBind(args *ArgsGetCertMoreByBind) (dataList []FieldsCert, err error) {
	var configData FieldsConfig
	configData, err = GetConfigByMark(&ArgsGetConfigByMark{
		Mark:  args.Mark,
		OrgID: -1,
	})
	if err != nil {
		return
	}
	var rawList []FieldsCert
	err = Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM org_cert2 WHERE bind_id = ANY($1) AND ($2 < 1 OR org_id = $2) AND ($3 = true OR delete_at < to_timestamp(1000000)) AND config_id = $4 LIMIT 1000", args.IDs, args.OrgID, args.HaveRemove, configData.ID)
	if err != nil || len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range rawList {
		vData := getCertByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetCertMore 获取一组证件参数
type ArgsGetCertMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
}

// GetCertMore 获取一组证件
func GetCertMore(args *ArgsGetCertMore) (dataList []FieldsCert, err error) {
	for _, v := range args.IDs {
		vData := getCertByID(v)
		if vData.ID < 1 || !CoreFilter.EqID2(args.OrgID, vData.OrgID) || (!args.HaveRemove && !CoreSQL.CheckTimeHaveData(vData.DeleteAt)) {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetCertCount 获取证件数量参数
type ArgsGetCertCount struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//标识码
	// 用于程序化识别处理机制
	Mark string `db:"mark" json:"mark" check:"mark"`
	//查询时间范围
	TimeBetween CoreSQLTime2.DataCoreTime `json:"timeBetween"`
}

// GetCertCount 获取证件数量
func GetCertCount(args *ArgsGetCertCount) (count int64, err error) {
	//获取配置
	var configData FieldsConfig
	configData = getConfigByMark(args.OrgID, args.Mark)
	if configData.ID < 1 {
		configData = getConfigByMark(0, args.Mark)
	}
	if configData.ID < 1 {
		err = errors.New("no config")
		return
	}
	//获取数量
	err = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) as count FROM org_cert2 WHERE config_id = $1 AND ($2 < 1 OR org_id = $2) AND delete_at < to_timestamp(1000000) AND audit_at >= to_timestamp(1000000) AND ($3 < to_timestamp(1000000) OR ($3 >= to_timestamp(1000000) AND create_at >= $3)) AND ($4 < to_timestamp(1000000) OR ($4 >= to_timestamp(1000000) OR create_at <= $4))", configData.ID, args.OrgID, args.TimeBetween.MinTime, args.TimeBetween.MaxTime)
	if err != nil {
		return
	}
	return
}

// GetCertCountNoErr 获取证件数量
func GetCertCountNoErr(args *ArgsGetCertCount) (count int64) {
	count, _ = GetCertCount(args)
	return
}

// GetCertCountExpire 获取证件已经过期数量
func GetCertCountExpire(args *ArgsGetCertCount) (count int64, err error) {
	//获取配置
	var configData FieldsConfig
	configData = getConfigByMark(args.OrgID, args.Mark)
	if configData.ID < 1 {
		configData = getConfigByMark(0, args.Mark)
	}
	if configData.ID < 1 {
		err = errors.New("no config")
		return
	}
	//获取数量
	err = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) as count FROM org_cert2 WHERE config_id = $1 AND ($2 < 1 OR org_id = $2) AND delete_at < to_timestamp(1000000) AND audit_at >= to_timestamp(1000000) AND ($3 < to_timestamp(1000000) OR ($3 >= to_timestamp(1000000) AND create_at >= $3)) AND ($4 < to_timestamp(1000000) OR ($4 >= to_timestamp(1000000) AND create_at <= $4)) AND expire_at < NOW()", configData.ID, args.OrgID, args.TimeBetween.MinTime, args.TimeBetween.MaxTime)
	if err != nil {
		return
	}
	return
}

// 获取ID
func getCertByID(id int64) (data FieldsCert) {
	cacheMark := getCertCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, expire_at, audit_at, audit_bind_id, audit_ban_at, audit_des, org_id, config_id, bind_id, name, sn, file_ids, pay_at, pay_failed, pay_id, currency, price, params FROM org_cert2 WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 21600)
	return
}
