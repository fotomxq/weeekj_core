package OrgCoreCore

import (
	"errors"
	"fmt"
	AnalysisAny2 "github.com/fotomxq/weeekj_core/v5/analysis/any2"
	CoreNats "github.com/fotomxq/weeekj_core/v5/core/nats"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLIDs "github.com/fotomxq/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	CoreSQL2 "github.com/fotomxq/weeekj_core/v5/core/sql2"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	UserCore "github.com/fotomxq/weeekj_core/v5/user/core"
	"github.com/lib/pq"
	"time"
)

// ArgsGetOrgList 获取列表参数
type ArgsGetOrgList struct {
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
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetOrgList 获取列表
func GetOrgList(args *ArgsGetOrgList) (dataList []FieldsOrg, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	if args.ParentID > -1 {
		where = where + "parent_id = :parent_id"
		maps["parent_id"] = args.ParentID
	}
	if args.IsRemove {
		if where != "" {
			where = where + " AND "
		}
		where = where + "delete_at > to_timestamp(1000000)"
	} else {
		if where != "" {
			where = where + " AND "
		}
		where = where + "delete_at < to_timestamp(1000000)"
	}
	if args.SortID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "sort_id = :sort_id"
		maps["sort_id"] = args.SortID
	}
	if args.UserID > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "user_id = :user_id"
		maps["user_id"] = args.UserID
	}
	if len(args.ParentFunc) > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "parent_func @> :parent_func"
		maps["parent_func"] = args.ParentFunc
	}
	if len(args.OpenFunc) > 0 {
		if where != "" {
			where = where + " AND "
		}
		where = where + "open_func @> :open_func"
		maps["open_func"] = args.OpenFunc
	}
	if args.Search != "" {
		if where != "" {
			where = where + " AND "
		}
		where = where + "(name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	if where == "" {
		where = "true"
	}
	var rawList []FieldsOrg
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		"org_core",
		"id",
		"SELECT id FROM org_core "+"WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at", "name"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getOrgByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// GetOrgByUser 获取用户拥有的所有商户
func GetOrgByUser(userID int64) (dataList []FieldsOrg) {
	err := orgSQL.Select().SetFieldsList([]string{"id"}).SelectList("user_id = $1 AND delete_at < to_timestamp(1000000)", userID).Result(&dataList)
	if err != nil || len(dataList) < 1 {
		return
	}
	for k, v := range dataList {
		dataList[k] = getOrgByID(v.ID)
	}
	return
}

// GetOrgListStep 内部遍历组织专用
func GetOrgListStep(page, max int64, openFunc pq.StringArray) (dataList []FieldsOrg) {
	var rawList []FieldsOrg
	_ = orgSQL.Select().SetFieldsList([]string{"id"}).SetFieldsSort([]string{"id"}).SetPages(CoreSQL2.ArgsPages{
		Page: page,
		Max:  max,
		Sort: "id",
		Desc: false,
	}).SelectList("delete_at < to_timestamp(1000000) AND open_func @> $1", openFunc).Result(&rawList)
	if len(rawList) < 1 {
		return
	}
	for _, v := range rawList {
		dataList = append(dataList, getOrgByID(v.ID))
	}
	return
}

// ArgsGetOrgSearch 查询组织方法参数
type ArgsGetOrgSearch struct {
	//上级ID锁定
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//要查询的ID列
	// 和搜素互斥，更优先查询该数据
	IDs pq.Int64Array `json:"ids" check:"ids" empty:"true"`
	//搜索
	// 查询名称、描述信息
	Search string `json:"search" check:"search" empty:"true"`
}

// DataGetOrgSearch 查询组织方法数据
type DataGetOrgSearch struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//构架名称，或组织名称
	Name string `db:"name" json:"name"`
	//组织描述
	Des string `db:"des" json:"des"`
}

