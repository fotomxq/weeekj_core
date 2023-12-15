package OrgCert

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLFrom "github.com/fotomxq/weeekj_core/v5/core/sql/from"
	CoreSQLIDs "github.com/fotomxq/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	CoreSQLTime "github.com/fotomxq/weeekj_core/v5/core/sql/time"
	FinancePay "github.com/fotomxq/weeekj_core/v5/finance/pay"
	FinancePayCreate "github.com/fotomxq/weeekj_core/v5/finance/pay_create"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
	"time"
)

// ArgsGetCertList 获取证书列表参数
type ArgsGetCertList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//子商户
	// 实际发起的商户ID，可以和顶级商户orgID一致
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID" check:"id" empty:"true"`
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id" empty:"true"`
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

// GetCertList 获取证书列表
func GetCertList(args *ArgsGetCertList) (dataList []FieldsCert, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQLField(args.IsRemove, where, "c.delete_at")
	if args.OrgID > -1 {
		where = where + " AND c.org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.ChildOrgID > -1 {
		where = where + " AND c.child_org_id = :child_org_id"
		maps["child_org_id"] = args.ChildOrgID
	}
	if args.ConfigID > 0 {
		where = where + " AND c.config_id = :config_id"
		maps["config_id"] = args.ConfigID
	}
	if args.AuditBindID > -1 {
		where = where + " AND c.audit_bind_id = :audit_bind_id"
		maps["audit_bind_id"] = args.AuditBindID
	}
	if args.BindID > -1 {
		where = where + " AND c.bind_id = :bind_id"
		maps["bind_id"] = args.BindID
	}
	if args.NeedIsAudit {
		if args.IsAudit {
			where = where + " AND c.audit_at >= to_timestamp(1000000)"
		} else {
			where = where + " AND c.audit_at < to_timestamp(1000000)"
		}
	}
	if args.NeedIsExpire {
		if args.IsExpire {
			where = where + " AND c.expire_at < NOW()"
		} else {
			where = where + " AND c.expire_at >= NOW()"
		}
	}
	maps["bind_from"] = args.BindFrom
	maps["mark"] = args.Mark
	if args.Search != "" {
		where = where + " AND (c.name ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	where = where + " AND c.config_id = d.id AND (:bind_from = '' OR d.bind_from = :bind_from) AND (:mark = '' OR d.mark = :mark)"
	tableName := "org_cert as c, org_cert_config as d"
	args.Pages.Sort = "c." + args.Pages.Sort
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"c.id",
		"SELECT c.id as id, c.create_at as create_at, c.update_at as update_at, c.delete_at as delete_at, c.expire_at as expire_at, c.audit_at as audit_at, c.audit_bind_id as audit_bind_id, c.audit_ban_at as audit_ban_at, c.audit_des as audit_des, c.org_id as org_id, c.child_org_id as child_org_id, c.config_id as config_id, c.bind_id as bind_id, c.name as name, c.sn as sn, c.file_ids as file_ids, c.pay_at as pay_at, c.pay_failed as pay_failed, c.pay_id as pay_id, c.currency as currency, c.price as price, c.params as params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"c.id", "c.create_at", "c.update_at", "c.delete_at", "c.expire_at", "c.audit_at"},
	)
	return
}

// ArgsGetCertChildGroupList 获取子公司统计列表参数
type ArgsGetCertChildGroupList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//配置ID
	ConfigID int64 `db:"config_id" json:"configID" check:"id" empty:"true"`
	//是否过期
	NeedIsExpire bool `json:"needIsExpire" check:"bool"`
	IsExpire     bool `json:"isExpire" check:"bool"`
	//是否审核
	NeedIsAudit bool `json:"needIsAudit" check:"bool"`
	IsAudit     bool `json:"isAudit" check:"bool"`
	//是否缴费
	NeedIsPay bool `json:"needIsPay" check:"bool"`
	IsPay     bool `json:"isPay" check:"bool"`
	//绑定来源
	// user 用户 / org 商户 / org_bind 商户成员 / finance_assets 财务资产 /
	BindFrom string `db:"bind_from" json:"bindFrom" check:"mark" empty:"true"`
	//标识码
	// 用于程序化识别处理机制
	Mark string `db:"mark" json:"mark" check:"mark" empty:"true"`
	//是否被删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool" empty:"true"`
}

