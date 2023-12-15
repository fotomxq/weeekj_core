package ToolsAppUpdate

import (
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetAppList 获取APP列表参数
type ArgsGetAppList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//搜索
	Search string `db:"search" json:"search" check:"search" empty:"true"`
}

// GetAppList 获取APP列表
func GetAppList(args *ArgsGetAppList) (dataList []FieldsApp, dataCount int64, err error) {
	maps := map[string]interface{}{}
	where := ""
	if args.OrgID > 0 {
		where = where + "org_id = :org_id"
		maps["org_id"] = args.OrgID
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
	tableName := "tools_app_update_app"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, org_id, name, app_mark, count FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at"},
	)
	return
}

// ArgsGetAppID 获取APP ID参数
type ArgsGetAppID struct {
	//APP ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// GetAppID 获取APP ID
func GetAppID(args *ArgsGetAppID) (data FieldsApp, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, org_id, name, des, des_files, app_mark, count FROM tools_app_update_app WHERE id = $1", args.ID)
	return
}

// ArgsCreateApp 创建新的APP参数
type ArgsCreateApp struct {
	//组织ID
	// 设备所属的组织，也可能为0
	// 组织也可以发布自己的APP
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//升级内容
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
	//描述附带文件
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
	//应用标识码
	AppMark string `db:"app_mark" json:"appMark" check:"mark"`
}

// CreateApp 创建新的APP
func CreateApp(args *ArgsCreateApp) (data FieldsApp, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id FROM tools_app_update_app WHERE app_mark = $1", args.AppMark)
	if err == nil && data.ID > 0 {
		err = errors.New("mark is exist")
		return
	}
	err = CoreSQL.CreateOneAndData(Router2SystemConfig.MainDB.DB, "tools_app_update_app", "INSERT INTO tools_app_update_app (org_id, name, des, des_files, app_mark, count) VALUES (:org_id,:name,:des,:des_files,:app_mark,0)", args, &data)
	return
}

// ArgsUpdateApp 修改APP参数
type ArgsUpdateApp struct {
	//APP ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	// 设备所属的组织，也可能为0
	// 组织也可以发布自己的APP
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//名称
	Name string `db:"name" json:"name" check:"name"`
	//升级内容
	Des string `db:"des" json:"des" check:"des" min:"1" max:"1000" empty:"true"`
	//描述附带文件
	DesFiles pq.Int64Array `db:"des_files" json:"desFiles" check:"ids" empty:"true"`
}

// UpdateApp 修改APP
func UpdateApp(args *ArgsUpdateApp) (err error) {
	if args.DesFiles == nil {
		args.DesFiles = []int64{}
	}
	_, err = CoreSQL.UpdateOne(Router2SystemConfig.MainDB.DB, "UPDATE tools_app_update_app SET name = :name, des = :des, des_files = :des_files WHERE id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// ArgsDeleteApp 删除APP参数
type ArgsDeleteApp struct {
	//ID
	ID int64 `db:"id" json:"id"`
	//组织ID
	// 设备所属的组织，也可能为0
	// 组织也可以发布自己的APP
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeleteApp 删除APP参数
func DeleteApp(args *ArgsDeleteApp) (err error) {
	_, err = CoreSQL.DeleteAll(Router2SystemConfig.MainDB.DB, "tools_app_update_app", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}