// GetOrgSearch 查询组织方法
// 该方法可以锁定组织，也可以不锁定
// 不锁定则可以全局查询；锁定则只能在组织内查询
// 最多反馈100条数据
func GetOrgSearch(args *ArgsGetOrgSearch) (dataList []DataGetOrgSearch, err error) {
	var rawList []FieldsOrg
	if len(args.IDs) > 0 {
		if args.OrgID > 0 {
			err = Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM org_core WHERE parent_id = $1 AND id = ANY($2) AND delete_at < to_timestamp(1000000) LIMIT 10", args.OrgID, args.IDs)
		} else {
			err = Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM org_core WHERE id = ANY($1) AND delete_at < to_timestamp(1000000) LIMIT 10", args.IDs)
		}
		if err != nil {
			return
		}
	} else {
		where := "delete_at < to_timestamp(1000000)"
		maps := map[string]interface{}{}
		if args.OrgID > 0 {
			where = where + " AND parent_id = :parent_id"
			maps["parent_id"] = args.OrgID
		}
		if args.Search != "" {
			where = where + " AND (name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
			maps["search"] = args.Search
		}
		err = CoreSQL.GetList(
			Router2SystemConfig.MainDB.DB,
			&rawList,
			"SELECT id FROM org_core "+"WHERE "+where+" LIMIT 10",
			maps,
		)
		if err != nil {
			return
		}
	}
	for _, v := range rawList {
		vData := getOrgByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, DataGetOrgSearch{
			ID:   vData.ID,
			Name: vData.Name,
			Des:  vData.Des,
		})
	}
	return
}

// ArgsGetOrg 查看组织参数
type ArgsGetOrg struct {
	//ID
	ID int64 `json:"id" check:"id"`
}