type DataGetCertChildGroupList struct {
	//子商户
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID"`
	//证件数量
	Count int64 `db:"count" json:"count"`
}

// GetCertChildGroupList 获取子公司统计列表
func GetCertChildGroupList(args *ArgsGetCertChildGroupList) (dataList []DataGetCertChildGroupList, dataCount int64, err error) {
	args.Pages.Sort = fmt.Sprint("c.", args.Pages.Sort)
	where := "($1 < 0 OR c.org_id = $1) AND ($2 < 0 OR c.config_id = $2) AND ($3 = false OR (($4 = true AND c.expire_at < NOW()) OR ($4 = false AND c.expire_at >= NOW()))) AND ($5 = false OR (($6 = true AND c.audit_at < NOW()) OR ($6 = false AND c.audit_at >= NOW()))) AND ($7 = false OR (($8 = true AND c.pay_at < NOW()) OR ($8 = false AND c.pay_at >= NOW()))) AND ($9 = '' OR d.bind_from = $9) AND ($10 = '' OR d.mark = $10) AND (($11 = false AND c.delete_at <= to_timestamp(1000000)) OR ($11 = true AND c.delete_at <= to_timestamp(1000000))) AND c.config_id = d.id GROUP BY c.child_org_id"
	tableName := "org_cert as c, org_cert_config as d"
	dataCount, err = CoreSQL.GetListPageAndCountArgs(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"c.child_org_id",
		"SELECT c.child_org_id as child_org_id, COUNT(c.id) as count FROM "+tableName+" WHERE "+where,
		where,
		&args.Pages,
		[]string{"c.child_org_id"},
		args.OrgID,
		args.ConfigID,
		args.NeedIsExpire,
		args.IsExpire,
		args.NeedIsAudit,
		args.IsAudit,
		args.NeedIsPay,
		args.IsPay,
		args.BindFrom,
		args.Mark,
		args.IsRemove,
	)
	return
}

// ArgsGetCert 获取证书参数
type ArgsGetCert struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//子商户
	// 实际发起的商户ID，可以和顶级商户orgID一致
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID" check:"id" empty:"true"`
}

