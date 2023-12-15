package UserRole

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetTypeList 获取角色配置列表参数
type ArgsGetTypeList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//是否被删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool" empty:"true"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetTypeList 获取角色配置列表
func GetTypeList(args *ArgsGetTypeList) (dataList []FieldsType, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.Search != "" {
		where = where + " AND (name ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "user_role_type"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, mark, name, group_ids, params FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// ArgsGetTypeID 获取指定配置ID参数
type ArgsGetTypeID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// GetTypeID 获取指定配置ID
func GetTypeID(args *ArgsGetTypeID) (data FieldsType, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, mark, name, group_ids, params FROM user_role_type WHERE id = $1", args.ID)
	if err == nil && data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

// ArgsGetTypeMark 获取指定配置Mark参数
type ArgsGetTypeMark struct {
	//标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
}

// GetTypeMark 获取指定配置Mark
func GetTypeMark(args *ArgsGetTypeMark) (data FieldsType, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, mark, name, group_ids, params FROM user_role_type WHERE mark = $1 AND delete_at < to_timestamp(1000000)", args.Mark)
	if err == nil && data.ID < 1 {
		err = errors.New("no data")
		return
	}
	return
}

func GetTypeMarkNoErr(mark string) (data FieldsType) {
	err := Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, mark, name, group_ids, params FROM user_role_type WHERE mark = $1 AND delete_at < to_timestamp(1000000)", mark)
	if err == nil && data.ID < 1 {
		return
	}
	return
}

// ArgsCreateType 创建新的配置参数
type ArgsCreateType struct {
	//标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
	//配置名称
	Name string `db:"name" json:"name" check:"name"`
	//分配的用户组
	GroupIDs pq.Int64Array `db:"group_ids" json:"groupIDs" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params" check:"params" empty:"true"`
}

// CreateType 创建新的配置
func CreateType(args *ArgsCreateType) (data FieldsType, err error) {
	//mark不能重复
	if err = checkTypeMark(args.Mark, -1); err != nil {
		return
	}
	//创建数据
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "user_role_type", "INSERT INTO user_role_type (mark, name, group_ids, params) VALUES (:mark,:name,:group_ids,:params)", args, &data)
	return
}

// ArgsUpdateType 修改配置参数
type ArgsUpdateType struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//标识码
	Mark string `db:"mark" json:"mark" check:"mark"`
	//配置名称
	Name string `db:"name" json:"name" check:"name"`
	//分配的用户组
	GroupIDs pq.Int64Array `db:"group_ids" json:"groupIDs" check:"ids" empty:"true"`
	//扩展参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params" check:"params" empty:"true"`
}

// UpdateType 修改配置
func UpdateType(args *ArgsUpdateType) (err error) {
	//mark不能重复
	if err = checkTypeMark(args.Mark, args.ID); err != nil {
		return
	}
	//更新数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE user_role_type SET update_at = NOW(), mark = :mark, name = :name, group_ids = :group_ids, params = :params WHERE id = :id", args)
	return
}

// ArgsDeleteType 删除配置参数
type ArgsDeleteType struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteType 删除配置
func DeleteType(args *ArgsDeleteType) (err error) {
	//删除所有角色
	var count int64
	err = Router2SystemConfig.MainDB.Get(&count, "SELECT count(id) FROM user_role WHERE delete_at < to_timestamp(1000000) AND role_type = $1", args.ID)
	if count > 0 {
		_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "user_role", "role_type = :id", args)
		if err != nil {
			return
		}
	}
	//删除数据
	_, err = CoreSQL.DeleteOneSoft(Router2SystemConfig.MainDB.DB, "user_role_type", "id", args)
	return
}

// 检查mark是否重复
func checkTypeMark(mark string, id int64) (err error) {
	var findID int64
	err = Router2SystemConfig.MainDB.Get(&findID, "SELECT id FROM user_role_type WHERE mark = $1 AND delete_at < to_timestamp(1000000) AND ($2 < 1 OR id != $2)", mark, id)
	if err != nil || findID < 1 {
		return nil
	}
	err = errors.New("data is exist")
	return
}
