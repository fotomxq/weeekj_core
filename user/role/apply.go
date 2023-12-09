package UserRole

import (
	"errors"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetApplyList 获取申请列表参数
type ArgsGetApplyList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//角色类型
	RoleType int64 `db:"role_type" json:"roleType" check:"id" empty:"true"`
	//用户ID
	// 允许为0，则该信息不属于任何用户，或不和任何用户关联
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
	//是否审核
	NeedIsAudit bool `json:"needIsAudit" check:"bool"`
	IsAudit     bool `json:"isAudit" check:"bool"`
	//是否被删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetApplyList 获取申请列表
func GetApplyList(args *ArgsGetApplyList) (dataList []FieldsApply, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.RoleType > -1 {
		where = where + " AND role_type = :role_type"
		maps["role_type"] = args.RoleType
	}
	if args.UserID > -1 {
		where = where + " AND user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if args.NeedIsAudit {
		if args.IsAudit {
			where = where + " AND audit_at > to_timestamp(1000000)"
		} else {
			where = where + " AND audit_at <= to_timestamp(1000000)"
		}
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "user_role_apply"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, audit_at, audit_des, audit_ban_des, role_type, user_id, name, country, city, gender, phone, cover_file_id, cert_files, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "audit_at"},
	)
	return
}

// ArgsGetApplyID 查看申请详情ID参数
type ArgsGetApplyID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// GetApplyID 查看申请详情ID
func GetApplyID(args *ArgsGetApplyID) (data FieldsApply, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, audit_at, audit_des, audit_ban_des, role_type, user_id, name, country, city, gender, phone, cover_file_id, cert_files, params FROM user_role_apply WHERE id = $1 AND ($2 < 1 OR user_id = $2)", args.ID, args.UserID)
	if err == nil && data.ID < 1 {
		err = errors.New("data not exist")
		return
	}
	return
}

// ArgsCreateApply 创建新的申请参数
type ArgsCreateApply struct {
	//申请描述
	AuditDes string `db:"audit_des" json:"auditDes" check:"des" min:"1" max:"600" empty:"true"`
	//角色类型
	RoleType int64 `db:"role_type" json:"roleType" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//姓名
	Name string `db:"name" json:"name" check:"name"`
	//所属国家 国家代码
	// eg: china => 86
	Country int `db:"country" json:"country" check:"country"`
	//城市编码
	City string `db:"city" json:"city" check:"cityCode"`
	//性别
	// 0 男 1 女 2 未知
	Gender int `db:"gender" json:"gender" check:"gender"`
	//联系电话
	Phone string `db:"phone" json:"phone" check:"phone"`
	//个人照片
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//证件列
	CertFiles pq.Int64Array `db:"cert_files" json:"certFiles" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params" check:"params" empty:"true"`
}

// CreateApply 创建新的申请
func CreateApply(args *ArgsCreateApply) (data FieldsApply, err error) {
	//找到对应的角色类型
	var id int64
	err = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM user_role_type WHERE id = $1", args.RoleType)
	if err != nil || id < 1 {
		err = errors.New("role type not exist")
		return
	}
	//检查是否存在正在申请处理
	var applyID int64
	err = Router2SystemConfig.MainDB.Get(&applyID, "SELECT id FROM user_role_apply WHERE role_type = $1 AND user_id = $2 AND delete_at < to_timestamp(1000000) AND audit_at < to_timestamp(1000000)", args.RoleType, args.UserID)
	if err == nil && applyID > 0 {
		err = errors.New("have apply")
		return
	}
	//写入数据
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "user_role_apply", "INSERT INTO user_role_apply (audit_des, audit_ban_des, role_type, user_id, name, country, city, gender, phone, cover_file_id, cert_files, params) VALUES (:audit_des,'',:role_type,:user_id,:name,:country,:city,:gender,:phone,:cover_file_id,:cert_files,:params)", args, &data)
	return
}

// ArgsAuditApply 审核申请参数
type ArgsAuditApply struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//是否审核通过
	IsAudit bool `db:"is_audit" json:"isAudit" check:"bool"`
	//拒绝原因
	AuditBanDes string `db:"audit_ban_des" json:"auditBanDes" check:"des" min:"1" max:"600" empty:"true"`
}

// AuditApply 审核申请
func AuditApply(args *ArgsAuditApply) (errCode string, err error) {
	if args.IsAudit {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE user_role_apply SET update_at = NOW(), audit_at = NOW() WHERE id = :id AND audit_at < to_timestamp(1000000)", map[string]interface{}{
			"id": args.ID,
		})
		if err != nil {
			errCode = "update_apply"
			return
		}
		var data FieldsApply
		data, err = GetApplyID(&ArgsGetApplyID{
			ID: args.ID,
		})
		if err != nil {
			errCode = "apply_not_exist"
			return
		}
		_, err = SetRole(&ArgsSetRole{
			RoleType:    data.RoleType,
			ApplyID:     data.ID,
			UserID:      data.UserID,
			Name:        data.Name,
			Country:     data.Country,
			City:        data.City,
			Gender:      data.Gender,
			Phone:       data.Phone,
			CoverFileID: data.CoverFileID,
			CertFiles:   data.CertFiles,
			Params:      data.Params,
		})
		if err != nil {
			errCode = "set_role"
			return
		}
	} else {
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE user_role_apply SET update_at = NOW(), audit_at = NOW(), audit_ban_des = :audit_ban_des WHERE id = :id AND audit_at < to_timestamp(1000000)", map[string]interface{}{
			"id":            args.ID,
			"audit_ban_des": args.AuditBanDes,
		})
		if err != nil {
			errCode = "update_apply"
			return
		}
		var data FieldsApply
		data, err = GetApplyID(&ArgsGetApplyID{
			ID: args.ID,
		})
		if err != nil {
			errCode = "apply_not_exist"
			return
		}
		_, _ = CoreSQL.DeleteOneSoft(Router2SystemConfig.MainDB.DB, "user_role", "apply_id", map[string]interface{}{
			"apply_id": data.ID,
		})
	}
	return
}

// ArgsDeleteApply 删除申请参数
type ArgsDeleteApply struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id" empty:"true"`
}

// DeleteApply 删除申请
func DeleteApply(args *ArgsDeleteApply) (err error) {
	_, err = CoreSQL.DeleteOneSoft(Router2SystemConfig.MainDB.DB, "user_role_apply", "id", args)
	return
}