// GetCert 获取证书
func GetCert(args *ArgsGetCert) (data FieldsCert, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, expire_at, audit_at, audit_bind_id, audit_ban_at, audit_des, org_id, child_org_id, config_id, bind_id, name, sn, file_ids, pay_at, pay_failed, pay_id, currency, price, params FROM org_cert WHERE id = $1 AND ($2 < 1 OR org_id = $2) AND ($3 < 1 OR child_org_id = $3)", args.ID, args.OrgID, args.ChildOrgID)
	if err == nil && data.ID < 1 {
		err = errors.New("data is empty")
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
	//子商户
	// 实际发起的商户ID，可以和顶级商户orgID一致
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID" check:"id" empty:"true"`
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
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, update_at, delete_at, expire_at, audit_at, audit_bind_id, audit_ban_at, audit_des, org_id, child_org_id, config_id, bind_id, name, sn, file_ids, pay_at, pay_failed, pay_id, currency, price, params FROM org_cert WHERE bind_id = ANY($1) AND ($2 < 1 OR org_id = $2) AND ($3 < 1 OR child_org_id = $3) AND ($4 = true OR delete_at < to_timestamp(1000000)) AND config_id = $5 LIMIT 1000", args.IDs, args.OrgID, args.ChildOrgID, args.HaveRemove, configData.ID)
	if err == nil && len(dataList) < 1 {
		err = errors.New("data is empty")
		return
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
	err = CoreSQLIDs.GetIDsOrgAndDelete(&dataList, "org_cert", "id, create_at, update_at, delete_at, expire_at, audit_at, audit_bind_id, audit_ban_at, audit_des, org_id, child_org_id, config_id, bind_id, name, sn, file_ids, pay_at, pay_failed, pay_id, currency, price, params", args.IDs, args.OrgID, args.HaveRemove)
	return
}

// ArgsGetCertCount 获取证件数量参数
type ArgsGetCertCount struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//子商户
	// 实际发起的商户ID，可以和顶级商户orgID一致
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID" check:"id" empty:"true"`
	//标识码
	// 用于程序化识别处理机制
	Mark string `db:"mark" json:"mark" check:"mark"`
	//查询时间范围
	TimeBetween CoreSQLTime.DataCoreTime `json:"timeBetween"`
}

// GetCertCount 获取证件数量
func GetCertCount(args *ArgsGetCertCount) (count int64, err error) {
	//时间结构
	var timeBetween CoreSQLTime.FieldsCoreTime
	if args.TimeBetween.MinTime != "" || args.TimeBetween.MaxTime != "" {
		timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
		if err != nil {
			return
		}
	}
	//获取配置
	var configData FieldsConfig
	configData, err = GetConfigByMark(&ArgsGetConfigByMark{
		Mark:  args.Mark,
		OrgID: -1,
	})
	if err != nil {
		return
	}
	//获取数量
	err = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) as count FROM org_cert WHERE config_id = $1 AND ($2 < 1 OR org_id = $2) AND ($3 < 1 OR child_org_id = $3) AND delete_at < to_timestamp(1000000) AND audit_at >= to_timestamp(1000000) AND ($4 < to_timestamp(1000000) OR create_at >= $4) AND ($5 < to_timestamp(1000000) OR create_at <= $5)", configData.ID, args.OrgID, args.ChildOrgID, timeBetween.MinTime, timeBetween.MaxTime)
	return
}

// GetCertCountExpire 获取证件已经过期数量
func GetCertCountExpire(args *ArgsGetCertCount) (count int64, err error) {
	//时间结构
	var timeBetween CoreSQLTime.FieldsCoreTime
	if args.TimeBetween.MinTime != "" || args.TimeBetween.MaxTime != "" {
		timeBetween, err = CoreSQLTime.GetBetweenByISO(args.TimeBetween)
		if err != nil {
			return
		}
	}
	//获取配置
	var configData FieldsConfig
	configData, err = GetConfigByMark(&ArgsGetConfigByMark{
		Mark:  args.Mark,
		OrgID: -1,
	})
	if err != nil {
		return
	}
	//获取数量
	err = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) as count FROM org_cert WHERE config_id = $1 AND ($2 < 1 OR org_id = $2) AND ($3 < 1 OR child_org_id = $3) AND delete_at < to_timestamp(1000000) AND audit_at >= to_timestamp(1000000) AND ($4 < to_timestamp(1000000) OR create_at >= $4) AND ($5 < to_timestamp(1000000) OR create_at <= $5) AND expire_at < NOW()", configData.ID, args.OrgID, args.ChildOrgID, timeBetween.MinTime, timeBetween.MaxTime)
	return
}

// ArgsCheckCertByMarks 通过一组mark查询证件到期情况参数
type ArgsCheckCertByMarks struct {
	//配置组
	Marks pq.StringArray `db:"marks" json:"marks" check:"marks"`
	//组织ID
	// 可选
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//子商户
	// 实际发起的商户ID，可以和顶级商户orgID一致
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID" check:"id" empty:"true"`
	//绑定ID
	// 根据配置决定，可能是用户ID\组织ID\或其他任意主体
	BindID int64 `db:"bind_id" json:"bindID" check:"id"`
}

// DataCheckCertByMarks 通过一组mark查询证件到期情况
type DataCheckCertByMarks struct {
	Mark string `json:"mark"`
	IsOK bool   `json:"isOK"`
}

