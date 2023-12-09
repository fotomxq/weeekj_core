package OrgCoreCore

import (
	"errors"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	UserCore "gitee.com/weeekj/weeekj_core/v5/user/core"
	"github.com/lib/pq"
)

// ArgsGetBindAuditList 获取审核列表参数
type ArgsGetBindAuditList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// 必填
	OrgID int64 `json:"orgID" check:"id"`
	//用户ID
	UserID int64 `json:"userID" check:"id" empty:"true"`
	//分组ID
	GroupID int64 `json:"groupID" check:"id" empty:"true"`
	//符合权限的
	Manager string `json:"manager" check:"mark" empty:"true"`
	//是否审核
	NeedIsAudit bool `db:"need_is_audit" json:"needIsAudit" check:"bool" empty:"true"`
	IsAudit     bool `db:"is_audit" json:"isAudit" check:"bool" empty:"true"`
	//是否拒绝
	IsBan bool `db:"is_ban" json:"isBan" check:"bool" empty:"true"`
	//审核人员
	AuditBindID int64 `db:"audit_bind_id" json:"auditBindID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetBindAuditList 获取审核列表
func GetBindAuditList(args *ArgsGetBindAuditList) (dataList []FieldsBindAudit, dataCount int64, err error) {
	where := "org_id = :org_id"
	maps := map[string]interface{}{
		"org_id": args.OrgID,
	}
	if args.IsRemove {
		where = where + " AND delete_at > to_timestamp(1000000)"
	} else {
		where = where + " AND delete_at < to_timestamp(1000000)"
	}
	if args.UserID > 0 {
		where = where + " AND user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.GroupID > 0 {
		where = where + " AND :group_id = ANY(group_ids)"
		maps["group_id"] = args.GroupID
	}
	if args.Manager != "" {
		where = where + " AND :manager = ANY(manager)"
		maps["manager"] = args.Manager
	}
	if args.NeedIsAudit {
		if args.IsAudit {
			where = where + " AND audit_at > to_timestamp(1000000) AND ban_at <= to_timestamp(1000000)"
		} else {
			if args.IsBan {
				where = where + " AND ban_at > to_timestamp(1000000)"
			} else {
				where = where + " AND audit_at <= to_timestamp(1000000) AND ban_at <= to_timestamp(1000000)"
			}
		}
	}
	if args.AuditBindID > 0 {
		where = where + " AND audit_bind_id = :audit_bind_id"
		maps["audit_bind_id"] = args.AuditBindID
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%' OR ban_des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "org_core_bind_audit"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, delete_at, audit_at, ban_at, ban_des, audit_bind_id, user_id, name, org_id, group_ids, manager, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "audit_at", "ban_at"},
	)
	return
}

// ArgsCreateBindAudit 创建新审核请求
type ArgsCreateBindAudit struct {
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//组织分组ID列
	GroupIDs pq.Int64Array `db:"group_ids" json:"groupIDs" check:"ids" empty:"true"`
	//权利主张
	Manager pq.StringArray `db:"manager" json:"manager" check:"marks" empty:"true"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateBindAudit 创建新审核请求
