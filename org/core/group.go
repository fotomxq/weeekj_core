package OrgCoreCore

import (
	"errors"
	"fmt"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLIDs "github.com/fotomxq/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetGroupList 获取分组数据参数
type ArgsGetGroupList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `json:"orgID" check:"id"`
	//上级部门
	ParentID int64 `db:"parent_id" json:"parentID"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetGroupList 获取分组数据
func GetGroupList(args *ArgsGetGroupList) (dataList []FieldsGroup, dataCount int64, err error) {
	where := "(org_id = :org_id or org_id = 0)"
	maps := map[string]interface{}{
		"org_id": args.OrgID,
	}
	if args.IsRemove {
		where = where + " AND delete_at > to_timestamp(1000000)"
	} else {
		where = where + " AND delete_at < to_timestamp(1000000)"
	}
	if args.ParentID > -1 {
		where = where + " AND parent_id = :parent_id"
		maps["parent_id"] = args.ParentID
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	var rawList []FieldsGroup
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		"org_core_group",
		"id",
		"SELECT id FROM org_core_group "+"WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getGroupByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// GetAllGroup 获取所有分组
func GetAllGroup(orgID int64) (dataList []FieldsGroup) {
	var rawList []FieldsGroup
	err := CoreSQL.GetList(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		"SELECT id FROM org_core_group org_id = :org_id AND delete_at < to_timestamp(1000000)",
		map[string]interface{}{
			"org_id": orgID,
		},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getGroupByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// 获取符合权限的分组
func getGroupByManager(orgID int64, managers pq.StringArray) (dataList []FieldsGroup) {
	var rawList []FieldsGroup
	err := Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM org_core_group WHERE org_id = $1 AND delete_at < to_timestamp(1000000) AND manager = ANY($2)", orgID, managers)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getGroupByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetGroup 获取某个分组参数
type ArgsGetGroup struct {
	//ID
	ID int64 `json:"id" check:"id"`
	//组织ID
	// 验证是否一致
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
}

// GetGroup 获取某个分组
func GetGroup(args *ArgsGetGroup) (data FieldsGroup, err error) {
	data = getGroupByID(args.ID)
	if data.ID < 1 || !CoreFilter.EqID2(args.OrgID, data.OrgID) || CoreSQL.CheckTimeHaveData(data.DeleteAt) {
		err = errors.New("no data")
		return
	}
	return
}

// GetGroupByName 通过名字获取分组
func GetGroupByName(orgID int64, name string) (data FieldsGroup) {
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_core_group WHERE org_id = $1 AND delete_at < to_timestamp(1000000) AND name = $2", orgID, name)
	if err != nil {
		return
	}
	data = getGroupByID(data.ID)
	return
}

// GetGroupNameByID 通过ID查询名字
func GetGroupNameByID(orgID int64, id int64) (name string) {
	data := getGroupByID(id)
	if !CoreFilter.EqID2(orgID, data.OrgID) {
		return
	}
	name = data.Name
	return
}

// ArgsGetGroupMore 获取一组分组参数
type ArgsGetGroupMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
}

// GetGroupMore 获取一组分组
func GetGroupMore(args *ArgsGetGroupMore) (dataList []FieldsGroup, err error) {
	var rawList []FieldsGroup
	err = CoreSQLIDs.GetIDsAndDelete(&rawList, "org_core_group", "id", args.IDs, args.HaveRemove)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getGroupByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

func GetGroupMoreNames(args *ArgsGetGroupMore) (data map[int64]string, err error) {
	data, err = CoreSQLIDs.GetIDsNameAndDelete("org_core_group", args.IDs, args.HaveRemove)
	return
}

func GetGroupMoreNamesNoErr(ids []int64) (data map[int64]string) {
	data, _ = GetGroupMoreNames(&ArgsGetGroupMore{
		IDs:        ids,
		HaveRemove: false,
	})
	return
}

func GetGroupMoreNamesStr(args *ArgsGetGroupMore) (data []string, err error) {
	var rawData map[int64]string
	rawData, err = GetGroupMoreNames(args)
	if err != nil {
		return
	}
	for k := 0; k < len(args.IDs); k++ {
		for k2, v2 := range rawData {
			if args.IDs[k] == k2 {
				continue
			}
			data = append(data, v2)
		}
	}
	return
}

func GetGroupMoreNamesStrNoErr(groupIDs pq.Int64Array) (data []string) {
	for k := 0; k < len(groupIDs); k++ {
		v := groupIDs[k]
		if v < 1 {
			continue
		}
		vData := getGroupByID(v)
		if vData.ID < 1 {
			continue
		}
		data = append(data, vData.Name)
	}
	return
}

// ArgsCheckGroupsHavePermission 检查一组分组内，是否存在某个权限？参数
type ArgsCheckGroupsHavePermission struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//权限
	Manager string `json:"manager" check:"mark"`
}

// CheckGroupsHavePermission 检查一组分组内，是否存在某个权限？
func CheckGroupsHavePermission(args *ArgsCheckGroupsHavePermission) (groupIDs pq.Int64Array, b bool) {
	var dataList []FieldsGroup
	err := CoreSQLIDs.GetIDsAndDelete(&dataList, "org_core_group", "id", args.IDs, false)
	if err != nil {
		return
	}
	for _, v := range dataList {
		vData := getGroupByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		for _, v2 := range vData.Manager {
			if v2 == args.Manager {
				groupIDs = append(groupIDs, vData.ID)
				break
			}
		}
	}
	if len(groupIDs) > 0 {
		b = true
	}
	return
}

// ArgsCreateGroup 创建分组参数
type ArgsCreateGroup struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID"`
	//分组名称
	Name string `db:"name" json:"name"`
	//权利主张
	Manager pq.StringArray `db:"manager" json:"manager"`
	//部门领导
	ManagerOrgBindID int64 `db:"manager_org_bind_id" json:"managerOrgBindID"`
	//上级部门
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateGroup 创建分组
func CreateGroup(args *ArgsCreateGroup) (data FieldsGroup, err error) {
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "org_core_group", "INSERT INTO org_core_group (org_id, name, manager, manager_org_bind_id, parent_id, params) VALUES (:org_id,:name,:manager,:manager_org_bind_id,:parent_id,:params)", args, &data)
	if err != nil {
		return
	}
	return
}