// CheckCertByMarks 通过一组mark查询证件到期情况
func CheckCertByMarks(args *ArgsCheckCertByMarks) (dataList []DataCheckCertByMarks, err error) {
	var configList []FieldsConfig
	err = Router2SystemConfig.MainDB.Select(&configList, "SELECT id, mark FROM org_cert_config WHERE mark = ANY($1) AND delete_at < to_timestamp(1000000)", args.Marks)
	if err != nil {
		return
	}
	var configIDs pq.Int64Array
	for _, v := range configList {
		configIDs = append(configIDs, v.ID)
	}
	var certList []FieldsCert
	err = Router2SystemConfig.MainDB.Select(&certList, "SELECT id, config_id FROM org_cert WHERE bind_id = $1 AND config_id = ANY($2) AND ($3 < 1 OR org_id = $3) AND ($4 < 1 OR child_org_id = $4) AND delete_at < to_timestamp(1000000) AND expire_at >= NOW() AND audit_at >= to_timestamp(1000000)", args.BindID, configIDs, args.OrgID, args.ChildOrgID)
	if err != nil {
		return
	}
	for _, v := range configList {
		isFind := false
		for _, v2 := range certList {
			if v.ID == v2.ConfigID {
				isFind = true
				dataList = append(dataList, DataCheckCertByMarks{
					Mark: v.Mark,
					IsOK: true,
				})
				break
			}
		}
		if !isFind {
			dataList = append(dataList, DataCheckCertByMarks{
				Mark: v.Mark,
				IsOK: false,
			})
		}
	}
	return
}

