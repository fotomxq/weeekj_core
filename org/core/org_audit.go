package OrgCoreCore

import (
	"errors"
	"fmt"
	BaseConfig "gitee.com/weeekj/weeekj_core/v5/base/config"
	CoreFilter "gitee.com/weeekj/weeekj_core/v5/core/filter"
	CoreLog "gitee.com/weeekj/weeekj_core/v5/core/log"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	OrgMap "gitee.com/weeekj/weeekj_core/v5/org/map"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	UserCore "gitee.com/weeekj/weeekj_core/v5/user/core"
	"github.com/lib/pq"
	"github.com/mozillazg/go-pinyin"
	"strings"
)

// ArgsGetOrgAuditList 获取审核列表参数
type ArgsGetOrgAuditList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//用户ID
	UserID int64 `json:"userID" check:"id" empty:"true"`
	//上级组织ID
	ParentID int64 `json:"parentID" check:"id" empty:"true"`
	//上级控制权限限制
	ParentFunc pq.StringArray `db:"parent_func" json:"parentFunc" check:"marks" empty:"true"`
	//开通业务
	// 该内容只有总管理员或订阅能进行控制
	OpenFunc pq.StringArray `db:"open_func" json:"openFunc" check:"marks" empty:"true"`
	//是否审核
	NeedIsAudit bool `db:"need_is_audit" json:"needIsAudit" check:"bool" empty:"true"`
	IsAudit     bool `db:"is_audit" json:"isAudit" check:"bool" empty:"true"`
	//是否拒绝
	IsBan bool `db:"is_ban" json:"isBan" check:"bool" empty:"true"`
	//审核人员
	AuditUserID int64 `db:"audit_user_id" json:"auditUserID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetOrgAuditList 获取审核列表
func GetOrgAuditList(args *ArgsGetOrgAuditList) (dataList []FieldsOrgAudit, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.ParentID > -1 {
		where = where + " AND :parent_id = :parent_id"
		maps["parent_id"] = args.ParentID
	}
	if args.UserID > 0 {
		where = where + " AND user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if len(args.ParentFunc) > 0 {
		where = where + " AND parent_func @> :parent_func"
		maps["parent_func"] = args.ParentFunc
	}
	if len(args.OpenFunc) > 0 {
		where = where + " AND open_func @> :open_func"
		maps["open_func"] = args.OpenFunc
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
	if args.AuditUserID > 0 {
		where = where + " AND audit_user_id = :audit_user_id"
		maps["audit_user_id"] = args.AuditUserID
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%' OR ban_des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "org_core_audit"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, delete_at, audit_at, ban_at, ban_des, audit_user_id, user_id, key, name, des, parent_id, parent_func, open_func, params "+"FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "audit_at", "ban_at"},
	)
	return
}

// ArgsCreateOrgAudit 创建新审核请求参数
type ArgsCreateOrgAudit struct {
	//所属用户
	// 掌管该数据的用户，创建人和根管理员，不可删除只能更换
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//企业唯一标识码
	// 用于特殊识别和登陆识别等操作
	Key string `db:"key" json:"key" check:"mark" empty:"true"`
	//构架名称，或组织名称
	Name string `db:"name" json:"name" check:"name"`
	//组织描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000"`
	//上级ID
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
	//上级控制权限限制
	ParentFunc pq.StringArray `db:"parent_func" json:"parentFunc" check:"marks" empty:"true"`
	//开通业务
	// 该内容只有总管理员或订阅能进行控制
	OpenFunc pq.StringArray `db:"open_func" json:"openFunc" check:"marks" empty:"true"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateOrgAudit 创建新审核请求
func CreateOrgAudit(args *ArgsCreateOrgAudit) (auditData FieldsOrgAudit, errCode string, err error) {
	//修正参数
	if args.ParentID < 1 {
		args.ParentID = 0
	}
	//用户禁止重复申请
	err = Router2SystemConfig.MainDB.Get(&auditData, "SELECT id FROM org_core_audit WHERE delete_at < to_timestamp(1000000) AND audit_at < to_timestamp(1000000) AND ban_at < to_timestamp(1000000) AND user_id = $1 LIMIT 1", args.UserID)
	if err == nil && auditData.ID > 0 {
		errCode = "err_org_replace_audit"
		err = errors.New("user replace audit")
		return
	}
	//生成key
	if args.Key == "" {
		args.Key = makeKey(args.Name)
	}
	//获取用户信息
	if args.UserID < 1 {
		errCode = "err_user"
		err = errors.New("user id is empty")
		return
	}
	_, err = UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
		ID:    args.UserID,
		OrgID: -1,
	})
	if err != nil {
		errCode = "err_user"
		err = errors.New("get user by id, " + err.Error())
		return
	}
	//检查上一级
	if args.ParentID > 0 {
		//获取上一级
		var parentData FieldsOrg
		parentData, err = GetOrg(&ArgsGetOrg{
			ID: args.ParentID,
		})
		if err != nil {
			errCode = "err_org_parent_not_exist"
			return
		}
		//检查上一级功能和本次开通功能
		err = checkOrgInParentFunc(parentData.OpenFunc, args.OpenFunc)
		if err != nil {
			errCode = "err_org_func_not_in_area"
			return
		}
	}
	//检查key
	if args.Key == "" {
		args.Key, err = CoreFilter.GetRandStr3(10)
		if err != nil {
			errCode = "err_key"
			return
		}
	}
	var orgData FieldsOrg
	orgData, err = GetOrgByKey(&ArgsGetOrgByKey{
		Key: args.Key,
	})
	if err == nil && orgData.ID > 0 {
		errCode = "err_org_replace"
		err = errors.New("key is exist")
		return
	}
	//生成数据
	if err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "org_core_audit", "INSERT INTO org_core_audit (audit_at, ban_at, ban_des, audit_user_id, user_id, key, name, des, parent_id, parent_func, open_func, params) VALUES (to_timestamp(0), to_timestamp(0), '', 0, :user_id, :key, :name, :des, :parent_id, :parent_func, :open_func, :params)", args, &auditData); err != nil {
		errCode = "err_insert"
		return
	}
	return
}

// ArgsUpdateOrgAuditPass 通过审核参数
type ArgsUpdateOrgAuditPass struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//上级组织ID
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
	//操作人
	AuditUserID int64 `db:"audit_user_id" json:"auditUserID" check:"id" empty:"true"`
}

// AuditOrgCreateMap 新增特殊逻辑判断平台是否开始同步创建组织地图开关
func AuditOrgCreateMap(auditData FieldsOrgAudit, orgID int64) {
	// 通过扩展参数获取创建地图相关数据
	orgMapProvince := auditData.Params.GetValIntNoBool("OrgMapProvince")
	orgMapCity := auditData.Params.GetValIntNoBool("OrgMapCity")
	newIndexShowGps := auditData.Params.GetValNoErr("IndexShowGPS")
	newGPSArr := strings.Split(newIndexShowGps, ",")
	var longitude, latitude float64
	if newIndexShowGps != "" && len(newGPSArr) >= 3 {
		longitude = CoreFilter.GetFloat64ByStringNoErr(newGPSArr[1])
		latitude = CoreFilter.GetFloat64ByStringNoErr(newGPSArr[2])
	}

	_, err := OrgMap.CreateMap(&OrgMap.ArgsCreateMap{
		OrgID:         orgID,
		UserID:        auditData.UserID,
		ParentID:      0,
		CoverFileID:   auditData.Params.GetValInt64NoBool("CoverFileID"),
		CoverFileIDs:  []int64{},
		Name:          auditData.Name,
		Des:           auditData.Des,
		Country:       86,
		Province:      orgMapProvince,
		City:          orgMapCity,
		Address:       auditData.Params.GetValNoErr("IndexShowAddress"),
		MapType:       1,
		Longitude:     longitude,
		Latitude:      latitude,
		AdCountLimit:  -1,
		ViewTimeLimit: 0,
		Params:        CoreSQLConfig.FieldsConfigsType{},
	})

	if err != nil {
		CoreLog.Error(err)
	}
}

// UpdateOrgAuditPass 通过审核
func UpdateOrgAuditPass(args *ArgsUpdateOrgAuditPass) (errCode string, err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_core_audit SET audit_at = NOW(), ban_at = to_timestamp(0), audit_user_id = :audit_user_id WHERE id = :id AND (:parent_id < 1 OR parent_id = :parent_id)", args)
	if err == nil {
		var auditData FieldsOrgAudit
		err = Router2SystemConfig.MainDB.Get(&auditData, "SELECT id, create_at, delete_at, audit_at, ban_at, ban_des, audit_user_id, user_id, key, name, des, parent_id, parent_func, open_func, params FROM org_core_audit WHERE id = $1 LIMIT 1", args.ID)
		if err != nil {
			errCode = "audit_empty"
			return
		}
		var newOrgData FieldsOrg
		newOrgData, errCode, err = CreateOrg(&ArgsCreateOrg{
			UserID:     auditData.UserID,
			Key:        auditData.Key,
			Name:       auditData.Name,
			Des:        auditData.Des,
			ParentID:   auditData.ParentID,
			ParentFunc: auditData.ParentFunc,
			OpenFunc:   auditData.OpenFunc,
			SortID:     0,
		})
		if err != nil {
			return
		}
		// 判断审核组织后是否同步创建组织地图
		orgAuditCreateMap := BaseConfig.GetDataBoolNoErr("OrgAuditCreateMap")
		// 如果开启了同步创建组织地图，则创建组织地图
		if orgAuditCreateMap {
			AuditOrgCreateMap(auditData, newOrgData.ID)
		}
		for _, v := range auditData.Params {
			if strings.Contains(v.Mark, "*") {
				continue
			}
			isFind := true
			switch v.Mark {
			case "CoverFileID":
			case "OrgPhone":
			case "OrgDes":
			case "OrgBusinessHours":
			case "IndexShowAddress":
			case "IndexShowGPS":
			case "OrgAuditLicenseFileID":
			case "OrgAuditIDCardAFileID":
			case "OrgAuditIDCardBFileID":
			case "OrgAuditLegalName":
			default:
				isFind = false
			}
			if isFind {
				_ = Config.SetConfigValSimple(newOrgData.ID, v.Mark, v.Val)
			}
		}
		auditData.Params = CoreSQLConfig.Set(auditData.Params, "NewOrgID", newOrgData.ID)
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_core_audit SET params = :params WHERE id = :id", map[string]interface{}{
			"id":     args.ID,
			"params": auditData.Params,
		})
		if err != nil {
			errCode = "err_update"
			return
		}
		return
	} else {
		errCode = "no_audit"
	}
	return
}

// ArgsUpdateOrgAuditBan 拒绝审核参数
type ArgsUpdateOrgAuditBan struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//上级组织ID
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
	//操作人
	AuditUserID int64 `db:"audit_user_id" json:"auditUserID" check:"id" empty:"true"`
	//拒绝审核原因
	BanDes string `db:"ban_des" json:"banDes" check:"des" min:"1" max:"600" empty:"true"`
}

// UpdateOrgAuditBan 拒绝审核
func UpdateOrgAuditBan(args *ArgsUpdateOrgAuditBan) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_core_audit SET audit_at = to_timestamp(0), ban_at = NOW(), audit_user_id = :audit_user_id, ban_des = :ban_des WHERE id = :id AND (:parent_id < 1 OR parent_id = :parent_id)", args)
	return
}

// ArgsDeleteOrgAudit 删除审核参数
type ArgsDeleteOrgAudit struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//上级组织ID
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
}

// DeleteOrgAudit 删除审核
func DeleteOrgAudit(args *ArgsDeleteOrgAudit) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "org_core_audit", "id = :id AND (:parent_id < 1 OR parent_id = :parent_id)", args)
	return
}

// 生成新的key
func makeKey(orgName string) string {
	//生成拼音
	newKeys := pinyin.LazyPinyin(orgName, pinyin.NewArgs())
	var newKey string
	for k := 0; k < len(newKeys); k++ {
		newKey = newKey + newKeys[k]
	}
	//检查key是否存在
	if newKey != "" {
		var id int64
		err := Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM org_core WHERE key = $1", newKey)
		if err != nil || id < 1 {
			err = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM org_core_audit WHERE key = $1", newKey)
			if err == nil && id > 0 {
				return newKey
			}
		}
	}
	step := 1
	for {
		if step > 100 {
			break
		}
		step += 1
		var id int64
		newKey2 := newKey + fmt.Sprint(step)
		err := Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM org_core WHERE key = $1", newKey2)
		if err == nil && id > 0 {
			continue
		}
		err = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM org_core_audit WHERE key = $1", newKey2)
		if err == nil && id > 0 {
			continue
		}
		break
	}
	//反馈数据
	return newKey
}
