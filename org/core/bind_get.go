package OrgCoreCore

import (
	"errors"
	"fmt"
	BaseConfig "github.com/fotomxq/weeekj_core/v5/base/config"
	ClassConfig "github.com/fotomxq/weeekj_core/v5/class/config"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLIDs "github.com/fotomxq/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	"github.com/lib/pq"
	"time"
)

// ArgsGetBindList 获取某个组织的所有绑定关系参数
type ArgsGetBindList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	// 必填
	OrgID int64 `json:"orgID" check:"id"`
	//用户ID
	UserID int64 `json:"userID" check:"id" empty:"true"`
	//分组ID
	GroupID int64 `json:"groupID" check:"id" empty:"true"`
	//角色配置列
	RoleConfigID int64 `json:"roleConfigID" check:"id" empty:"true"`
	//符合权限的
	Manager string `json:"manager" check:"mark" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetBindList 获取某个组织的所有绑定关系
func GetBindList(args *ArgsGetBindList) (dataList []FieldsBind, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > 0 {
		where = where + " AND :org_id = org_id"
		maps["org_id"] = args.OrgID
	}
	if args.UserID > 0 {
		where = where + " AND :user_id = user_id"
		maps["user_id"] = args.UserID
	}
	if args.GroupID > 0 {
		where = where + " AND :group_id = ANY(group_ids)"
		maps["group_id"] = args.GroupID
	}
	if args.RoleConfigID > 0 {
		where = where + " AND :role_config_id = ANY(role_config_ids)"
		maps["role_config_id"] = args.RoleConfigID
	}
	if args.Manager != "" {
		where = where + " AND :manager = ANY(manager)"
		maps["manager"] = args.Manager
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%' OR phone ILIKE '%' || :search || '%' OR email ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	var rawList []FieldsBind
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		"org_core_bind",
		"id",
		"SELECT id FROM org_core_bind WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "last_at", "name"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getBindByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsSearchBind 搜索组织成员参数
type ArgsSearchBind struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

type DataSearchBind struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//名称
	Name string `db:"name" json:"name"`
}

// SearchBind 搜索组织成员
func SearchBind(args *ArgsSearchBind) (dataList []DataSearchBind, err error) {
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, name FROM org_core_bind WHERE org_id = $1 AND ($2 != '' OR name ILIKE '%' || $2 || '%') ORDER BY id LIMIT 10", args.OrgID, args.Search)
	if err != nil {
		return
	}
	if len(dataList) < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsGetBind 查看绑定关系参数
type ArgsGetBind struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//用户ID
	UserID int64 `json:"userID" check:"id" empty:"true"`
}

// GetBind 查看绑定关系
func GetBind(args *ArgsGetBind) (data FieldsBind, err error) {
	data = getBindByID(args.ID)
	if data.ID < 1 || !CoreFilter.EqID2(args.OrgID, data.OrgID) || !CoreFilter.EqID2(args.UserID, data.UserID) {
		err = errors.New("no data")
		return
	}
	return
}

func GetBindNoErr(id int64, orgID int64, userID int64) (data FieldsBind) {
	data = getBindByID(id)
	if data.ID < 1 || !CoreFilter.EqID2(orgID, data.OrgID) || !CoreFilter.EqID2(userID, data.UserID) {
		data = FieldsBind{}
		return
	}
	return
}

// ArgsGetBindByUserAndOrg 获取绑定来源在指定组织中的关系参数
type ArgsGetBindByUserAndOrg struct {
	//用户ID
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// GetBindByUserAndOrg 获取绑定来源在指定组织中的关系
func GetBindByUserAndOrg(args *ArgsGetBindByUserAndOrg) (data FieldsBind, err error) {
	err = CoreSQL.GetOne(Router2SystemConfig.MainDB.DB, &data, "SELECT id FROM org_core_bind WHERE org_id = :org_id AND user_id = :user_id AND delete_at < to_timestamp(1000000)", map[string]interface{}{
		"user_id": args.UserID,
		"org_id":  args.OrgID,
	})
	if err != nil {
		err = errors.New(fmt.Sprint(err, ", org id: ", args.OrgID, ", user id: ", args.UserID, ", err: ", err))
		return
	}
	data = getBindByID(data.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsGetBindByUser 查看某个绑定来源的所有绑定关系参数
type ArgsGetBindByUser struct {
	//用户ID
	UserID int64 `json:"userID" check:"id"`
}

// GetBindByUser 查看某个绑定来源的所有绑定关系
func GetBindByUser(args *ArgsGetBindByUser) (dataList []FieldsBind, err error) {
	where := "user_id = :user_id"
	where = CoreSQL.GetDeleteSQL(false, where)
	maps := map[string]interface{}{
		"user_id": args.UserID,
	}
	var rawList []FieldsBind
	err = CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		"SELECT id FROM org_core_bind WHERE "+where,
		maps,
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getBindByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// GetBindByUserIDOnly 通过组织和用户ID查询成员
func GetBindByUserIDOnly(orgID int64, userID int64) (data FieldsBind) {
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_core_bind WHERE org_id = $1 AND user_id = $2 AND delete_at < to_timestamp(1000000) LIMIT 1", orgID, userID)
	if err != nil {
		return
	}
	data = getBindByID(data.ID)
	return
}

// ArgsGetBinds 查询一组ID参数
type ArgsGetBinds struct {
	//一组ID
	IDs pq.Int64Array `db:"ids" json:"ids"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove"`
}