// ArgsCreateCert 创建新的证书请求参数
type ArgsCreateCert struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//子商户
	// 实际发起的商户ID，可以和顶级商户orgID一致
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID" check:"id" empty:"true"`
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
	ExpireAt string `db:"expire_at" json:"expireAt" check:"isoTime" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateCert 创建新的证书请求
func CreateCert(args *ArgsCreateCert) (data FieldsCert, errCode string, err error) {
	var configData FieldsConfig
	if args.ConfigID > 0 {
		configData, err = GetConfigByID(&ArgsGetConfigByID{
			ID:    args.ConfigID,
			OrgID: -1,
		})
	} else {
		configData, err = GetConfigByMark(&ArgsGetConfigByMark{
			Mark:  args.ConfigMark,
			OrgID: -1,
		})
	}
	if err != nil || configData.ID < 1 {
		errCode = "config_not_find"
		err = errors.New(fmt.Sprint("get config data, ", err))
		return
	}
	if configData.BindFrom != args.BindFrom {
		errCode = "bind_from_error"
		err = errors.New("bind from error")
		return
	}
	if len(args.SN) > configData.SNLen {
		errCode = "sn_error"
		err = errors.New("sn too more")
		return
	}
	if len(args.FileIDs) > 100 {
		errCode = "files_too_many"
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
		expireAt, err = CoreFilter.GetTimeByISO(args.ExpireAt)
	}
	if err != nil {
		errCode = "expire_error"
		return
	}
	var auditAt time.Time
	if configData.AuditType == "none" {
		auditAt = CoreFilter.GetNowTime()
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_cert WHERE bind_id = $1 AND org_id = $2 AND child_org_id = $3 AND config_id = $4 AND delete_at < to_timestamp(1000000)", args.BindID, args.OrgID, args.ChildOrgID, configData.ID)
	if err == nil && data.ID > 0 {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_cert SET update_at = NOW(), audit_at = :audit_at, name = :name, sn = :sn, expire_at = :expire_at, file_ids = :file_ids, params = :params WHERE id = :id", map[string]interface{}{
			"id":        data.ID,
			"audit_at":  auditAt,
			"name":      args.Name,
			"sn":        args.SN,
			"expire_at": expireAt,
			"file_ids":  args.FileIDs,
			"params":    args.Params,
		})
		if err != nil {
			errCode = "update_cert"
			err = errors.New(fmt.Sprint("update cert id: ", data.ID, ", expireAt: ", expireAt, ", err: ", err))
			return
		}
		data, err = GetCert(&ArgsGetCert{
			ID:         data.ID,
			OrgID:      -1,
			ChildOrgID: -1,
		})
		return
	}
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "org_cert", "INSERT INTO org_cert (expire_at, audit_at, audit_bind_id, audit_ban_at, audit_des, org_id, child_org_id, config_id, bind_id, name, sn, file_ids, pay_at, pay_failed, pay_id, currency, price, params) VALUES (:expire_at,:audit_at,:audit_bind_id,to_timestamp(0),:audit_des,:org_id,:child_org_id,:config_id,:bind_id,:name,:sn,:file_ids,:pay_at,false,:pay_id,:currency,:price,:params)", map[string]interface{}{
		"expire_at":     expireAt,
		"audit_at":      auditAt,
		"audit_bind_id": 0,
		"audit_des":     "",
		"org_id":        args.OrgID,
		"child_org_id":  args.ChildOrgID,
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
		errCode = "insert_cert"
	}
	return
}

// ArgsUpdateCert 修改证书信息参数
type ArgsUpdateCert struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//子商户
	// 实际发起的商户ID，可以和顶级商户orgID一致
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID" check:"id" empty:"true"`
	//绑定ID
	// 根据配置决定，可能是用户ID\组织ID\或其他任意主体
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//证件序列号
	SN string `db:"sn" json:"sn"`
	//过期时间
	ExpireAt string `db:"expire_at" json:"expireAt" check:"isoTime" empty:"true"`
	//拍照文件ID序列
	FileIDs pq.Int64Array `db:"file_ids" json:"fileIDs" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateCert 修改证书信息参数
func UpdateCert(args *ArgsUpdateCert) (err error) {
	var expireAt time.Time
	if args.ExpireAt != "" {
		expireAt, err = CoreFilter.GetTimeByISO(args.ExpireAt)
		if err != nil {
			return
		}
	} else {
		var certData FieldsCert
		certData, err = GetCert(&ArgsGetCert{
			ID:         args.ID,
			OrgID:      args.OrgID,
			ChildOrgID: args.ChildOrgID,
		})
		if err != nil {
			return
		}
		var configData FieldsConfig
		configData, err = GetConfigByID(&ArgsGetConfigByID{
			ID:    certData.ID,
			OrgID: certData.OrgID,
		})
		if err != nil {
			return
		}
		expireAt, err = CoreFilter.GetTimeByAdd(configData.DefaultExpire)
		if err != nil {
			return
		}
	}
	if args.OrgID < 1 && args.BindID > 0 {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_cert SET update_at = NOW(), audit_at = to_timestamp(0), name = :name, sn = :sn, expire_at = :expire_at, file_ids = :file_ids, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id) AND (:bind_id < 1 OR bind_id = :bind_id) AND (:child_org_id < 1 OR child_org_id = :child_org_id)", map[string]interface{}{
			"id":           args.ID,
			"org_id":       args.OrgID,
			"bind_id":      args.BindID,
			"child_org_id": args.ChildOrgID,
			"name":         args.Name,
			"sn":           args.SN,
			"expire_at":    expireAt,
			"file_ids":     args.FileIDs,
			"params":       args.Params,
		})
		if err != nil {
			//fmt.Println("update cert id: ", args.ID, ", expire at: ", expireAt, ", org id and bind id mix, org id: ", args.OrgID, ", bind id: ", args.BindID, ", err: ", err)
		}
	} else {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_cert SET update_at = NOW(), name = :name, sn = :sn, expire_at = :expire_at, file_ids = :file_ids, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id) AND (:bind_id < 1 OR bind_id = :bind_id) AND (:child_org_id < 1 OR child_org_id = :child_org_id)", map[string]interface{}{
			"id":           args.ID,
			"org_id":       args.OrgID,
			"bind_id":      args.BindID,
			"child_org_id": args.ChildOrgID,
			"name":         args.Name,
			"sn":           args.SN,
			"expire_at":    expireAt,
			"file_ids":     args.FileIDs,
			"params":       args.Params,
		})
		if err != nil {
			//fmt.Println("update cert id: ", args.ID, ", expire at: ", expireAt, ", err: ", err)
		}
	}
	return
}

// ArgsUpdateCertExpire 修改证件过期时间参数
type ArgsUpdateCertExpire struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//子商户
	// 实际发起的商户ID，可以和顶级商户orgID一致
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID" check:"id" empty:"true"`
	//过期时间
	ExpireAt string `db:"expire_at" json:"expireAt" check:"isoTime" empty:"true"`
}

