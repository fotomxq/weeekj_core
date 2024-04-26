package UserCore

import (
	"encoding/json"
	"errors"
	"fmt"
	BaseFileSys2 "github.com/fotomxq/weeekj_core/v5/base/filesys2"
	CoreFilter "github.com/fotomxq/weeekj_core/v5/core/filter"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLIDs "github.com/fotomxq/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetUserList 获取用户列表参数
type ArgsGetUserList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//状态
	Status int `json:"status" check:"intThan0" empty:"true"`
	//组织ID
	// 可以留空，则表明为平台
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//上级系统来源
	ParentSystem string `json:"parentSystem" check:"mark" empty:"true"`
	ParentID     int64  `json:"parentID" check:"id" empty:"true"`
	//分类
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//是否删除
	IsRemove bool `json:"isRemove" check:"bool" empty:"true"`
	//搜素
	Search string `json:"search" check:"search" empty:"true"`
}

// GetUserList 获取用户列表
func GetUserList(args *ArgsGetUserList) (dataList []FieldsUserType, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > 0 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.Status > -1 && args.Status < 3 {
		where = where + " AND status = :status"
		maps["status"] = args.Status
	}
	if args.ParentSystem != "" {
		var parents []byte
		if args.ParentID > 0 {
			type parentArg struct {
				System   string `json:"system"`
				ParentID int64  `json:"parentID"`
			}
			parents, err = json.Marshal([]parentArg{
				{
					System:   args.ParentSystem,
					ParentID: args.ParentID,
				},
			})
			if err != nil {
				return
			}
		} else {
			type parentArg struct {
				System string `json:"system"`
			}
			parents, err = json.Marshal([]parentArg{
				{
					System: args.ParentSystem,
				},
			})
			if err != nil {
				return
			}
		}
		where = where + " AND parents = :parents"
		maps["parents"] = string(parents)
	}
	if args.SortID > 0 {
		where = where + " AND sort_id = :sort_id"
		maps["sort_id"] = args.SortID
	}
	if len(args.Tags) > 0 {
		where = where + " AND tags @> :tags"
		maps["tags"] = args.Tags
	}
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%' OR phone ILIKE '%' || :search || '%' OR email ILIKE '%' || :search || '%' OR username ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	var rawList []FieldsUserType
	tableName := "user_core"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		fmt.Sprint("SELECT id ", "FROM ", tableName, " WHERE ", where),
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	if err != nil {
		return
	}
	for _, v := range rawList {
		vData := getUserByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		dataList = append(dataList, vData)
	}
	return
}

// ArgsGetOrgByPhone 获取手机号对应的所有组织参数
type ArgsGetOrgByPhone struct {
	//联系方式
	NationCode string `db:"nation_code" json:"nationCode" check:"nationCode"`
	Phone      string `db:"phone" json:"phone" check:"phone"`
}

// DataGetOrgByPhone 获取手机号对应的所有组织数据
type DataGetOrgByPhone struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 如果为空，则说明是平台的用户；否则为对应组织的用户
	// 所有获取的方法，都需要给与该ID参数，也可以留空，否则禁止获取
	OrgID int64 `db:"org_id" json:"orgID"`
}

// GetOrgByPhone 获取手机号对应的所有组织
func GetOrgByPhone(args *ArgsGetOrgByPhone) (dataList []DataGetOrgByPhone, err error) {
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, org_id FROM user_core WHERE nation_code = $1 AND phone = $2 AND delete_at < to_timestamp(1000000) ORDER BY id LIMIT 1", args.NationCode, args)
	return
}

// ArgsGetUserByID 获取某个用户信息参数
type ArgsGetUserByID struct {
	//用户ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可选，用于验证
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

func GetUserByID(args *ArgsGetUserByID) (data FieldsUserType, err error) {
	data = getUserByID(args.ID)
	if data.ID < 1 || !CoreFilter.EqID2(args.OrgID, data.OrgID) || CoreSQL.CheckTimeHaveData(data.DeleteAt) {
		data = FieldsUserType{}
		err = errors.New("no data")
		return
	}
	return
}

func GetUserByIDHaveDelete(args *ArgsGetUserByID) (data FieldsUserType, err error) {
	data = getUserByID(args.ID)
	if data.ID < 1 || !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		data = FieldsUserType{}
		err = errors.New("no data")
		return
	}
	return
}

