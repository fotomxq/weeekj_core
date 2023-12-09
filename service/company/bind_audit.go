package ServiceCompany

import (
	"errors"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetBindAuditList 获取绑定审核列表参数
type ArgsGetBindAuditList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// -1 跳过
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//绑定用户
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//公司ID
	CompanyID int64 `db:"company_id" json:"companyID" check:"id" empty:"true"`
	//是否审核
	// -1 跳过; 0 尚未审核; 1 审核通过; 2 审核拒绝
	AuditStatus int `json:"auditStatus" check:"intThan0" empty:"true"`
}

// GetBindAuditList 获取绑定审核列表
func GetBindAuditList(args *ArgsGetBindAuditList) (dataList []FieldsBindAudit, dataCount int64, err error) {
	//获取数据
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
	if args.CompanyID > -1 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "company_id = :company_id"
		maps["company_id"] = args.CompanyID
	}
	switch args.AuditStatus {
	case 0:
		where = CoreSQL.GetDeleteSQLField(false, where, "audit_at")
	case 1:
		where = CoreSQL.GetDeleteSQLField(true, where, "audit_at")
	case 2:
		where = CoreSQL.GetDeleteSQLField(false, where, "audit_at")
		where = CoreSQL.GetDeleteSQLField(true, where, "ban_at")
	}
	if where == "" {
		where = "true"
	}
	var rawList []FieldsBindAudit
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		"service_company_bind_audit",
		"id",
		"SELECT id FROM service_company_bind_audit WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	if err != nil {
		return
	}
	if len(rawList) < 1 {
		err = errors.New("no data")
		return
	}
	//遍历重组数据
	for _, v := range rawList {
		vData := getBindAudit(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	//反馈
	return
}

// ArgsCreateBindAudit 创建新审核请求参数
type ArgsCreateBindAudit struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//绑定用户
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//公司ID
	CompanyID int64 `db:"company_id" json:"companyID" check:"id"`
	//绑定原因
	Des string `db:"des" json:"des" check:"des" min:"1" max:"600"`
}

// CreateBindAudit 创建新审核请求
func CreateBindAudit(args *ArgsCreateBindAudit) (errCode string, err error) {
	//检查绑定限制
	if !checkBindLimit(args.OrgID, args.UserID, "", "") {
		errCode = "err_limit"
		err = errors.New("bind limit")
		return
	}
	//检查是否具备绑定关系，如果存在则静止提交审核
	var bindData FieldsBind
	err = Router2SystemConfig.MainDB.Get(&bindData, "SELECT id FROM service_company_bind WHERE org_id = $1 AND user_id = $2 AND company_id = $3", args.OrgID, args.UserID, args.CompanyID)
	if bindData.ID > 0 {
		errCode = "err_audit_no_need"
		err = errors.New("have bind")
		return
	}
	//检查是否存在待审核请求
	var id int64
	err = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM service_company_bind_audit WHERE user_id = $1 AND company_id = $2 AND audit_at < to_timestamp(1000000) AND ban_at < to_timestamp(1000000)", args.UserID, args.CompanyID)
	if err == nil && id > 0 {
		errCode = "err_audit_replace"
		err = errors.New("have bind audit wait")
		return
	}
	//创建请求
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO service_company_bind_audit(org_id, user_id, company_id, des, audit_at, ban_at, audit_des) VALUES(:org_id, :user_id, :company_id, :des, to_timestamp(0), to_timestamp(0), '')", args)
	if err != nil {
		errCode = "err_insert"
		return
	}
	//反馈
	return
}

// ArgsAuditBind 审核绑定关系请求
type ArgsAuditBind struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否审核
	// 1 审核通过; 2 审核拒绝
	AuditStatus int `json:"auditStatus" check:"intThan0"`
	//通过或拒绝原因
	AuditDes string `db:"audit_at" json:"auditDes"`
	//赋予能力
	// 系统约定的几个特定能力，平台无法编辑该范围，只能授权
	Managers pq.StringArray `db:"managers" json:"managers"`
}

// AuditBind 审核绑定关系
func AuditBind(args *ArgsAuditBind) (err error) {
	auditData := getBindAudit(args.ID)
	if auditData.ID < 1 || !CoreFilter.EqID2(args.OrgID, auditData.OrgID) || CoreSQL.CheckTimeHaveData(auditData.AuditAt) || CoreSQL.CheckTimeHaveData(auditData.BanAt) {
		err = errors.New("no data")
		return
	}
	if args.AuditStatus == 1 {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_company_bind_audit SET audit_at = NOW(), ban_at = to_timestamp(0), audit_des = :audit_des WHERE id = :id", map[string]interface{}{
			"id":        auditData.ID,
			"audit_des": args.AuditDes,
		})
		if err != nil {
			return
		}
		_, err = SetBind(&ArgsSetBind{
			OrgID:      auditData.OrgID,
			UserID:     auditData.UserID,
			NationCode: "",
			Phone:      "",
			CompanyID:  auditData.CompanyID,
			Managers:   args.Managers,
		})
		if err != nil {
			return
		}
	} else {
		_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE service_company_bind_audit SET audit_at = to_timestamp(0), ban_at = NOW(), audit_des = :audit_des WHERE id = :id", map[string]interface{}{
			"id":        auditData.ID,
			"audit_des": args.AuditDes,
		})
		if err != nil {
			return
		}
	}
	deleteBindAuditCache(auditData.ID)
	return
}

func getBindAudit(id int64) (data FieldsBindAudit) {
	cacheMark := getBindAuditCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, user_id, company_id, des, audit_at, ban_at, audit_des FROM service_company_bind_audit WHERE id = $1", id)
	if err != nil {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, 1800)
	return
}
