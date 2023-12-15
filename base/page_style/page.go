package BasePageStyle

import (
	"encoding/json"
	"errors"
	CoreSQL "github.com/fotomxq/weeekj_core/v5/core/sql"
	CoreSQLConfig "github.com/fotomxq/weeekj_core/v5/core/sql/config"
	CoreSQLIDs "github.com/fotomxq/weeekj_core/v5/core/sql/ids"
	CoreSQLPages "github.com/fotomxq/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "github.com/fotomxq/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetPageList 获取页面列表参数
type ArgsGetPageList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//系统
	System string `db:"system" json:"system" check:"mark" empty:"true"`
	//页面识别码
	// 同一个系统下唯一
	Page string `db:"page" json:"page" check:"mark_page" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//是否删除
	IsRemove bool `db:"is_remove" json:"isRemove" check:"bool"`
	//搜索
	Search string `json:"search" check:"search" empty:"true"`
}

// GetPageList 获取页面列表
func GetPageList(args *ArgsGetPageList) (dataList []FieldsPage, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
	if args.OrgID > -1 {
		where = where + " AND org_id = :org_id"
		maps["org_id"] = args.OrgID
	}
	if args.System != "" {
		where = where + " AND system = :system"
		maps["system"] = args.System
	}
	if args.Page != "" {
		where = where + " AND page = :page"
		maps["page"] = args.Page
	}
	if args.SortID > -1 {
		where = where + " AND sort_id = :sort_id"
		maps["sort_id"] = args.SortID
	}
	if len(args.Tags) > 0 {
		where = where + " AND tags @> :tags"
		maps["tags"] = args.Tags
	}
	if args.Search != "" {
		where = where + " AND (title ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "core_page_style_page"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, org_id, system, page, title, component_list FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// ArgsGetPageIDs 获取一组页面参数
type ArgsGetPageIDs struct {
	//ID列
	IDs pq.Int64Array `json:"ids" check:"ids"`
	//是否包含删除数据
	HaveRemove bool `json:"haveRemove" check:"bool"`
}

// GetPageIDs 获取一组页面
func GetPageIDs(args *ArgsGetPageIDs) (dataList []FieldsPage, err error) {
	err = CoreSQLIDs.GetIDsAndDelete(&dataList, "core_page_style_page", "id, create_at, update_at, delete_at, org_id, system, page, title, data, component_list, params", args.IDs, args.HaveRemove)
	return
}

// ArgsGetPageMark 获取指定页面参数
type ArgsGetPageMark struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id"`
	//系统
	System string `db:"system" json:"system" check:"mark"`
	//页面识别码
	// 同一个系统下唯一
	Page string `db:"page" json:"page" check:"mark_page"`
}

// GetPageMark 获取指定页面
func GetPageMark(args *ArgsGetPageMark) (data FieldsPage, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, system, page, title, data, component_list, params FROM core_page_style_page WHERE system = $1 AND page = $2 AND org_id = $3 AND delete_at < to_timestamp(1000000)", args.System, args.Page, args.OrgID)
	if err == nil && data.ID > 0 {
		return
	}
	if args.OrgID != 0 {
		err = errors.New("no page")
		return
	}
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, org_id, system, page, title, data, component_list, params FROM core_page_style_page WHERE system = $1 AND page = $2 AND org_id = 0 AND delete_at < to_timestamp(1000000)", args.System, args.Page)
	return
}

// ArgsSetPage 修改页面参数
type ArgsSetPage struct {
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
	//系统
	System string `db:"system" json:"system" check:"mark"`
	//页面识别码
	// 同一个系统和商户ID下唯一
	Page string `db:"page" json:"page" check:"mark_page"`
	//标题
	Title string `db:"title" json:"title" check:"name" min:"1" max:"300"`
	//样式结构内容
	Data string `db:"data" json:"data"`
	//组件列
	ComponentList FieldsPageComponentList `db:"component_list" json:"componentList"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// SetPage 修改页面
func SetPage(args *ArgsSetPage) (err error) {
	//检查数据是否存在
	var id int64
	err = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM core_page_style_page WHERE org_id = $1 AND system = $2 AND page = $3 AND delete_at < to_timestamp(1000000)", args.OrgID, args.System, args.Page)
	if err == nil && id > 0 {
		//更新数据
		_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE core_page_style_page SET update_at = NOW(), title = :title, data = :data, component_list = :component_list, params = :params WHERE id = :id", map[string]interface{}{
			"id":             id,
			"title":          args.Title,
			"data":           args.Data,
			"component_list": args.ComponentList,
			"params":         args.Params,
		})
		return
	}
	//写入数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO core_page_style_page (org_id, system, page, title, data, component_list, params) VALUES (:org_id,:system,:page,:title,:data,:component_list,:params)", args)
	return
}

// ArgsDeletePage 删除页面参数
type ArgsDeletePage struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//组织ID
	OrgID int64 `db:"org_id" json:"orgID" check:"id" empty:"true"`
}

// DeletePage 删除页面
func DeletePage(args *ArgsDeletePage) (err error) {
	_, err = CoreSQL.DeleteAllSoft(Router2SystemConfig.MainDB.DB, "core_page_style_page", "id = :id AND (:org_id < 1 OR org_id = :org_id)", args)
	return
}

// OutputPage 导出页面数据
func OutputPage() (data string, err error) {
	var dataList []FieldsPage
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, update_at, delete_at, org_id, system, page, title, data, component_list, params FROM core_page_style_page WHERE delete_at < to_timestamp(1000000)")
	if err != nil || len(dataList) < 1 {
		err = errors.New("no data")
		return
	}
	var dataByte []byte
	dataByte, err = json.Marshal(dataList)
	if err != nil {
		return
	}
	data = string(dataByte)
	return
}

// ArgsImportPage 导入数据参数
type ArgsImportPage struct {
	//数据
	Data string `json:"data"`
}

// ImportPage 导入数据
func ImportPage(args *ArgsImportPage) (err error) {
	var dataList []FieldsPage
	err = json.Unmarshal([]byte(args.Data), &dataList)
	if err != nil {
		return
	}
	//准备写入数据包
	var insertData []interface{}
	for _, v := range dataList {
		//查询mark是否已经创建
		var id int64
		err = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM core_page_style_page WHERE system = $1 AND page = $2 AND delete_at < to_timestamp(1000000)", v.System, v.Page)
		if err == nil && id > 0 {
			continue
		}
		//写入数据
		insertData = append(insertData, ArgsSetPage{
			OrgID:         v.OrgID,
			System:        v.System,
			Page:          v.Page,
			Title:         v.Title,
			Data:          v.Data,
			ComponentList: v.ComponentList,
			Params:        v.Params,
		})
	}
	err = CoreSQL.CreateMore(Router2SystemConfig.MainDB.DB, "INSERT INTO core_page_style_page (org_id, system, page, title, data, component_list, params) VALUES (:org_id,:system,:page,:title,:data,:component_list,:params)", insertData)
	return
}