// GetBinds 查询一组ID
func GetBinds(args *ArgsGetBinds) (dataList []FieldsBind, err error) {
	var rawList []FieldsBind
	err = CoreSQLIDs.GetIDsOrgAndDelete(&rawList, "org_core_bind", "id", args.IDs, args.OrgID, args.HaveRemove)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getBindByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// GetBindsName 查询一组ID的名称
// 该方法将绕过删除项
func GetBindsName(args *ArgsGetBinds) (data map[int64]string, err error) {
	data, err = CoreSQLIDs.GetIDsOrgNameAndDelete("org_core_bind", args.IDs, args.OrgID, args.HaveRemove)
	return
}

func GetBindName(bindID int64) string {
	if bindID < 1 {
		return ""
	}
	data := getBindByID(bindID)
	if data.ID < 1 {
		return ""
	}
	return data.Name
}

func GetBindPhone(bindID int64) string {
	if bindID < 1 {
		return ""
	}
	data := getBindByID(bindID)
	if data.ID < 1 {
		return ""
	}
	return data.Phone
}

// GetBindNameAndAvatar 获取姓名和头像
func GetBindNameAndAvatar(bindID int64) (name string, avatar int64) {
	data := getBindByID(bindID)
	if data.ID < 1 {
		return
	}
	name = data.Name
	avatar = data.Avatar
	return
}

// ArgsGetOrgBinds 查询一组ID的参数
type ArgsGetOrgBinds struct {
	//一组ID
	IDs pq.Int64Array `db:"ids" json:"ids" check:"ids"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
}

// GetOrgBindsName 查询一组ID的名称
func GetOrgBindsName(args *ArgsGetOrgBinds) (data map[int64]string, err error) {
	data, err = CoreSQLIDs.GetIDsOrgNameAndDelete("org_core_bind", args.IDs, args.OrgID, args.HaveRemove)
	return
}

// GetOrgBindsNameAndAvatar 获取成员姓名和头像
func GetOrgBindsNameAndAvatar(args *ArgsGetOrgBinds) (data []CoreSQLIDs.DataGetIDsOrgNameAndDeleteAndAvatar, err error) {
	data, err = CoreSQLIDs.GetIDsOrgNameAndDeleteAndAvatar("org_core_bind", args.IDs, args.OrgID, args.HaveRemove)
	return
}

type DataGetBindByUserMarge struct {
	//组织ID
	OrgID int64 `json:"orgID"`
	//组织所有人用户ID
	OrgUserID int64 `json:"orgUserID"`
	//组织key
	OrgKey string `json:"orgKey"`
	//组织名称
	OrgName string `json:"orgName"`
	//组织描述
	OrgDes string `json:"orgDes"`
	//组织封面
	OrgCoverFileID int64 `json:"orgCoverFileID"`
	//上级ID
	OrgParentID int64 `json:"orgParentID"`
	//上级控制权限限制
	OrgParentFunc pq.StringArray `json:"orgParentFunc"`
	//开通业务
	// 该内容只有总管理员或订阅能进行控制
	OrgOpenFunc pq.StringArray `json:"orgOpenFunc"`
	//绑定关系ID
	BindID int64 `json:"bindID"`
	//上次登陆时间
	LastAt time.Time `json:"lastAt"`
	//用户头像
	Avatar int64 `db:"avatar" json:"avatar" check:"id" empty:"true"`
	//绑定关系名称
	Name string `json:"name"`
	//绑定权限
	Manager []string `json:"manager"`
	//参与的组织分组
	Groups []DataGetBindByUserMargeGroup `json:"groups"`
	//扩展信息
	Params CoreSQLConfig.FieldsConfigsType `json:"params"`
}

// DataGetBindByUserMargeGroup 组织内的分组结构
type DataGetBindByUserMargeGroup struct {
	//分组ID
	ID int64 `json:"id"`
	//分组名称
	Name string `json:"name"`
}

func GetBindByUserMarge(args *ArgsGetBindByUser) (dataList []DataGetBindByUserMarge, err error) {
	var filterResult []DataGetBindByUserMarge
	//获取数据集合
	var bindList []FieldsBind
	bindList, err = GetBindByUser(args)
	if err != nil {
		return
	}
	//重组数据
	var orgIDs []int64
	for _, v := range bindList {
		orgIDs = append(orgIDs, v.OrgID)
	}
	var orgList []FieldsOrg
	orgList, err = GetOrgMore(&ArgsGetOrgMore{
		IDs:        orgIDs,
		HaveRemove: false,
	})
	if err != nil {
		return
	}
	//遍历组织成员清单
	for _, v := range bindList {
		var vOrgData FieldsOrg
		isFind := false
		for _, v2 := range orgList {
			if v2.ID == v.OrgID {
				vOrgData = v2
				isFind = true
				break
			}
		}
		if !isFind {
			continue
		}
		var orgCoverFileID int64
		orgCoverFileID, err = Config.GetConfigValInt64(&ClassConfig.ArgsGetConfig{
			BindID:    vOrgData.ID,
			Mark:      "CoverFileID",
			VisitType: "public",
		})
		if err != nil {
			err = nil
			orgCoverFileID = 0
		}
		var groupList []FieldsGroup
		var groups []DataGetBindByUserMargeGroup
		groupList, err = getGroups(v.GroupIDs)
		if err != nil {
			err = nil
		} else {
			for _, vGroup := range groupList {
				groups = append(groups, DataGetBindByUserMargeGroup{
					ID:   vGroup.ID,
					Name: vGroup.Name,
				})
			}
		}
		//重组权限
		vManager := GetPermissionByBindID(v.ID)
		//写入数据
		filterResult = append(filterResult, DataGetBindByUserMarge{
			OrgID:          vOrgData.ID,
			OrgUserID:      vOrgData.UserID,
			OrgKey:         vOrgData.Key,
			OrgName:        vOrgData.Name,
			OrgDes:         vOrgData.Des,
			OrgCoverFileID: orgCoverFileID,
			OrgParentID:    vOrgData.ParentID,
			OrgParentFunc:  vOrgData.ParentFunc,
			OrgOpenFunc:    vOrgData.OpenFunc,
			BindID:         v.ID,
			LastAt:         v.LastAt,
			Avatar:         v.Avatar,
			Name:           v.Name,
			Manager:        vManager,
			Groups:         groups,
			Params:         v.Params,
		})
	}
	//去除dataList重复的OrgID项目
	for _, v := range filterResult {
		isFind := false
		for _, v2 := range dataList {
			if v.OrgID == v2.OrgID {
				isFind = true
				break
			}
		}
		if !isFind {
			dataList = append(dataList, v)
		}
	}
	//反馈
	return
}

// getBindAllByUserID 获取所有用户绑定的成员
func getBindAllByUserID(userID int64) (dataList []FieldsBind) {
	var rawList []FieldsBind
	err := Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM org_core_bind WHERE user_id = $1 AND delete_at < to_timestamp(1000000)", userID)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getBindByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetBindIDsByGroupID 获取指定分组所有人员ID列参数
type ArgsGetBindIDsByGroupID struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//分组ID
	GroupID int64 `json:"groupID" check:"id" empty:"true"`
}

// GetBindIDsByGroupID 获取指定分组所有人员ID列
func GetBindIDsByGroupID(args *ArgsGetBindIDsByGroupID) (ids pq.Int64Array, err error) {
	var dataList []FieldsBind
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id FROM org_core_bind WHERE org_id = $1 AND ($2 < 1 OR $2 = ANY(group_ids))", args.OrgID, args.GroupID)
	if err != nil || len(dataList) < 1 {
		err = errors.New("no data")
		return
	}
	for _, v := range dataList {
		ids = append(ids, v.ID)
	}
	return
}

// GetBindID 通过ID查询成员信息
func GetBindID(orgID int64, bindID int64) (data FieldsBind) {
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_core_bind WHERE ($1 < 1 OR org_id = $1) AND id = $2", orgID, bindID)
	if err != nil || data.ID < 1 {
		data = FieldsBind{}
		return
	}
	data = getBindByID(data.ID)
	return
}

// DataGetBindMoreName 获取一组成员姓名数据
type DataGetBindMoreName struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// GetBindMoreName 获取一组成员姓名
func GetBindMoreName(orgID int64, ids pq.Int64Array) (dataList []DataGetBindMoreName) {
	var rawList []FieldsBind
	_ = Router2SystemConfig.MainDB.Select(&rawList, "SELECT id, name FROM org_core_bind WHERE ($1 < 0 OR org_id = $1) AND id = ANY($2)", orgID, ids)
	if len(rawList) < 1 {
		return
	}
	for _, v := range rawList {
		dataList = append(dataList, DataGetBindMoreName{
			ID:   v.ID,
			Name: v.Name,
		})
	}
	return
}

func GetBindMoreNameMap(orgID int64, ids pq.Int64Array) (dataList map[int64]string) {
	rawList := GetBindMoreName(orgID, ids)
	dataList = map[int64]string{}
	for _, v := range rawList {
		dataList[v.ID] = v.Name
	}
	return
}

// GetBindIDsByManager 获取符合权限的所有成员
func GetBindIDsByManager(orgID int64, manager string) (ids pq.Int64Array, err error) {
	var groupList []FieldsGroup
	_ = Router2SystemConfig.MainDB.Select(&groupList, "SELECT id FROM org_core_group WHERE org_id = $1 AND $2 = ANY(manager) AND delete_at < to_timestamp(1000000)", orgID, manager)
	var dataList []FieldsBind
	if len(groupList) > 0 {
		var groupIDs pq.Int64Array
		for _, v := range groupList {
			groupIDs = append(groupIDs, v.ID)
		}
		err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id FROM org_core_bind WHERE org_id = $1 AND ($2 = ANY(manager) OR $3 @> group_ids) AND delete_at < to_timestamp(1000000)", orgID, manager, groupIDs)
	} else {
		err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id FROM org_core_bind WHERE org_id = $1 AND $2 = ANY(manager) AND delete_at < to_timestamp(1000000)", orgID, manager)
	}
	for _, v := range dataList {
		ids = append(ids, v.ID)
	}
	return
}

// GetBindByPhone 通过手机号查询成员
func GetBindByPhone(orgID int64, nationCode string, phone string) (data FieldsBind) {
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_core_bind WHERE org_id = $1 AND nation_code = $2 AND phone = $3 AND delete_at < to_timestamp(1000000)", orgID, nationCode, phone)
	if err != nil || data.ID < 1 {
		data = FieldsBind{}
		return
	}
	data = getBindByID(data.ID)
	return
}

// getBindAllByPhone 获取所有绑定手机号的成员
func getBindAllByPhone(nationCode string, phone string) (dataList []FieldsBind) {
	var rawList []FieldsBind
	err := Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM org_core_bind WHERE nation_code = $1 AND phone = $2 AND delete_at < to_timestamp(1000000)", nationCode, phone)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getBindByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// GetBindByEmail 通过邮箱查询成员
func GetBindByEmail(orgID int64, email string) (data FieldsBind) {
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_core_bind WHERE org_id = $1 AND email = $2 AND delete_at < to_timestamp(1000000)", orgID, email)
	if err != nil || data.ID < 1 {
		data = FieldsBind{}
		return
	}
	data = getBindByID(data.ID)
	return
}

// GetBindBySync 通过同步数据查询成员
func GetBindBySync(orgID int64, syncSystem string, syncID int64, syncHash string) (data FieldsBind) {
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_core_bind WHERE org_id = $1 AND sync_system = $2 AND sync_id = $3 AND sync_hash = $4 AND delete_at < to_timestamp(1000000)", orgID, syncSystem, syncID, syncHash)
	if err != nil || data.ID < 1 {
		data = FieldsBind{}
		return
	}
	data = getBindByID(data.ID)
	return
}

// GetBindCountByOrgID 获取组织下成员数量
func GetBindCountByOrgID(orgID int64) (count int64) {
	_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM org_core_bind WHERE org_id = $1 AND delete_at < to_timestamp(1000000)", orgID)
	return
}

// ArgsCheckBindAndOrg 检查绑定关系是否和组织一致参数
type ArgsCheckBindAndOrg struct {
	//绑定ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
}

// CheckBindAndOrg 检查绑定关系是否和组织一致
func CheckBindAndOrg(args *ArgsCheckBindAndOrg) (err error) {
	type dataType struct {
		//ID
		ID int64 `db:"id" json:"id"`
	}
	var data dataType
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_core_bind WHERE id = $1 AND org_id = $2", args.ID, args.OrgID)
	if err == nil && data.ID < 1 {
		err = errors.New("not exist")
	}
	return
}

// GetBindByGroupAndRole 获取同时满足分组和角色的成员
func GetBindByGroupAndRole(orgID int64, groupIDs, roleIDs pq.Int64Array) (dataList []FieldsBind) {
	var rawList []FieldsBind
	err := Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM org_core_bind WHERE org_id = $1 AND delete_at < to_timestamp(1000000) AND group_ids && $2 AND role_config_ids && $3", orgID, groupIDs, roleIDs)
	if err != nil || len(rawList) < 1 {
		return
	}
	for _, v := range rawList {
		dataList = append(dataList, getBindByID(v.ID))
	}
	return
}

// 检查用户是否具备行政权限，并自动授权
func checkAndUpdateUserGroup(userID int64) (err error) {
	//检查用户所属组
	var userData UserCore.FieldsUserType
	userData, err = UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
		ID:    userID,
		OrgID: -1,
	})
	if err != nil {
		err = errors.New(fmt.Sprint("user not exist, ", err))
		return
	}
	var userOrgDefaultGroupID int64
	userOrgDefaultGroupID, err = BaseConfig.GetDataInt64("UserOrgDefaultGroupID")
	if err == nil && userOrgDefaultGroupID > 0 {
		isFind := false
		for _, v := range userData.Groups {
			if v.GroupID == userOrgDefaultGroupID && v.ExpireAt.Unix() > CoreFilter.GetNowTime().Unix() {
				isFind = true
				break
			}
		}
		if !isFind {
			err = UserCore.UpdateUserGroupByID(&UserCore.ArgsUpdateUserGroupByID{
				ID:       userData.ID,
				OrgID:    -1,
				GroupID:  userOrgDefaultGroupID,
				ExpireAt: CoreFilter.GetNowTimeCarbon().AddYears(999).Time,
				IsRemove: false,
			})
			if err != nil {
				err = errors.New(fmt.Sprint("update user group, ", err))
				return
			}
		}
	}
	return nil
}

// 获取指定ID
func getBindByID(id int64) (data FieldsBind) {
	cacheMark := getBindCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, last_at, user_id, name, org_id, group_ids, manager, params, avatar, nation_code, phone, email, sync_system, sync_id, sync_hash, role_config_ids FROM org_core_bind WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, bindCacheTime)
	return
}