// UpdateCertExpire 修改证件过期时间
func UpdateCertExpire(args *ArgsUpdateCertExpire) (err error) {
	var expireAt time.Time
	expireAt, err = CoreFilter.GetTimeByISO(args.ExpireAt)
	if err != nil {
		return
	}
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_cert SET update_at = NOW(), expire_at = :expire_at WHERE id = :id AND (:org_id < 1 OR org_id = :org_id) AND (:child_org_id < 1 OR child_org_id = :child_org_id)", map[string]interface{}{
		"id":           args.ID,
		"org_id":       args.OrgID,
		"child_org_id": args.ChildOrgID,
		"expire_at":    expireAt,
	})
	return
}

// ArgsPayCert 支付费用参数
type ArgsPayCert struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//绑定ID
	// 根据配置决定，可能是用户ID\组织ID\或其他任意主体
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//支付方式
	// system: cash 现金 ; deposit 存储模块 ; weixin 微信支付 ; alipay 支付宝
	// mark: 子渠道信息，例如 weixin 的wxx/merchant
	PaymentChannel CoreSQLFrom.FieldsFrom `db:"payment_channel" json:"paymentChannel"`
	//支付备注
	// 用户环节可根据实际业务需求开放此项
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600" empty:"true"`
}

// PayCert 支付费用
func PayCert(args *ArgsPayCert) (payData FinancePay.FieldsPayType, errCode string, err error) {
	//获取请求
	var certData FieldsCert
	certData, err = GetCert(&ArgsGetCert{
		ID:    args.ID,
		OrgID: -1,
	})
	if err != nil {
		errCode = "cert_not_exist"
		return
	}
	if certData.DeleteAt.Unix() < 100000 {
		errCode = "cert_not_exist"
		err = errors.New("cert is delete")
		return
	}
	//构建支付请求
	payData, errCode, err = FinancePayCreate.CreateUserToOrg(&FinancePayCreate.ArgsCreateUserToOrg{
		UserID:         args.UserID,
		OrgID:          args.OrgID,
		IsRefund:       false,
		Currency:       certData.Currency,
		Price:          certData.Price,
		PaymentChannel: args.PaymentChannel,
		ExpireAt:       CoreFilter.GetNowTimeCarbon().AddMinutes(30).Time,
		Des:            args.Des,
	})
	if err != nil {
		return
	}
	//计算支付方式
	payFromSystem := fmt.Sprint(payData.PaymentChannel.System)
	if payData.PaymentChannel.Mark != "" {
		payFromSystem = payFromSystem + "_" + payData.PaymentChannel.Mark
	}
	certData.Params = CoreSQLConfig.Set(certData.Params, "paySystem", payFromSystem)
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_cert SET update_at = NOW(), pay_failed = false, pay_id = :pay_id, params = :params WHERE id = :id", map[string]interface{}{
		"id":     certData.ID,
		"pay_id": payData.ID,
		"params": certData.Params,
	})
	if err != nil {
		errCode = "update"
		return
	}
	return
}

// ArgsCheckPay 检查支付状态参数
type ArgsCheckPay struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//绑定ID
	// 根据配置决定，可能是用户ID\组织ID\或其他任意主体
	BindID int64 `db:"bind_id" json:"bindID" check:"id" empty:"true"`
}