// GetUserByIDNoErr 直接查询用户ID
func GetUserByIDNoErr(userID int64, orgID int64) (data FieldsUserType) {
	data, _ = GetUserByID(&ArgsGetUserByID{
		ID:    userID,
		OrgID: orgID,
	})
	return
}

// GetUserByIDNoErrHaveDelete 直接查询用户ID
func GetUserByIDNoErrHaveDelete(userID int64, orgID int64) (data FieldsUserType) {
	data, _ = GetUserByIDHaveDelete(&ArgsGetUserByID{
		ID:    userID,
		OrgID: orgID,
	})
	return
}

// ArgsCheckUserExist 检查用户是否存在参数
type ArgsCheckUserExist struct {
	//用户ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 可以留空，则表明为平台
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// CheckUserExist 检查用户是否存在
func CheckUserExist(args *ArgsCheckUserExist) (err error) {
	data := getUserByID(args.ID)
	if data.ID < 1 || CoreSQL.CheckTimeHaveData(data.DeleteAt) || !CoreFilter.EqID2(args.OrgID, data.OrgID) {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsGetUserByPhone 通过手机号找到唯一的用户参数
type ArgsGetUserByPhone struct {
	//组织ID
	// 可以留空，则表明为平台
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//联系方式
	NationCode string `db:"nation_code" json:"nationCode" check:"nationCode"`
	Phone      string `db:"phone" json:"phone" check:"phone"`
	//是否无视是否已经验证
	IgnoreVerify bool `db:"ignore_verify" json:"ignoreVerify" check:"bool" empty:"true"`
}

// GetUserByPhone 通过手机号找到唯一的用户
func GetUserByPhone(args *ArgsGetUserByPhone) (data FieldsUserType, err error) {
	if args.Phone == "" {
		err = errors.New("phone error")
		return
	}
	if args.IgnoreVerify {
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM user_core WHERE nation_code = $1 AND phone = $2 AND ($3 < 0 OR org_id = $3) AND delete_at < to_timestamp(1000000) ORDER BY id LIMIT 1", args.NationCode, args.Phone, args.OrgID)
	} else {
		err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM user_core WHERE nation_code = $1 AND phone = $2 AND ($3 < 0 OR org_id = $3) AND delete_at < to_timestamp(1000000) AND phone_verify > to_timestamp(1000000) ORDER BY id LIMIT 1", args.NationCode, args.Phone, args.OrgID)
	}
	if err != nil || data.ID < 1 {
		err = errors.New("no data")
		return
	}
	data = getUserByID(data.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// GetAllUserByPhone 获取已经被移除的关联的用户手机号
func GetAllUserByPhone(nationCode string, phone string) (userList []FieldsUserType) {
	if nationCode == "" || phone == "" {
		return
	}
	var rawList []FieldsUserType
	err := Router2SystemConfig.MainDB.Select(&rawList, "SELECT id FROM user_core WHERE nation_code = $1 AND phone = $2", nationCode, phone)
	if err != nil || len(rawList) < 1 {
		return
	}
	for _, v := range rawList {
		vData := getUserByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		userList = append(userList, vData)
	}
	return
}

// ArgsGetUserByEmail 通过邮箱找到用户参数
type ArgsGetUserByEmail struct {
	//组织ID
	// 可以留空，则表明为平台
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//邮箱
	Email string `db:"email" json:"email" check:"email"`
}

// GetUserByEmail 通过邮箱找到用户
func GetUserByEmail(args *ArgsGetUserByEmail) (data FieldsUserType, err error) {
	if args.Email == "" {
		err = errors.New("email error")
		return
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM user_core WHERE email = $1 AND email_verify > to_timestamp(1000000) AND ($2 < 0 OR org_id = $2) AND delete_at < to_timestamp(1000000) ORDER BY id LIMIT 1", args.Email, args.OrgID)
	if err != nil || data.ID < 1 {
		err = errors.New("no data")
		return
	}
	data = getUserByID(data.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsGetUserByUsername 通过用户名找到用户参数
type ArgsGetUserByUsername struct {
	//组织ID
	// 可以留空，则表明为平台
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//用户
	Username string `db:"username" json:"username" check:"username"`
}

// GetUserByUsername 通过用户名找到用户
func GetUserByUsername(args *ArgsGetUserByUsername) (data FieldsUserType, err error) {
	if args.Username == "" {
		err = errors.New("username error")
		return
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM user_core WHERE username = $1 AND (org_id = $2 OR $2 < 0) AND delete_at < to_timestamp(1000000) ORDER BY id LIMIT 1", args.Username, args.OrgID)
	if err != nil || data.ID < 1 {
		err = errors.New("no data")
		return
	}
	data = getUserByID(data.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsGetUserByInfos 通过infos找到用户参数
type ArgsGetUserByInfos struct {
	//组织ID
	// 可以留空，则表明为平台
	OrgID int64 `db:"org_id" json:"orgID"`
	//标识码
	Mark string `db:"mark" json:"mark"`
	//值
	Val string `db:"val" json:"val"`
}

// GetUserByInfos 通过infos找到用户
func GetUserByInfos(args *ArgsGetUserByInfos) (data FieldsUserType, err error) {
	var findJSON []byte
	type arg struct {
		//标识码
		Mark string `db:"mark" json:"mark"`
		//值
		Val string `db:"val" json:"val"`
	}
	findJSON, err = json.Marshal([]arg{
		{
			Mark: args.Mark,
			Val:  args.Val,
		},
	})
	if err != nil {
		return
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM user_core WHERE ($1 < 0 OR org_id = $1) AND infos @> $2 AND delete_at < to_timestamp(1000000)", args.OrgID, string(findJSON))
	if err != nil || data.ID < 1 {
		err = errors.New("no data")
		return
	}
	data = getUserByID(data.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsGetUserByLogin 通过扩展登陆信息结构体，找到用户参数
type ArgsGetUserByLogin struct {
	//组织ID
	// 可以留空，则表明为平台
	OrgID int64 `db:"org_id" json:"orgID"`
	//标识码
	Mark string `db:"mark" json:"mark"`
	//值
	Val string `db:"val" json:"val"`
}

// GetUserByLogin 通过扩展登陆信息结构体，找到用户
func GetUserByLogin(args *ArgsGetUserByLogin) (data FieldsUserType, err error) {
	if args.Mark == "" || args.Val == "" {
		err = errors.New("mark and val error")
		return
	}
	var findJSON []byte
	type arg struct {
		//标识码
		Mark string `db:"mark" json:"mark"`
		//值
		Val string `db:"val" json:"val"`
	}
	findJSON, err = json.Marshal([]arg{
		{
			Mark: args.Mark,
			Val:  args.Val,
		},
	})
	if err != nil {
		return
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM user_core WHERE ($1 < 0 OR org_id = $1) AND logins @> $2 AND delete_at < to_timestamp(1000000)", args.OrgID, string(findJSON))
	if err != nil || data.ID < 1 {
		err = errors.New("no data")
		return
	}
	data = getUserByID(data.ID)
	if data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsSearchUser 搜索用户
// 组件专用
type ArgsSearchUser struct {
	//最大反馈数量
	Max int64 `json:"max" check:"max"`
	//所属组织ID
	OrgID int64 `json:"orgID" check:"id" empty:"true"`
	//状态
	Status int `json:"status"`
	//上级系统及ID
	ParentSystem string `json:"parentSystem" check:"mark" empty:"true"`
	ParentID     int64  `json:"parentID" check:"id" empty:"true"`
	//是否包含删除的数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
	//搜索内容
	Search string `json:"search" check:"search" empty:"true"`
}

type DataSearchUser struct {
	//用户ID
	ID int64 `db:"id" json:"id"`
	//用户昵称
	Name string `db:"name" json:"name"`
	//用户电话
	NationCode string `db:"nation_code" json:"nationCode"`
	Phone      string `db:"phone" json:"phone"`
	//头像
	Avatar int64 `db:"avatar" json:"avatar"`
}

func SearchUser(args *ArgsSearchUser) (dataList []DataSearchUser, dataCount int64, err error) {
	where := "org_id = :org_id"
	maps := map[string]interface{}{
		"org_id": args.OrgID,
	}
	if args.Status > -1 || args.Status < 3 {
		where = where + " AND status = :status"
		maps["status"] = args.Status
	}
	if args.ParentSystem != "" {
		var parents []byte
		if args.ParentID > 0 {
			type parentArg struct {
				System   string `json:"system"`
				ParentID int64  `json:"parentID"`
			}
			parents, err = json.Marshal([]parentArg{
				{
					System:   args.ParentSystem,
					ParentID: args.ParentID,
				},
			})
			if err != nil {
				return
			}
		} else {
			type parentArg struct {
				System string `json:"system"`
			}
			parents, err = json.Marshal([]parentArg{
				{
					System: args.ParentSystem,
				},
			})
			if err != nil {
				return
			}
		}
		where = where + " AND parents = :parents"
		maps["parents"] = string(parents)
	}
	if !args.HaveRemove {
		where = where + " AND delete_at < to_timestamp(1000000)"
	}
	if args.Search != "" {
		maps["search"] = args.Search
		where = where + " AND (name ILIKE '%' || :search || '%' OR phone ILIKE '%' || :search || '%' OR email ILIKE '%' || :search || '%' OR username ILIKE '%' || :search || '%')"
	}
	var rawList []FieldsUserType
	tableName := "user_core"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&rawList,
		tableName,
		"id",
		fmt.Sprint("SELECT id ", "FROM ", tableName, " WHERE ", where),
		where,
		maps,
		&CoreSQLPages.ArgsDataList{
			Page: 1,
			Max:  args.Max,
			Sort: "update_at",
			Desc: true,
		},
		[]string{"update_at"},
	)
	if err != nil {
		return
	}
	searchInt64 := CoreFilter.GetInt64ByStringNoErr(args.Search)
	var findUserIDData FieldsUserType
	for _, v := range rawList {
		vData := getUserByID(v.ID)
		if vData.ID < 1 {
			continue
		}
		if vData.ID == searchInt64 {
			findUserIDData = vData
		}
		dataList = append(dataList, DataSearchUser{
			ID:         vData.ID,
			Name:       vData.Name,
			NationCode: vData.NationCode,
			Phone:      vData.Phone,
			Avatar:     vData.Avatar,
		})
	}
	if searchInt64 > 0 && findUserIDData.ID < 1 {
		vData := getUserByID(searchInt64)
		if vData.ID > 0 {
			dataList = append(dataList, DataSearchUser{
				ID:         vData.ID,
				Name:       vData.Name,
				NationCode: vData.NationCode,
				Phone:      vData.Phone,
				Avatar:     vData.Avatar,
			})
		}
	}
	return
}

type ArgsSearchUserByID struct {
	//用户ID
	ID int64 `db:"id" json:"id" check:"id"`
}

func SearchUserByID(args *ArgsSearchUserByID) (data DataSearchUser, err error) {
	var rawData FieldsUserType
	err = Router2SystemConfig.MainDB.Get(&rawData, "SELECT id FROM user_core WHERE id = $1", args.ID)
	if err != nil || rawData.ID < 1 {
		err = errors.New("no data")
		return
	}
	rawData = getUserByID(rawData.ID)
	if rawData.ID < 1 {
		err = errors.New("no data")
		return
	}
	data = DataSearchUser{
		ID:         rawData.ID,
		Name:       rawData.Name,
		NationCode: rawData.NationCode,
		Phone:      rawData.Phone,
		Avatar:     rawData.Avatar,
	}
	return
}

// ArgsGetUsers 获取一组用户ID参数
type ArgsGetUsers struct {
	//一组ID
	IDs pq.Int64Array `db:"ids" json:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove"`
	//是否不要显示电话
	NoPhone bool `json:"noPhone"`
}

// GetUsersName 获取一组用户ID
func GetUsersName(args *ArgsGetUsers) (data map[int64]string, err error) {
	//去重复
	newIDs := pq.Int64Array{}
	for _, v := range args.IDs {
		isFind := false
		for _, v2 := range newIDs {
			if v == v2 {
				isFind = true
				break
			}
		}
		if isFind {
			continue
		}
		newIDs = append(newIDs, v)
	}
	//获取数据
	type dataType struct {
		//ID
		ID int64 `db:"id" json:"id"`
		//名称
		Name string `db:"name" json:"name"`
		//电话
		Phone string `db:"phone" json:"phone"`
	}
	var dataList []dataType
	if args.NoPhone {
		err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, name FROM user_core WHERE id = ANY($1) AND ($2 = true OR delete_at < to_timestamp(1000000)) LIMIT 100", newIDs, args.HaveRemove)
	} else {
		err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, name, phone FROM user_core WHERE id = ANY($1) AND ($2 = true OR delete_at < to_timestamp(1000000)) LIMIT 100", newIDs, args.HaveRemove)
	}
	if err == nil {
		data = map[int64]string{}
		for _, v := range dataList {
			if args.NoPhone {
				data[v.ID] = v.Name
			} else {
				data[v.ID] = fmt.Sprint(v.Name, "(", v.Phone, ")")
			}
		}
	}
	return
}

// GetUserName 单独获取一个用户的信息参数
func GetUserName(id int64, noPhone bool) (data string) {
	if id < 1 {
		return
	}
	rawData, _ := GetUsersName(&ArgsGetUsers{
		IDs:        []int64{id},
		HaveRemove: false,
		NoPhone:    noPhone,
	})
	if len(rawData) > 0 {
		for _, v := range rawData {
			data = v
			return
		}
	}
	return
}

// 获取用户手机号
func GetUserPhone(id int64) (data string) {
	rawData := getUserByID(id)
	data = rawData.Phone
	return
}

func GetUserNameByID(id int64) string {
	if id < 1 {
		return ""
	}
	data := getUserByID(id)
	if data.ID < 1 {
		return ""
	}
	return data.Name
}

// ArgsGetUsersAndAvatar 获取一组用户ID和头像参数
type ArgsGetUsersAndAvatar struct {
	//一组ID
	IDs pq.Int64Array `db:"ids" json:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove"`
}

// DataGetUsersAndAvatar 获取一组用户ID和头像数据
type DataGetUsersAndAvatar struct {
	//用户ID
	ID int64 `db:"id" json:"id"`
	//用户昵称
	Name string `db:"name" json:"name"`
	//头像
	Avatar int64 `db:"avatar" json:"avatar"`
}

// GetUsersAndAvatar 获取一组用户ID和头像
func GetUsersAndAvatar(args *ArgsGetUsersAndAvatar) (dataList []DataGetUsersAndAvatar, err error) {
	err = CoreSQLIDs.GetIDsAndDelete(&dataList, "user_core", "id, name, avatar", args.IDs, args.HaveRemove)
	return
}

func GetUserAndAvatar(userID int64) (data DataGetUsersAndAvatar) {
	if userID < 1 {
		return
	}
	rawData := getUserByID(userID)
	if rawData.ID < 1 {
		return
	}
	data = DataGetUsersAndAvatar{
		ID:     rawData.ID,
		Name:   rawData.Name,
		Avatar: rawData.Avatar,
	}
	return
}

func GetUserAvatarUrl(userID int64) (data string) {
	if userID < 1 {
		return
	}
	rawData := getUserByID(userID)
	if rawData.ID < 1 {
		return
	}
	if rawData.Avatar < 1 {
		return
	}
	data = BaseFileSys2.GetPublicURLByClaimID(rawData.Avatar)
	return
}

// GetUsers 获取一组用户
func GetUsers(args *ArgsGetUsers) (dataList []DataGetUsersAndAvatar, err error) {
	err = CoreSQLIDs.GetIDsAndDelete(&dataList, "user_core", "id, name, avatar", args.IDs, args.HaveRemove)
	return
}

// GetUserCountByOrgID 获取组织下总人数
func GetUserCountByOrgID(orgID int64) (count int64) {
	if orgID < 0 {
		_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM user_core WHERE delete_at < to_timestamp(1000000)")
	} else {
		_ = Router2SystemConfig.MainDB.Get(&count, "SELECT COUNT(id) FROM user_core WHERE org_id = $1 AND delete_at < to_timestamp(1000000)", orgID)
	}
	return
}

// 获取指定用户数据
func getUserByID(userID int64) (data FieldsUserType) {
	cacheMark := getUserCacheMark(userID)
	if err := Router2SystemConfig.MainCache.GetStruct(cacheMark, &data); err == nil && data.ID > 0 {
		return
	}
	_ = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, status, org_id, name, password, nation_code, phone, email, username, avatar, parents, groups, infos, logins, sort_id, tags, phone_verify, email_verify FROM user_core WHERE id = $1", userID)
	if data.ID < 1 {
		return
	}
	Router2SystemConfig.MainCache.SetStruct(cacheMark, data, cacheUserTime)
	return
}