// GetOrg 查看组织
func GetOrg(args *ArgsGetOrg) (data FieldsOrg, err error) {
	data = getOrgByID(args.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

func GetOrgByID(id int64) (data FieldsOrg) {
	if id < 1 {
		return
	}
	data = getOrgByID(id)
	if data.ID < 1 {
		return
	}
	return
}

func GetOrgName(args *ArgsGetOrg) (data string, err error) {
	if args.ID < 1 {
		err = errors.New("no data")
		return
	}
	orgData := getOrgByID(args.ID)
	if orgData.ID < 1 {
		err = errors.New("no data")
		return
	}
	data = orgData.Name
	return
}

func GetOrgNameByID(id int64) string {
	if id < 1 {
		return ""
	}
	cacheMark := getOrgNameCacheMark(id)
	result, err := Router2SystemConfig.MainCache.GetString(cacheMark)
	if err == nil {
		return result
	}
	data := getOrgByID(id)
	if data.ID < 1 {
		return ""
	}
	Router2SystemConfig.MainCache.SetString(cacheMark, data.Name, 3600)
	return data.Name
}

func GetOrgKeyByID(id int64) string {
	if id < 1 {
		return ""
	}
	data := getOrgByID(id)
	if data.ID < 1 {
		return ""
	}
	return data.Key
}

// GetOrgByName 根据名称获取企业
func GetOrgByName(name string, parentID int64) (orgData FieldsOrg) {
	err := Router2SystemConfig.MainDB.Get(&orgData, "SELECT id FROM org_core WHERE name = $1 AND parent_id = $2 AND delete_at < to_timestamp(1000000)", name, parentID)
	if err != nil || orgData.ID < 1 {
		return
	}
	orgData = getOrgByID(orgData.ID)
	return
}

// ArgsGetOrgByKey 通过key获取企业参数
type ArgsGetOrgByKey struct {
	//企业唯一标识码
	// 用于特殊识别和登陆识别等操作
	Key string `db:"key" json:"key" check:"key"`
}

// GetOrgByKey 通过key获取企业
func GetOrgByKey(args *ArgsGetOrgByKey) (data FieldsOrg, err error) {
	if args.Key == "" {
		err = errors.New("key is empty")
		return
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM org_core WHERE key = $1", args.Key)
	if err != nil {
		return
	}
	data = getOrgByID(data.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsGetOrgMore 通过一组ID查询组织参数
type ArgsGetOrgMore struct {
	//ID列
	IDs pq.Int64Array `json:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove"`
}

// GetOrgMore 通过一组ID查询组织
// 反馈组织数量，最多不能超出100个
func GetOrgMore(args *ArgsGetOrgMore) (dataList []FieldsOrg, err error) {
	var rawList []FieldsOrg
	err = CoreSQLIDs.GetIDsAndDelete(&rawList, "org_core", "id", args.IDs, args.HaveRemove)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getOrgByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// GetOrgMoreName 获取组织名称信息
func GetOrgMoreName(args *ArgsGetOrgMore) (data map[int64]string, err error) {
	data, err = CoreSQLIDs.GetIDsNameAndDelete("org_core", args.IDs, args.HaveRemove)
	return
}

// ArgsGetOrgChildCount 获取子组织数量参数
type ArgsGetOrgChildCount struct {
	//上级组织ID
	ParentOrgID int64 `json:"parentOrgID" check:"id"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
}

// GetOrgChildCount 获取子组织数量
func GetOrgChildCount(args *ArgsGetOrgChildCount) (count int64, err error) {
	err = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) as count FROM org_core WHERE parent_id = $1 AND ($2 = true OR delete_at < to_timestamp(1000000))", args.ParentOrgID, args.HaveRemove)
	return
}

// AddOrgUserVisit 增加组织访问人数
func AddOrgUserVisit(orgID int64) {
	if orgID < 1 {
		return
	}
	AnalysisAny2.AppendData("add", "org_user_visit_count", time.Time{}, orgID, 0, 0, 0, 0, 1)
}

// ArgsUpdateOrgUserID 修改组织的所有归属权参数
type ArgsUpdateOrgUserID struct {
	//组织ID
	ID int64 `json:"id" check:"id"`
	//上级关系
	// 可选，用于筛选
	ParentID int64 `json:"parentID" check:"id" empty:"true"`
	//原用户ID
	OldUserID int64 `json:"oldUserID" check:"id"`
	//修改目标用户ID
	UserID int64 `json:"userID" check:"id"`
	//是否修改子组织所有人
	// 注意，子组织所有人原先必须一致，如果是其他人则不会修改
	// 本设计只能修改下一级，多级需进入下一级继续修改
	AllowChild bool `json:"allowChild" check:"bool"`
}

// UpdateOrgUserID 修改组织的所有归属权
// 将修改本组织和旗下所有子组织的所属权
// 本操作只能管理员或原所有人用户修改
func UpdateOrgUserID(args *ArgsUpdateOrgUserID) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
			return
		}
	}()
	var data FieldsOrg
	data, err = GetOrg(&ArgsGetOrg{
		ID: args.ID,
	})
	if err != nil {
		return
	}
	if args.ParentID > 0 && data.ParentID != args.ParentID {
		err = errors.New("no child org operate permission")
		return
	}
	if args.ID == args.ParentID {
		err = errors.New("child id is now org id")
		return
	}
	if args.OldUserID != data.UserID {
		err = errors.New("user not old user")
		return
	}
	//新用户ID
	var newUserData UserCore.FieldsUserType
	newUserData, err = UserCore.GetUserByID(&UserCore.ArgsGetUserByID{
		ID:    args.UserID,
		OrgID: -1,
	})
	if err != nil {
		return
	}
	//开始转移商户数据
	tx := Router2SystemConfig.MainDB.MustBegin()
	tx.MustExec("UPDATE org_core SET update_at = NOW(), user_id = $1 WHERE id = $2;", args.UserID, args.ID)
	var childList []FieldsOrg
	err = Router2SystemConfig.MainDB.Select(&childList, "SELECT id FROM org_core WHERE parent_id = $1", data.ID)
	if err == nil {
		var childIDs pq.Int64Array
		for _, v := range childList {
			childIDs = append(childIDs, v.ID)
		}
		if len(childIDs) > 0 {
			tx.MustExec("UPDATE org_core SET update_at = NOW(), user_id = $1 WHERE id = ANY($2);", args.UserID, childIDs)
		}
	}
	err = tx.Commit()
	if err != nil {
		if err = tx.Rollback(); err != nil {
			err = errors.New("roll back failed, " + err.Error())
			return
		}
		err = errors.New("update failed, " + err.Error())
		return
	}
	//如果旧的用户和新的用户不同
	if args.OldUserID != args.UserID {
		//删除旧用户的绑定关系
		err = DeleteBindByUser(&ArgsDeleteBindByUser{
			UserID: args.OldUserID,
			OrgID:  data.ID,
		})
		if err != nil {
			return
		}
		_, err = SetBind(&ArgsSetBind{
			UserID:     newUserData.ID,
			Avatar:     0,
			Name:       newUserData.Name,
			OrgID:      data.ID,
			GroupIDs:   []int64{},
			Manager:    []string{"member", "all"},
			NationCode: newUserData.NationCode,
			Phone:      newUserData.Phone,
			Email:      newUserData.Email,
			SyncSystem: "",
			SyncID:     0,
			SyncHash:   "",
			Params:     CoreSQLConfig.FieldsConfigsType{},
		})
		if err != nil {
			return
		}
	}
	deleteOrgCache(args.ID)
	return
}

// ArgsUpdateOrg 修改组织参数
type ArgsUpdateOrg struct {
	//组织ID
	ID int64 `db:"id" json:"id" check:"id"`
	//所属用户
	// 可选，用于验证
	UserID int64 `db:"user_id" json:"userID" check:"id"`
	//企业唯一标识码
	// 用于特殊识别和登陆识别等操作
	Key string `db:"key" json:"key" check:"mark" empty:"true"`
	//构架名称，或组织名称
	Name string `db:"name" json:"name" check:"name"`
	//组织描述
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
	//上级ID
	ParentID int64 `db:"parent_id" json:"parentID" check:"id" empty:"true"`
	//上级控制权限限制
	ParentFunc pq.StringArray `db:"parent_func" json:"parentFunc" check:"marks" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
}

// UpdateOrg 修改组织
func UpdateOrg(args *ArgsUpdateOrg) (errCode string, err error) {
	//禁止上级为自己
	if args.ID == args.ParentID {
		errCode = "parent_is_self"
		err = errors.New("child id is now org id")
		return
	}
	//生成key
	if args.Key == "" {
		args.Key = makeKey(args.Name)
	}
	//检查key
	if args.Key != "" {
		var data FieldsOrg
		data, err = GetOrgByKey(&ArgsGetOrgByKey{
			Key: args.Key,
		})
		if err == nil && data.ID > 0 && data.ID != args.ID {
			errCode = "key_exist"
			err = errors.New("key is exist")
			return
		}
	}
	//检查名称
	oldOrgID := getOrgByName(args.Name)
	if oldOrgID > 0 {
		if oldOrgID != args.ID {
			errCode = "name_exist"
			err = errors.New("name is exist")
			return
		}
	}
	//更新数据
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_core SET update_at = NOW(), key = :key, name = :name, des = :des, parent_id = :parent_id, parent_func = :parent_func, sort_id = :sort_id WHERE id = :id AND (:user_id < 1 OR user_id = :user_id);", args)
	if err != nil {
		errCode = "update"
		return
	}
	deleteOrgCache(args.ID)
	return
}

// ArgsUpdateOrgFunc 修改组织开通业务参数
type ArgsUpdateOrgFunc struct {
	//组织ID
	ID int64 `db:"id" json:"id"`
	//开通业务
	// 该内容只有总管理员或订阅能进行控制
	OpenFunc pq.StringArray `db:"open_func" json:"openFunc"`
}

// UpdateOrgFunc 修改组织开通业务
func UpdateOrgFunc(args *ArgsUpdateOrgFunc) (err error) {
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE org_core SET update_at = NOW(), open_func = :open_func WHERE id = :id;", args)
	if err != nil {
		return
	}
	deleteOrgCache(args.ID)
	return
}

// ArgsDeleteOrg 删除组织参数
type ArgsDeleteOrg struct {
	//组织ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteOrg 删除组织
func DeleteOrg(args *ArgsDeleteOrg) (err error) {
	_, err = CoreSQL.DeleteOneSoft(Router2SystemConfig.MainDB.DB, "org_core", "id", args)
	if err != nil {
		return
	}
	_ = DeleteBindByOrg(&ArgsDeleteBindByOrg{
		OrgID: args.ID,
	})
	deleteOrgCache(args.ID)
	//推送NATS
	CoreNats.PushDataNoErr("org_core_org", "/org/core/org", "delete", args.ID, "", nil)
	return
}

// ArgsDeleteOrgParent 删除下级组织参数
type ArgsDeleteOrgParent struct {
	//组织ID
	ID int64 `db:"id" json:"id" check:"id"`
	//上级ID
	ParentID int64 `db:"parent_id" json:"parentID" check:"id"`
}

// DeleteOrgParent 删除下级组织
func DeleteOrgParent(args *ArgsDeleteOrgParent) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "org_core", "id = :id AND parent_id = :parent_id", args)
	if err != nil {
		return
	}
	deleteOrgCache(args.ID)
	return
}

// checkOrgInParentFunc 检查一个级别功能是否在另外一个范围内
// 只要有一个超出，则禁止
func checkOrgInParentFunc(parentFunc []string, childFunc []string) (err error) {
	//遍历子组织的功能
	for _, v := range childFunc {
		isFind := false
		for _, v2 := range parentFunc {
			if v == v2 {
				isFind = true
				break
			}
		}
		if isFind {
			continue
		}
		//如果没有找到，则说明子未包含在上级，禁止访问
		err = errors.New("func not in parent area")
		return
	}
	return
}

// checkOrgParent 递归检查，是否存在死循环
func checkOrgParent(id int64, parentID int64, checkParentIDs []int64) (err error) {
	checkParentIDs = append(checkParentIDs, id)
	if id == parentID {
		return errors.New("parent id is cycle")
	}
	var parentData FieldsOrg
	parentData, err = GetOrg(&ArgsGetOrg{
		ID: parentID,
	})
	if err != nil {
		//上级不存在
		err = errors.New("parent is not exist")
		return
	}
	//上级存在，则检查其上级是否在序列内？
	for _, v := range checkParentIDs {
		if v == parentID {
			err = errors.New("parent id is cycle")
			return
		}
	}
	//不存在上级，则跳出
	if parentID < 1 {
		return nil
	}
	//不存在则继续
	if err = checkOrgParent(parentData.ID, parentData.ParentID, checkParentIDs); err != nil {
		return err
	}
	return nil
}

// 根据商户名称找到商户
func getOrgByName(name string) (orgID int64) {
	_ = Router2SystemConfig.MainDB.Get(&orgID, "SELECT id FROM org_core WHERE name = $1 AND delete_at < to_timestamp(1000000)", name)
	return
}

// 获取指定组织信息
func getOrgByID(id int64) (data FieldsOrg) {
	cacheMark := getOrgCacheMark(id)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, user_id, key, name, des, parent_id, parent_func, open_func, sort_id FROM org_core WHERE id = $1", id)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, orgCacheTime)
	return
}

// 缓冲
func getOrgCacheMark(id int64) string {
	return fmt.Sprint("org:core:org:id:", id)
}

func getOrgNameCacheMark(id int64) string {
	return fmt.Sprint("org:core:org:name:", id)
}

func deleteOrgCache(id int64) {
	Router2SystemConfig.MainCache.DeleteMark(getOrgCacheMark(id))
	Router2SystemConfig.MainCache.DeleteMark(getOrgNameCacheMark(id))
	deletePermissionByOrgCache(id)
}