func CreateBindAudit(args *ArgsCreateBindAudit) (auditData FieldsBindAudit, errCode string, err error) {
	//用户禁止重复申请
	err = Router2SystemConfig.MainDB.Get(&auditData, "SELECT id FROM org_core_bind_audit WHERE delete_at < to_timestamp(1000000) AND audit_at < to_timestamp(1000000) AND ban_at < to_timestamp(1000000) AND user_id = $1 LIMIT 1", args.UserID)
	if err == nil && auditData.ID > 0 {
		errCode = "replace_audit"
		err = errors.New("user replace audit")
		return
	}
	//检查用户是否存在绑定关系
	_, err = GetBindByUserAndOrg(&ArgsGetBindByUserAndOrg{
		UserID: args.UserID,
		OrgID:  args.OrgID,
	})
	if err == nil {
		errCode = "have_bind"
		err = errors.New("user have bind")
		return
	}
	//获取用户信息
	if args.UserID < 1 {
		errCode = "user_not_exist"
		err = errors.New("user id is empty")
		return
	}
	_, err = UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
		ID:    args.UserID,
		OrgID: -1,
	})
	if err != nil {
		errCode = "user_not_exist"
		err = errors.New("get user by id, " + err.Error())
		return
	}
	//检查组织
	var orgData FieldsOrg
	orgData, err = GetOrg(&ArgsGetOrg{
		ID: args.OrgID,
	})
	if err != nil || orgData.ID < 1 {
		errCode = "org_not_exist"
		err = errors.New("org not exist")
		return
	}
	//生成数据
	if err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "org_core_bind_audit", "INSERT INTO org_core_bind_audit (audit_at, ban_at, ban_des, audit_bind_id, user_id, name, org_id, group_ids, manager, params) VALUES (to_timestamp(0), to_timestamp(0), '', 0, :user_id, :name, :org_id, :group_ids, :manager, :params)", args, &auditData); err != nil {
		errCode = "insert"
		return
	}
	return
}

// ArgsUpdateBindAuditPass 通过审核参数
type ArgsUpdateBindAuditPass struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//操作人
	AuditBindID int64 `db:"audit_bind_id" json:"auditBindID" check:"id" empty:"true"`
}

// UpdateBindAuditPass 通过审核
func UpdateBindAuditPass(args *ArgsUpdateBindAuditPass) (err error) {
	var auditData FieldsBindAudit
	err = Router2SystemConfig.MainDB.Get(&auditData, "SELECT id FROM org_core_bind_audit WHERE id = $1 AND delete_at < to_timestamp(1000000)", args.ID)
	if err != nil {
		return
	}
	userData, _ := UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
		ID:    auditData.UserID,
		OrgID: auditData.OrgID,
	})
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_core_bind_audit SET audit_at = NOW(), ban_at = to_timestamp(0), audit_bind_id = :audit_bind_id WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err == nil {
		err = Router2SystemConfig.MainDB.Get(&auditData, "SELECT id, create_at, delete_at, audit_at, ban_at, ban_des, audit_bind_id, user_id, name, org_id, group_ids, manager, params FROM org_core_bind_audit WHERE id = $1 LIMIT 1", args.ID)
		if err != nil {
			return
		}
		_, err = SetBind(&ArgsSetBind{
			UserID:     auditData.UserID,
			Avatar:     0,
			Name:       auditData.Name,
			OrgID:      auditData.OrgID,
			GroupIDs:   auditData.GroupIDs,
			Manager:    auditData.Manager,
			NationCode: userData.NationCode,
			Phone:      userData.Phone,
			Email:      userData.Email,
			SyncSystem: "",
			SyncID:     0,
			SyncHash:   "",
			Params:     auditData.Params,
		})
		return
	}
	return
}

// ArgsUpdateBindAuditBan 拒绝审核参数
type ArgsUpdateBindAuditBan struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//操作人
	AuditBindID int64 `db:"audit_bind_id" json:"auditBindID" check:"id" empty:"true"`
	//拒绝审核原因
	BanDes string `db:"ban_des" json:"banDes" check:"des" min:"1" max:"600" empty:"true"`
}

// UpdateBindAuditBan 拒绝审核
func UpdateBindAuditBan(args *ArgsUpdateBindAuditBan) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_core_bind_audit SET audit_at = to_timestamp(0), ban_at = NOW(), audit_bind_id = :audit_bind_id, ban_des = :ban_des WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// ArgsDeleteBindAudit 删除审核参数
type ArgsDeleteBindAudit struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// DeleteBindAudit 删除审核
func DeleteBindAudit(args *ArgsDeleteBindAudit) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "org_core_bind_audit", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}
