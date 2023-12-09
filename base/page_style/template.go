package BasePageStyle

import (
	"encoding/json"
	"errors"
	CoreSQL "gitee.com/weeekj/weeekj_core/v5/core/sql"
	CoreSQLConfig "gitee.com/weeekj/weeekj_core/v5/core/sql/config"
	CoreSQLPages "gitee.com/weeekj/weeekj_core/v5/core/sql/pages"
	Router2SystemConfig "gitee.com/weeekj/weeekj_core/v5/router2/system_config"
	"github.com/lib/pq"
)

// ArgsGetTemplateList 获取模版列表参数
type ArgsGetTemplateList struct {
	//分页
	Pages CoreSQLPages.ArgsDataList `json:"pages"`
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

// GetTemplateList 获取模版列表
func GetTemplateList(args *ArgsGetTemplateList) (dataList []FieldsTemplate, dataCount int64, err error) {
	where := ""
	maps := map[string]interface{}{}
	where = CoreSQL.GetDeleteSQL(args.IsRemove, where)
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
		where = where + " AND (page ILIKE '%' || :search || '%' OR name ILIKE '%' || :search || '%' OR des ILIKE '%' || :search || '%')"
		maps["search"] = args.Search
	}
	tableName := "core_page_style_template"
	dataCount, err = CoreSQL.GetListPageAndCount(
		Router2SystemConfig.MainDB.DB,
		&dataList,
		tableName,
		"id",
		"SELECT id, create_at, update_at, delete_at, system, page, name, des, cover_file_id, sort_id, tags, org_sub_config_id, org_func_list, component_ids, default_component_ids FROM "+tableName+" WHERE "+where,
		where,
		maps,
		&args.Pages,
		[]string{"id", "create_at", "update_at", "delete_at"},
	)
	return
}

// ArgsGetTemplateID 获取指定模版参数
type ArgsGetTemplateID struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// GetTemplateID 获取指定模版
func GetTemplateID(args *ArgsGetTemplateID) (data FieldsTemplate, err error) {
	err = Router2SystemConfig.MainDB.Get(&data, "SELECT id, create_at, update_at, delete_at, system, page, name, des, cover_file_id, sort_id, tags, org_sub_config_id, org_func_list, component_ids, default_component_ids, data, default_data, params FROM core_page_style_template WHERE id = $1", args.ID)
	return
}