// CheckPay 检查支付状态
func CheckPay(args *ArgsCheckPay) (isOk bool, err error) {
	//获取请求
	var certData FieldsCert
	certData, err = GetCert(&ArgsGetCert{
		ID:    args.ID,
		OrgID: args.OrgID,
	})
	if err != nil {
		return
	}
	if certData.DeleteAt.Unix() < 100000 {
		err = errors.New("cert is delete")
		return
	}
	if certData.PayID < 1 {
		err = errors.New("no pay")
		return
	}
	//检查支付状态
	if certData.PayFailed {
		return
	} else {
		if certData.PayAt.Unix() > 100000 {
			isOk = true
			return
		}
	}
	//检查支付ID
	var payStatus []FinancePay.DataCheckFinish
	payStatus, err = FinancePay.CheckFinishByIDs(&FinancePay.ArgsCheckFinishByIDs{
		IDs: []int64{certData.PayID},
	})
	if err != nil {
		return
	}
	for _, v := range payStatus {
		if v.ID == certData.PayID {
			if v.IsFinish {
				_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_cert SET update_at = NOW(), pay_at = NOW() WHERE id = :id", map[string]interface{}{
					"id": certData.ID,
				})
				isOk = true
				return
			} else {
				if v.FailedCode != "" {
					_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_cert SET update_at = NOW(), pay_failed = true WHERE id = :id", map[string]interface{}{
						"id": certData.ID,
					})
					return
				}
			}
		}
	}
	return
}

// ArgsUpdateCertAudit 处理审核参数
type ArgsUpdateCertAudit struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//子商户
	// 实际发起的商户ID，可以和顶级商户orgID一致
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID" check:"id" empty:"true"`
	//审核人
	AuditBindID int64 `db:"audit_bind_id" json:"auditBindID" check:"id" empty:"true"`
	//是否拒绝
	IsBan bool `db:"is_ban" json:"isBan" check:"bool"`
	//审核留言
	AuditDes string `db:"audit_des" json:"auditDes" check:"des" min:"1" max:"300" empty:"true"`
}

// UpdateCertAudit 处理审核
func UpdateCertAudit(args *ArgsUpdateCertAudit) (err error) {
	if args.IsBan {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_cert SET update_at = NOW(), audit_bind_id = :audit_bind_id, audit_ban_at = NOW(), audit_des = :audit_des WHERE id = :id AND (:org_id < 1 OR org_id = :org_id) AND (:child_org_id < 1 OR child_org_id = :child_org_id) AND audit_at < to_timestamp(1000000)", args)
	} else {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_cert SET update_at = NOW(), audit_bind_id = :audit_bind_id, audit_at = NOW(), audit_des = :audit_des WHERE id = :id AND (:org_id < 1 OR org_id = :org_id) AND (:child_org_id < 1 OR child_org_id = :child_org_id) AND audit_ban_at < to_timestamp(1000000)", args)
	}
	return
}

// ArgsDeleteCert 删除证书参数
type ArgsDeleteCert struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//子商户
	// 实际发起的商户ID，可以和顶级商户orgID一致
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID" check:"id" empty:"true"`
}

// DeleteCert 删除证书
func DeleteCert(args *ArgsDeleteCert) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "org_cert", "id = :id AND (:org_id < 1 OR org_id = :org_id) AND (:child_org_id < 1 OR child_org_id = :child_org_id)", args)
	return
}

// ArgsDeleteCerts 批量删除证书参数
type ArgsDeleteCerts struct {
	//ID
	IDs pq.Int64Array `db:"ids" json:"ids" check:"ids"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//子商户
	// 实际发起的商户ID，可以和顶级商户orgID一致
	ChildOrgID int64 `db:"child_org_id" json:"childOrgID" check:"id" empty:"true"`
}

// DeleteCerts 批量删除证书
func DeleteCerts(args *ArgsDeleteCerts) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "org_cert", "id = ANY(:ids) AND (:org_id < 1 OR org_id = :org_id) AND (:child_org_id < 1 OR child_org_id = :child_org_id)", args)
	return
}