// ArgsUpdateGroup 修改分组参数
type ArgsUpdateGroup struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//组织ID
	// 无法修改，用于验证
	OrgID int64 `db:"org_id" json:"orgID"`
	//分组名称
	Name string `db:"name" json:"name"`
	//权利主张
	Manager pq.StringArray `db:"manager" json:"manager"`
	//部门领导
	ManagerOrgBindID int64 `db:"manager_org_bind_id" json:"managerOrgBindID"`
	//上级部门
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateGroup 修改分组
func UpdateGroup(args *ArgsUpdateGroup) (err error) {
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE org_core_group SET update_at = NOW(), name = :name, manager = :manager, manager_org_bind_id = :manager_org_bind_id, parent_id = :parent_id, params = :params WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	deleteGroupCache(args.ID)
	return
}

// ArgsDeleteGroup 删除分组参数
type ArgsDeleteGroup struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//组织ID
	// 无法修改，用于验证
	OrgID int64 `db:"org_id" json:"orgID"`
}

// DeleteGroup 删除分组
func DeleteGroup(args *ArgsDeleteGroup) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "org_core_group", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	if err != nil {
		return
	}
	deleteGroupCache(args.ID)
	return
}

// 获取多个分组
func getGroups(ids pq.Int64Array) (dataList []FieldsGroup, err error) {
	var rawList []FieldsGroup
	err = Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM org_core_group WHERE id = ANY($1) AND delete_at < to_timestamp(1000000);", ids)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getGroupByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// 获取ID
func getGroupByID(id int64) (data FieldsGroup) {
	cacheMark := getGroupCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, name, manager, manager_org_bind_id, parent_id, params FROM org_core_group WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, bindGroupCacheTime)
	return
}

// 获取缓冲
func getGroupCacheMark(id int64) string {
	return fmt.Sprint("org:core:group:id:", id)
}

func deleteGroupCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getGroupCacheMark(id))
}