// ArgsCreateTemplate 创建新模版参数
type ArgsCreateTemplate struct {
	//系统
	System string `db:"system" json:"system" check:"mark"`
	//页面识别码
	// 同一个系统下，可能有多个相同页面的模版
	Page string `db:"page" json:"page" check:"mark_page"`
	//名称
	Name string `db:"name" json:"name" check:"name" min:"1" max:"300"`
	//介绍
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签组
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//商户订阅
	// 必须存在商户订阅配置的订阅，才能使用该组件
	OrgSubConfigID pq.Int64Array `db:"org_sub_config_id" json:"orgSubConfigID" check:"ids" empty:"true"`
	//商户功能
	// 只有开通相关功能，才能使用使用该组件
	OrgFuncList pq.StringArray `db:"org_func_list" json:"orgFuncList" check:"marks" empty:"true"`
	//可用组件列
	ComponentIDs pq.Int64Array `db:"component_ids" json:"componentIDs" check:"ids" empty:"true"`
	//默认呈现的组件排序
	DefaultComponentIDs pq.Int64Array `db:"default_component_ids" json:"defaultComponentIDs" check:"ids" empty:"true"`
	//样式结构内容
	Data string `db:"data" json:"data"`
	//默认样式结构
	DefaultData string `db:"default_data" json:"defaultData"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// CreateTemplate 创建新模版
func CreateTemplate(args *ArgsCreateTemplate) (err error) {
	//写入数据
	_, err = CoreSQL.CreateOne(Router2SystemConfig.MainDB.DB, "INSERT INTO core_page_style_template (system, page, name, des, cover_file_id, sort_id, tags, org_sub_config_id, org_func_list, component_ids, default_component_ids, data, default_data, params) VALUES (:system,:page,:name,:des,:cover_file_id,:sort_id,:tags,:org_sub_config_id,:org_func_list,:component_ids,:default_component_ids,:data,:default_data,:params)", args)
	return
}

// ArgsUpdateTemplate 修改模版参数
type ArgsUpdateTemplate struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
	//系统
	System string `db:"system" json:"system" check:"mark"`
	//页面识别码
	// 同一个系统下唯一
	Page string `db:"page" json:"page" check:"mark_page"`
	//名称
	Name string `db:"name" json:"name" check:"name" min:"1" max:"300"`
	//介绍
	Des string `db:"des" json:"des" check:"des" min:"1" max:"3000" empty:"true"`
	//封面
	CoverFileID int64 `db:"cover_file_id" json:"coverFileID" check:"id" empty:"true"`
	//分类ID
	SortID int64 `db:"sort_id" json:"sortID" check:"id" empty:"true"`
	//标签组
	Tags pq.Int64Array `db:"tags" json:"tags" check:"ids" empty:"true"`
	//商户订阅
	// 必须存在商户订阅配置的订阅，才能使用该组件
	OrgSubConfigID pq.Int64Array `db:"org_sub_config_id" json:"orgSubConfigID" check:"ids" empty:"true"`
	//商户功能
	// 只有开通相关功能，才能使用使用该组件
	OrgFuncList pq.StringArray `db:"org_func_list" json:"orgFuncList" check:"marks" empty:"true"`
	//可用组件列
	ComponentIDs pq.Int64Array `db:"component_ids" json:"componentIDs" check:"ids" empty:"true"`
	//默认呈现的组件排序
	DefaultComponentIDs pq.Int64Array `db:"default_component_ids" json:"defaultComponentIDs" check:"ids" empty:"true"`
	//样式结构内容
	Data string `db:"data" json:"data"`
	//默认样式结构
	DefaultData string `db:"default_data" json:"defaultData"`
	//附加参数
	Params CoreSQLConfig.FieldsConfigsType `db:"params" json:"params"`
}

// UpdateTemplate 修改模版
func UpdateTemplate(args *ArgsUpdateTemplate) (err error) {
	//更新数据
	_, err = CoreSQL.UpdateOneSoft(Router2SystemConfig.MainDB.DB, "UPDATE core_page_style_template SET update_at = NOW(), system = :system, page = :page, name = :name, des = :des, cover_file_id = :cover_file_id, sort_id = :sort_id, tags = :tags, org_sub_config_id = :org_sub_config_id, org_func_list = :org_func_list, component_ids = :component_ids, default_component_ids = :default_component_ids, data = :data, default_data = :default_data, params = :params WHERE id = :id", args)
	return
}

// ArgsDeleteTemplate 删除模版参数
type ArgsDeleteTemplate struct {
	//ID
	ID int64 `db:"id" json:"id" check:"id"`
}

// DeleteTemplate 删除模版
func DeleteTemplate(args *ArgsDeleteTemplate) (err error) {
	_, err = CoreSQL.DeleteOneSoft(Router2SystemConfig.MainDB.DB, "core_page_style_template", "id", args)
	return
}

// OutputTemplate 导出模版数据
func OutputTemplate() (data string, err error) {
	var dataList []FieldsTemplate
	err = Router2SystemConfig.MainDB.Select(&dataList, "SELECT id, create_at, update_at, delete_at, system, page, name, des, cover_file_id, sort_id, tags, org_sub_config_id, org_func_list, component_ids, default_component_ids, data, default_data, params FROM core_page_style_template WHERE delete_at < to_timestamp(1000000)")
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

// ArgsImportTemplate 导入模版参数
type ArgsImportTemplate struct {
	//数据
	Data string `json:"data"`
}

// ImportTemplate 导入数据
func ImportTemplate(args *ArgsImportTemplate) (err error) {
	var dataList []FieldsTemplate
	err = json.Unmarshal([]byte(args.Data), &dataList)
	if err != nil {
		return
	}
	//准备写入数据包
	var insertData []interface{}
	for _, v := range dataList {
		//查询mark是否已经创建
		var id int64
		err = Router2SystemConfig.MainDB.Get(&id, "SELECT id FROM core_page_style_template WHERE system = $1 AND page = $2 AND name = $3 AND delete_at < to_timestamp(1000000)", v.System, v.Page, v.Name)
		if err == nil && id > 0 {
			continue
		}
		//写入数据
		insertData = append(insertData, ArgsCreateTemplate{
			System:              v.System,
			Page:                v.Page,
			Name:                v.Name,
			Des:                 v.Des,
			CoverFileID:         v.CoverFileID,
			SortID:              v.SortID,
			Tags:                v.Tags,
			OrgSubConfigID:      v.OrgSubConfigID,
			OrgFuncList:         v.OrgFuncList,
			ComponentIDs:        v.ComponentIDs,
			DefaultComponentIDs: v.DefaultComponentIDs,
			Data:                v.Data,
			DefaultData:         v.DefaultData,
			Params:              v.Params,
		})
	}
	err = CoreSQL.CreateMore(Router2SystemConfig.MainDB.DB, "INSERT INTO core_page_style_template (system, page, name, des, cover_file_id, sort_id, tags, org_sub_config_id, org_func_list, component_ids, default_component_ids, data, default_data, params) VALUES (:system,:page,:name,:des,:cover_file_id,:sort_id,:tags,:org_sub_config_id,:org_func_list,:component_ids,:default_component_ids,:data,:default_data,:params)", insertData)
	return
}
